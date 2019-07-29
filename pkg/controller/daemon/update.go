/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package daemon

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"

	"k8s.io/klog"

	apps "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	intstrutil "k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/sets"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
	"k8s.io/kubernetes/pkg/controller"
	"k8s.io/kubernetes/pkg/controller/daemon/util"
	labelsutil "k8s.io/kubernetes/pkg/util/labels"
)

// rollingUpdate deletes old daemon set pods making sure that no more than
// ds.Spec.UpdateStrategy.RollingUpdate.MaxUnavailable pods are unavailable
func (dsc *DaemonSetsController) rollingUpdate(ds *apps.DaemonSet, nodeList []*v1.Node, hash string) error {
	nodeToDaemonPods, err := dsc.getNodesToDaemonPods(ds)

	if err != nil {
		return fmt.Errorf("couldn't get node to daemon pod mapping for daemon set %q: %v", ds.Name, err)
	}

	// respect gray update choice.
	if ds.Spec.UpdateStrategy.RollingUpdate.Selector == nil {
		if ds.Spec.UpdateStrategy.RollingUpdate.Partition != nil &&  *ds.Spec.UpdateStrategy.RollingUpdate.Partition != 0 {
			// respect partitioned nodes to keep old versions.
			nodeToDaemonPods = dsc.getNodeToDaemonPodsByPartition(ds, nodeToDaemonPods)
		}
	} else {
		// respect selected nodes to update.
		nodeToDaemonPods = dsc.getNodeToDaemonPodsBySelector(ds, nodeToDaemonPods)
	}

	_, oldPods := dsc.getAllDaemonSetPods(ds, nodeToDaemonPods, hash)
	maxUnavailable, numUnavailable, err := dsc.getUnavailableNumbers(ds, nodeList, nodeToDaemonPods)
	if err != nil {
		return fmt.Errorf("couldn't get unavailable numbers: %v", err)
	}
	oldAvailablePods, oldUnavailablePods := util.SplitByAvailablePods(ds.Spec.MinReadySeconds, oldPods)

	// for oldPods delete all not running pods
	var oldPodsToDelete []string
	klog.V(4).Infof("Marking all unavailable old pods for deletion")
	for _, pod := range oldUnavailablePods {
		// Skip terminating pods. We won't delete them again
		if pod.DeletionTimestamp != nil {
			continue
		}
		klog.V(4).Infof("Marking pod %s/%s for deletion", ds.Name, pod.Name)
		oldPodsToDelete = append(oldPodsToDelete, pod.Name)
	}

	klog.V(4).Infof("Marking old pods for deletion")
	for _, pod := range oldAvailablePods {
		if numUnavailable >= maxUnavailable {
			klog.V(4).Infof("Number of unavailable DaemonSet pods: %d, is equal to or exceeds allowed maximum: %d", numUnavailable, maxUnavailable)
			break
		}
		klog.V(4).Infof("Marking pod %s/%s for deletion", ds.Name, pod.Name)
		oldPodsToDelete = append(oldPodsToDelete, pod.Name)
		numUnavailable++
	}
	return dsc.syncNodes(ds, oldPodsToDelete, []string{}, hash)
}

func (dsc *DaemonSetsController) getNodeToDaemonPodsByPartition(ds *apps.DaemonSet, nodeToDaemonPods map[string][]*v1.Pod) map[string][]*v1.Pod {
	// sort Nodes by their Names
	var nodeNames []string
	for nodeName := range nodeToDaemonPods {
		nodeNames = append(nodeNames, nodeName)
	}
	sort.Strings(nodeNames)

	var partition int32
	switch ds.Spec.UpdateStrategy.Type {
	case apps.RollingUpdateDaemonSetStrategyType:
		partition = *ds.Spec.UpdateStrategy.RollingUpdate.Partition
	case apps.SurgingRollingUpdateDaemonSetStrategyType:
		partition = *ds.Spec.UpdateStrategy.SurgingRollingUpdate.Partition
	}

	// keep the old version pods whose count is no more than Partition value.
	for i := int32(0); i < int32(len(nodeNames)); i++ {
		if i == partition {
			break
		}
		delete(nodeToDaemonPods, nodeNames[i])
	}
	return nodeToDaemonPods
}

func (dsc *DaemonSetsController) getNodeToDaemonPodsBySelector(ds *apps.DaemonSet, nodeToDaemonPods map[string][]*v1.Pod) map[string][]*v1.Pod {
	var selector labels.Selector
	var err error
	switch ds.Spec.UpdateStrategy.Type {
	case apps.OnDeleteDaemonSetStrategyType:
		return nodeToDaemonPods
	case apps.RollingUpdateDaemonSetStrategyType:
		selector, err = metav1.LabelSelectorAsSelector(ds.Spec.UpdateStrategy.RollingUpdate.Selector)
		if err != nil {
			return nodeToDaemonPods
		}
	case apps.SurgingRollingUpdateDaemonSetStrategyType:
		selector, err = metav1.LabelSelectorAsSelector(ds.Spec.UpdateStrategy.SurgingRollingUpdate.Selector)
		if err != nil {
			return nodeToDaemonPods
		}
	}

	for nodeName := range nodeToDaemonPods {
		node, err := dsc.nodeLister.Get(nodeName)
		if err != nil {
			klog.Errorf("could not get node: %s nodeInfo", nodeName)
			continue
		}
		if !selector.Matches(labels.Set(node.Labels)) {
			delete(nodeToDaemonPods, nodeName)
		}
	}

	return nodeToDaemonPods
}

// getSurgeNumbers returns the max allowable number of surging pods and the current number of
// surging pods. The number of surging pods is computed as the total number pods above the first
// on each node.
func (dsc *DaemonSetsController) getSurgeNumbers(ds *apps.DaemonSet, nodeToDaemonPods map[string][]*v1.Pod, hash string) (int, int, error) {
	klog.V(4).Infof("Getting surge numbers")
	// TODO: get nodeList once in syncDaemonSet and pass it to other functions
	nodeList, err := dsc.nodeLister.List(labels.Everything())
	if err != nil {
		return -1, -1, fmt.Errorf("couldn't get list of nodes during surging rolling update of daemon set %#v: %v", ds, err)
	}

	generation, err := util.GetTemplateGeneration(ds)
	if err != nil {
		generation = nil
	}
	var desiredNumberScheduled, numSurge int
	for i := range nodeList {
		node := nodeList[i]
		wantToRun, _, _, err := dsc.nodeShouldRunDaemonPod(node, ds)
		if err != nil {
			return -1, -1, err
		}
		if !wantToRun {
			continue
		}

		desiredNumberScheduled++

		for _, pod := range nodeToDaemonPods[node.Name] {
			if util.IsPodUpdated(pod, hash, generation) && len(nodeToDaemonPods[node.Name]) > 1 {
				numSurge += 1
				break
			}
		}
	}

	maxSurge, err := intstrutil.GetValueFromIntOrPercent(ds.Spec.UpdateStrategy.SurgingRollingUpdate.MaxSurge, desiredNumberScheduled, true)
	if err != nil {
		return -1, -1, fmt.Errorf("Invalid value for MaxSurge: %v", err)
	}
	klog.V(4).Infof("DaemonSet %s/%s, maxSurge: %d, numSurge: %d", ds.Namespace, ds.Name, maxSurge, numSurge)
	return maxSurge, numSurge, nil
}

// get all nodes where the ds should run
func (dsc *DaemonSetsController) getNodesShouldRunDaemonPod(ds *apps.DaemonSet) (
	nodesWantToRun sets.String,
	nodesShouldContinueRunning sets.String,
	err error) {
	nodeList, err := dsc.nodeLister.List(labels.Everything())
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't get list of nodes when updating daemon set %#v: %v", ds, err)
	}

	nodesWantToRun = sets.String{}
	nodesShouldContinueRunning = sets.String{}
	for _, node := range nodeList {
		wantToRun, _, shouldContinueRunning, err := dsc.nodeShouldRunDaemonPod(node, ds)
		if err != nil {
			return nil, nil, err
		}

		if wantToRun {
			nodesWantToRun.Insert(node.Name)
		}

		if shouldContinueRunning {
			nodesShouldContinueRunning.Insert(node.Name)
		}
	}
	return nodesWantToRun, nodesShouldContinueRunning, nil
}

// pruneSurgingDaemonPods prunes the list of pods to only contain pods belonging to the current
// generation. This method only applies when the update strategy is SurgingRollingUpdate.
// This allows the daemon set controller to temporarily break its contract that only one daemon
// pod can run per node by ignoring the pods that belong to previous generations, which are
// cleaned up by the surgingRollineUpdate() method above.
func (dsc *DaemonSetsController) pruneSurgingDaemonPods(ds *apps.DaemonSet, pods []*v1.Pod, hash string) []*v1.Pod {
	if len(pods) <= 1 || ds.Spec.UpdateStrategy.Type != apps.SurgingRollingUpdateDaemonSetStrategyType {
		return pods
	}
	var currentPods []*v1.Pod
	for _, pod := range pods {
		generation, err := util.GetTemplateGeneration(ds)
		if err != nil {
			generation = nil
		}
		if util.IsPodUpdated(pod, hash, generation) {
			klog.V(5).Infof("Pod %s of ds %s/%s already updated(%s)", pod.Name, ds.Namespace, ds.Name, hash)
			currentPods = append(currentPods, pod)
		}
	}
	// Escape hatch if no new pods of the current generation are present yet.
	if len(currentPods) == 0 {
		klog.V(5).Infof("No new version(%s) pods for %s/%s", hash, ds.Namespace, ds.Name)
		currentPods = pods
	}
	return currentPods
}

// surgingRollingUpdate creates new daemon set pods to replace old ones making sure that no more
// than ds.Spec.UpdateStrategy.SurgingRollingUpdate.MaxSurge extra pods are scheduled at any
// given time.
func (dsc *DaemonSetsController) surgingRollingUpdate(ds *apps.DaemonSet, nodeList []*v1.Node, hash string) error {
	nodeToDaemonPods, err := dsc.getNodesToDaemonPods(ds)
	if err != nil {
		return fmt.Errorf("couldn't get node to daemon pod mapping for daemon set %q: %v", ds.Name, err)
	}

	// respect gray update choice.
	if ds.Spec.UpdateStrategy.SurgingRollingUpdate.Selector == nil {
		if ds.Spec.UpdateStrategy.SurgingRollingUpdate.Partition != nil && *ds.Spec.UpdateStrategy.SurgingRollingUpdate.Partition !=0 {
			// respect partitioned nodes to keep old versions.
			nodeToDaemonPods = dsc.getNodeToDaemonPodsByPartition(ds, nodeToDaemonPods)
		}
	} else {
		// respect selected nodes to update.
		nodeToDaemonPods = dsc.getNodeToDaemonPodsBySelector(ds, nodeToDaemonPods)
	}

	maxSurge, numSurge, err := dsc.getSurgeNumbers(ds, nodeToDaemonPods, hash)
	if err != nil {
		return fmt.Errorf("Couldn't get surge numbers: %v", err)
	}

	nodesWantToRun, nodesShouldContinueRunning, err := dsc.getNodesShouldRunDaemonPod(ds)
	if err != nil {
		return fmt.Errorf("Couldn't get nodes which want to run ds pod: %v", err)
	}

	klog.V(4).Infof("Surging new pods and deleting obsolete old pods")
	var nodesToSurge []string
	var oldPodsToDelete []string
	for node, pods := range nodeToDaemonPods {
		var newPod, oldPod *v1.Pod
		wantToRun := nodesWantToRun.Has(node)
		shouldContinueRunning := nodesShouldContinueRunning.Has(node)
		// if node has new taint, then the already existed pod does not want to run but should continue running
		if !wantToRun && !shouldContinueRunning {
			for _, pod := range pods {
				klog.V(4).Infof("Marking pod %s/%s on unsuitable node %s for deletion", ds.Name, pod.Name, node)
				if pod.DeletionTimestamp == nil {
					oldPodsToDelete = append(oldPodsToDelete, pod.Name)
				}
			}
			continue
		}
		foundAvailable := false
		for _, pod := range pods {
			generation, err := util.GetTemplateGeneration(ds)
			if err != nil {
				generation = nil
			}
			if util.IsPodUpdated(pod, hash, generation) {
				if newPod != nil {
					klog.Warningf("Multiple new pods on node %s: %s, %s", node, newPod.Name, pod.Name)
				}
				newPod = pod
			} else {
				if oldPod != nil {
					klog.Warningf("Multiple old pods on node %s: %s, %s", node, oldPod.Name, pod.Name)
				}
				oldPod = pod
			}
			if !foundAvailable && podutil.IsPodAvailable(pod, ds.Spec.MinReadySeconds, metav1.Now()) {
				foundAvailable = true
			}
		}

		if newPod == nil && numSurge < maxSurge && wantToRun {
			if !foundAvailable && len(pods) >= 2 {
				klog.Warningf("Node %s already has %d unavailble pods, need clean first, skip surge new pod", node, len(pods))
			} else {
				klog.V(4).Infof("Surging new pod on node %s", node)
				numSurge++
				nodesToSurge = append(nodesToSurge, node)
			}
		} else if newPod != nil && podutil.IsPodAvailable(newPod, ds.Spec.MinReadySeconds, metav1.Now()) && oldPod != nil && oldPod.DeletionTimestamp == nil {
			klog.V(4).Infof("Marking pod %s/%s for deletion", ds.Name, oldPod.Name)
			oldPodsToDelete = append(oldPodsToDelete, oldPod.Name)
		}
	}
	return dsc.syncNodes(ds, oldPodsToDelete, nodesToSurge, hash)
}

// constructHistory finds all histories controlled by the given DaemonSet, and
// update current history revision number, or create current history if need to.
// It also deduplicates current history, and adds missing unique labels to existing histories.
func (dsc *DaemonSetsController) constructHistory(ds *apps.DaemonSet) (cur *apps.ControllerRevision, old []*apps.ControllerRevision, err error) {
	var histories []*apps.ControllerRevision
	var currentHistories []*apps.ControllerRevision
	histories, err = dsc.controlledHistories(ds)
	if err != nil {
		return nil, nil, err
	}
	for _, history := range histories {
		// Add the unique label if it's not already added to the history
		// We use history name instead of computing hash, so that we don't need to worry about hash collision
		if _, ok := history.Labels[apps.DefaultDaemonSetUniqueLabelKey]; !ok {
			toUpdate := history.DeepCopy()
			toUpdate.Labels[apps.DefaultDaemonSetUniqueLabelKey] = toUpdate.Name
			history, err = dsc.kubeClient.AppsV1().ControllerRevisions(ds.Namespace).Update(toUpdate)
			if err != nil {
				return nil, nil, err
			}
		}
		// Compare histories with ds to separate cur and old history
		found := false
		found, err = Match(ds, history)
		if err != nil {
			return nil, nil, err
		}
		if found {
			currentHistories = append(currentHistories, history)
		} else {
			old = append(old, history)
		}
	}

	currRevision := maxRevision(old) + 1
	switch len(currentHistories) {
	case 0:
		// Create a new history if the current one isn't found
		cur, err = dsc.snapshot(ds, currRevision)
		if err != nil {
			return nil, nil, err
		}
	default:
		cur, err = dsc.dedupCurHistories(ds, currentHistories)
		if err != nil {
			return nil, nil, err
		}
		// Update revision number if necessary
		if cur.Revision < currRevision {
			toUpdate := cur.DeepCopy()
			toUpdate.Revision = currRevision
			_, err = dsc.kubeClient.AppsV1().ControllerRevisions(ds.Namespace).Update(toUpdate)
			if err != nil {
				return nil, nil, err
			}
		}
	}
	return cur, old, err
}

func (dsc *DaemonSetsController) cleanupHistory(ds *apps.DaemonSet, old []*apps.ControllerRevision) error {
	nodesToDaemonPods, err := dsc.getNodesToDaemonPods(ds)
	if err != nil {
		return fmt.Errorf("couldn't get node to daemon pod mapping for daemon set %q: %v", ds.Name, err)
	}

	toKeep := int(*ds.Spec.RevisionHistoryLimit)
	toKill := len(old) - toKeep
	if toKill <= 0 {
		return nil
	}

	// Find all hashes of live pods
	liveHashes := make(map[string]bool)
	for _, pods := range nodesToDaemonPods {
		for _, pod := range pods {
			if hash := pod.Labels[apps.DefaultDaemonSetUniqueLabelKey]; len(hash) > 0 {
				liveHashes[hash] = true
			}
		}
	}

	// Clean up old history from smallest to highest revision (from oldest to newest)
	sort.Sort(historiesByRevision(old))
	for _, history := range old {
		if toKill <= 0 {
			break
		}
		if hash := history.Labels[apps.DefaultDaemonSetUniqueLabelKey]; liveHashes[hash] {
			continue
		}
		// Clean up
		err := dsc.kubeClient.AppsV1().ControllerRevisions(ds.Namespace).Delete(history.Name, nil)
		if err != nil {
			return err
		}
		toKill--
	}
	return nil
}

// maxRevision returns the max revision number of the given list of histories
func maxRevision(histories []*apps.ControllerRevision) int64 {
	max := int64(0)
	for _, history := range histories {
		if history.Revision > max {
			max = history.Revision
		}
	}
	return max
}

func (dsc *DaemonSetsController) dedupCurHistories(ds *apps.DaemonSet, curHistories []*apps.ControllerRevision) (*apps.ControllerRevision, error) {
	if len(curHistories) == 1 {
		return curHistories[0], nil
	}
	var maxRevision int64
	var keepCur *apps.ControllerRevision
	for _, cur := range curHistories {
		if cur.Revision >= maxRevision {
			keepCur = cur
			maxRevision = cur.Revision
		}
	}
	// Clean up duplicates and relabel pods
	for _, cur := range curHistories {
		if cur.Name == keepCur.Name {
			continue
		}
		// Relabel pods before dedup
		pods, err := dsc.getDaemonPods(ds)
		if err != nil {
			return nil, err
		}
		for _, pod := range pods {
			if pod.Labels[apps.DefaultDaemonSetUniqueLabelKey] != keepCur.Labels[apps.DefaultDaemonSetUniqueLabelKey] {
				toUpdate := pod.DeepCopy()
				if toUpdate.Labels == nil {
					toUpdate.Labels = make(map[string]string)
				}
				toUpdate.Labels[apps.DefaultDaemonSetUniqueLabelKey] = keepCur.Labels[apps.DefaultDaemonSetUniqueLabelKey]
				_, err = dsc.kubeClient.CoreV1().Pods(ds.Namespace).Update(toUpdate)
				if err != nil {
					return nil, err
				}
			}
		}
		// Remove duplicates
		err = dsc.kubeClient.AppsV1().ControllerRevisions(ds.Namespace).Delete(cur.Name, nil)
		if err != nil {
			return nil, err
		}
	}
	return keepCur, nil
}

// controlledHistories returns all ControllerRevisions controlled by the given DaemonSet.
// This also reconciles ControllerRef by adopting/orphaning.
// Note that returned histories are pointers to objects in the cache.
// If you want to modify one, you need to deep-copy it first.
func (dsc *DaemonSetsController) controlledHistories(ds *apps.DaemonSet) ([]*apps.ControllerRevision, error) {
	selector, err := metav1.LabelSelectorAsSelector(ds.Spec.Selector)
	if err != nil {
		return nil, err
	}

	// List all histories to include those that don't match the selector anymore
	// but have a ControllerRef pointing to the controller.
	histories, err := dsc.historyLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}
	// If any adoptions are attempted, we should first recheck for deletion with
	// an uncached quorum read sometime after listing Pods (see #42639).
	canAdoptFunc := controller.RecheckDeletionTimestamp(func() (metav1.Object, error) {
		fresh, err := dsc.kubeClient.AppsV1().DaemonSets(ds.Namespace).Get(ds.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if fresh.UID != ds.UID {
			return nil, fmt.Errorf("original DaemonSet %v/%v is gone: got uid %v, wanted %v", ds.Namespace, ds.Name, fresh.UID, ds.UID)
		}
		return fresh, nil
	})
	// Use ControllerRefManager to adopt/orphan as needed.
	cm := controller.NewControllerRevisionControllerRefManager(dsc.crControl, ds, selector, controllerKind, canAdoptFunc)
	return cm.ClaimControllerRevisions(histories)
}

// Match check if the given DaemonSet's template matches the template stored in the given history.
func Match(ds *apps.DaemonSet, history *apps.ControllerRevision) (bool, error) {
	patch, err := getPatch(ds)
	if err != nil {
		return false, err
	}
	return bytes.Equal(patch, history.Data.Raw), nil
}

// getPatch returns a strategic merge patch that can be applied to restore a Daemonset to a
// previous version. If the returned error is nil the patch is valid. The current state that we save is just the
// PodSpecTemplate. We can modify this later to encompass more state (or less) and remain compatible with previously
// recorded patches.
func getPatch(ds *apps.DaemonSet) ([]byte, error) {
	dsBytes, err := json.Marshal(ds)
	if err != nil {
		return nil, err
	}
	var raw map[string]interface{}
	err = json.Unmarshal(dsBytes, &raw)
	if err != nil {
		return nil, err
	}
	objCopy := make(map[string]interface{})
	specCopy := make(map[string]interface{})

	// Create a patch of the DaemonSet that replaces spec.template
	spec := raw["spec"].(map[string]interface{})
	template := spec["template"].(map[string]interface{})
	specCopy["template"] = template
	template["$patch"] = "replace"
	objCopy["spec"] = specCopy
	patch, err := json.Marshal(objCopy)
	return patch, err
}

func (dsc *DaemonSetsController) snapshot(ds *apps.DaemonSet, revision int64) (*apps.ControllerRevision, error) {
	patch, err := getPatch(ds)
	if err != nil {
		return nil, err
	}
	hash := controller.ComputeHash(&ds.Spec.Template, ds.Status.CollisionCount)
	name := ds.Name + "-" + hash
	history := &apps.ControllerRevision{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       ds.Namespace,
			Labels:          labelsutil.CloneAndAddLabel(ds.Spec.Template.Labels, apps.DefaultDaemonSetUniqueLabelKey, hash),
			Annotations:     ds.Annotations,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(ds, controllerKind)},
		},
		Data:     runtime.RawExtension{Raw: patch},
		Revision: revision,
	}

	history, err = dsc.kubeClient.AppsV1().ControllerRevisions(ds.Namespace).Create(history)
	if outerErr := err; errors.IsAlreadyExists(outerErr) {
		// TODO: Is it okay to get from historyLister?
		existedHistory, getErr := dsc.kubeClient.AppsV1().ControllerRevisions(ds.Namespace).Get(name, metav1.GetOptions{})
		if getErr != nil {
			return nil, getErr
		}
		// Check if we already created it
		done, matchErr := Match(ds, existedHistory)
		if matchErr != nil {
			return nil, matchErr
		}
		if done {
			return existedHistory, nil
		}

		// Handle name collisions between different history
		// Get the latest DaemonSet from the API server to make sure collision count is only increased when necessary
		currDS, getErr := dsc.kubeClient.AppsV1().DaemonSets(ds.Namespace).Get(ds.Name, metav1.GetOptions{})
		if getErr != nil {
			return nil, getErr
		}
		// If the collision count used to compute hash was in fact stale, there's no need to bump collision count; retry again
		if !reflect.DeepEqual(currDS.Status.CollisionCount, ds.Status.CollisionCount) {
			return nil, fmt.Errorf("found a stale collision count (%d, expected %d) of DaemonSet %q while processing; will retry until it is updated", ds.Status.CollisionCount, currDS.Status.CollisionCount, ds.Name)
		}
		if currDS.Status.CollisionCount == nil {
			currDS.Status.CollisionCount = new(int32)
		}
		*currDS.Status.CollisionCount++
		_, updateErr := dsc.kubeClient.AppsV1().DaemonSets(ds.Namespace).UpdateStatus(currDS)
		if updateErr != nil {
			return nil, updateErr
		}
		klog.V(2).Infof("Found a hash collision for DaemonSet %q - bumping collisionCount to %d to resolve it", ds.Name, *currDS.Status.CollisionCount)
		return nil, outerErr
	}
	return history, err
}

func (dsc *DaemonSetsController) getAllDaemonSetPods(ds *apps.DaemonSet, nodeToDaemonPods map[string][]*v1.Pod, hash string) ([]*v1.Pod, []*v1.Pod) {
	var newPods []*v1.Pod
	var oldPods []*v1.Pod

	for _, pods := range nodeToDaemonPods {
		for _, pod := range pods {
			// If the returned error is not nil we have a parse error.
			// The controller handles this via the hash.
			generation, err := util.GetTemplateGeneration(ds)
			if err != nil {
				generation = nil
			}
			if util.IsPodUpdated(pod, hash, generation) {
				newPods = append(newPods, pod)
			} else {
				oldPods = append(oldPods, pod)
			}
		}
	}
	return newPods, oldPods
}

func (dsc *DaemonSetsController) getUnavailableNumbers(ds *apps.DaemonSet, nodeList []*v1.Node, nodeToDaemonPods map[string][]*v1.Pod) (int, int, error) {
	klog.V(4).Infof("Getting unavailable numbers")
	var numUnavailable, desiredNumberScheduled int
	for i := range nodeList {
		node := nodeList[i]
		wantToRun, _, _, err := dsc.nodeShouldRunDaemonPod(node, ds)
		if err != nil {
			return -1, -1, err
		}
		if !wantToRun {
			continue
		}
		desiredNumberScheduled++
		daemonPods, exists := nodeToDaemonPods[node.Name]
		if !exists {
			numUnavailable++
			continue
		}
		available := false
		for _, pod := range daemonPods {
			//for the purposes of update we ensure that the Pod is both available and not terminating
			if podutil.IsPodAvailable(pod, ds.Spec.MinReadySeconds, metav1.Now()) && pod.DeletionTimestamp == nil {
				available = true
				break
			}
		}
		if !available {
			numUnavailable++
		}
	}
	maxUnavailable, err := intstrutil.GetValueFromIntOrPercent(ds.Spec.UpdateStrategy.RollingUpdate.MaxUnavailable, desiredNumberScheduled, true)
	if err != nil {
		return -1, -1, fmt.Errorf("invalid value for MaxUnavailable: %v", err)
	}
	klog.V(4).Infof(" DaemonSet %s/%s, maxUnavailable: %d, numUnavailable: %d", ds.Namespace, ds.Name, maxUnavailable, numUnavailable)
	return maxUnavailable, numUnavailable, nil
}

type historiesByRevision []*apps.ControllerRevision

func (h historiesByRevision) Len() int      { return len(h) }
func (h historiesByRevision) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h historiesByRevision) Less(i, j int) bool {
	return h[i].Revision < h[j].Revision
}
