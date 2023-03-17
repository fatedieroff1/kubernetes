/*
Copyright 2020 The Kubernetes Authors.

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

package e2enode

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/prometheus/common/model"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/test/e2e/framework"
	imageutils "k8s.io/kubernetes/test/utils/image"
	admissionapi "k8s.io/pod-security-admission/api"
	utilpointer "k8s.io/utils/pointer"
)

var _ = SIGDescribe("MirrorPodWithGracePeriod", func() {
	f := framework.NewDefaultFramework("mirror-pod-with-grace-period")
	f.NamespacePodSecurityEnforceLevel = admissionapi.LevelBaseline
	ginkgo.Context("when create a mirror pod ", func() {
		var ns, podPath, staticPodName, mirrorPodName string
		var staticPod *v1.Pod
		ginkgo.BeforeEach(func(ctx context.Context) {
			ns = f.Namespace.Name
			staticPodName = "graceful-pod-" + string(uuid.NewUUID())
			mirrorPodName = staticPodName + "-" + framework.TestContext.NodeName

			podPath = framework.TestContext.KubeletConfig.StaticPodPath

			ginkgo.By("create the static pod")
			staticPod = createStaticPodWithGracePeriod(podPath, staticPodName, ns)
			err := writeStaticPodManifest(podPath, staticPodName, ns, staticPod)
			framework.ExpectNoError(err)

			ginkgo.By("wait for the mirror pod to be running")
			gomega.Eventually(ctx, func(ctx context.Context) error {
				return checkMirrorPodRunning(ctx, f.ClientSet, mirrorPodName, ns, staticPod)
			}, 2*time.Minute, time.Second*4).Should(gomega.BeNil())
		})

		ginkgo.It("mirror pod termination should satisfy grace period when static pod is deleted [NodeConformance]", func(ctx context.Context) {
			ginkgo.By("get mirror pod uid")
			pod, err := f.ClientSet.CoreV1().Pods(ns).Get(ctx, mirrorPodName, metav1.GetOptions{})
			framework.ExpectNoError(err)
			uid := pod.UID

			ginkgo.By("delete the static pod")
			file := staticPodPath(podPath, staticPodName, ns)
			framework.Logf("deleting static pod manifest %q", file)
			err = os.Remove(file)
			framework.ExpectNoError(err)

			ginkgo.By("wait for the mirror pod to be running for grace period")
			gomega.Consistently(ctx, func(ctx context.Context) error {
				return checkMirrorPodRunningWithUID(ctx, f.ClientSet, mirrorPodName, ns, uid)
			}, 19*time.Second, 200*time.Millisecond).Should(gomega.BeNil())
		})

		ginkgo.It("mirror pod termination should satisfy grace period when static pod is updated [NodeConformance]", func(ctx context.Context) {
			ginkgo.By("get mirror pod uid")
			pod, err := f.ClientSet.CoreV1().Pods(ns).Get(ctx, mirrorPodName, metav1.GetOptions{})
			framework.ExpectNoError(err)
			uid := pod.UID

			ginkgo.By("update the static pod container image")
			image := imageutils.GetPauseImageName()
			staticPod = basicStaticPod(podPath, staticPodName, ns, image, v1.RestartPolicyAlways)
			err = writeStaticPodManifest(podPath, staticPodName, ns, staticPod)
			framework.ExpectNoError(err)

			ginkgo.By("wait for the mirror pod to be running for grace period")
			gomega.Consistently(ctx, func(ctx context.Context) error {
				return checkMirrorPodRunningWithUID(ctx, f.ClientSet, mirrorPodName, ns, uid)
			}, 19*time.Second, 200*time.Millisecond).Should(gomega.BeNil())

			ginkgo.By("wait for the mirror pod to be updated")
			gomega.Eventually(ctx, func(ctx context.Context) error {
				return checkMirrorPodRecreatedAndRunning(ctx, f.ClientSet, mirrorPodName, ns, uid, staticPod)
			}, 2*time.Minute, time.Second*4).Should(gomega.BeNil())

			ginkgo.By("check the mirror pod container image is updated")
			pod, err = f.ClientSet.CoreV1().Pods(ns).Get(ctx, mirrorPodName, metav1.GetOptions{})
			framework.ExpectNoError(err)
			framework.ExpectEqual(len(pod.Spec.Containers), 1)
			framework.ExpectEqual(pod.Spec.Containers[0].Image, image)
		})

		ginkgo.It("should update a static pod when the static pod is updated multiple times during the graceful termination period [NodeConformance]", func(ctx context.Context) {
			ginkgo.By("get mirror pod uid")
			pod, err := f.ClientSet.CoreV1().Pods(ns).Get(ctx, mirrorPodName, metav1.GetOptions{})
			framework.ExpectNoError(err)
			uid := pod.UID

			ginkgo.By("update the pod manifest multiple times during the graceful termination period")
			for i := 0; i < 300; i++ {
				staticPod = basicStaticPod(podPath, staticPodName, ns,
					fmt.Sprintf("image-%d", i), v1.RestartPolicyAlways)
				err = writeStaticPodManifest(podPath, staticPodName, ns, staticPod)
				framework.ExpectNoError(err)
				time.Sleep(100 * time.Millisecond)
			}
			image := imageutils.GetPauseImageName()
			staticPod = basicStaticPod(podPath, staticPodName, ns, image, v1.RestartPolicyAlways)
			err = writeStaticPodManifest(podPath, staticPodName, ns, staticPod)
			framework.ExpectNoError(err)

			ginkgo.By("wait for the mirror pod to be updated")
			gomega.Eventually(ctx, func(ctx context.Context) error {
				return checkMirrorPodRecreatedAndRunning(ctx, f.ClientSet, mirrorPodName, ns, uid, staticPod)
			}, 2*time.Minute, time.Second*4).Should(gomega.BeNil())

			ginkgo.By("check the mirror pod container image is updated")
			pod, err = f.ClientSet.CoreV1().Pods(ns).Get(ctx, mirrorPodName, metav1.GetOptions{})
			framework.ExpectNoError(err)
			framework.ExpectEqual(len(pod.Spec.Containers), 1)
			framework.ExpectEqual(pod.Spec.Containers[0].Image, image)
		})

		ginkgo.Context("and the container runtime is temporarily down during pod termination [NodeConformance] [Serial] [Disruptive]", func() {
			ginkgo.It("the mirror pod should terminate successfully", func(ctx context.Context) {
				ginkgo.By("verifying the pod is described as syncing in metrics")
				gomega.Eventually(ctx, getKubeletMetrics, 5*time.Second, time.Second).Should(gstruct.MatchKeys(gstruct.IgnoreExtras, gstruct.Keys{
					"kubelet_working_pods": gstruct.MatchElements(sampleLabelID, 0, gstruct.Elements{
						`kubelet_working_pods{config="desired", lifecycle="sync", static=""}`:                    timelessSample(0),
						`kubelet_working_pods{config="desired", lifecycle="sync", static="true"}`:                timelessSample(1),
						`kubelet_working_pods{config="orphan", lifecycle="sync", static=""}`:                     timelessSample(0),
						`kubelet_working_pods{config="orphan", lifecycle="sync", static="true"}`:                 timelessSample(0),
						`kubelet_working_pods{config="runtime_only", lifecycle="sync", static="unknown"}`:        timelessSample(0),
						`kubelet_working_pods{config="desired", lifecycle="terminating", static=""}`:             timelessSample(0),
						`kubelet_working_pods{config="desired", lifecycle="terminating", static="true"}`:         timelessSample(0),
						`kubelet_working_pods{config="orphan", lifecycle="terminating", static=""}`:              timelessSample(0),
						`kubelet_working_pods{config="orphan", lifecycle="terminating", static="true"}`:          timelessSample(0),
						`kubelet_working_pods{config="runtime_only", lifecycle="terminating", static="unknown"}`: timelessSample(0),
						`kubelet_working_pods{config="desired", lifecycle="terminated", static=""}`:              timelessSample(0),
						`kubelet_working_pods{config="desired", lifecycle="terminated", static="true"}`:          timelessSample(0),
						`kubelet_working_pods{config="orphan", lifecycle="terminated", static=""}`:               timelessSample(0),
						`kubelet_working_pods{config="orphan", lifecycle="terminated", static="true"}`:           timelessSample(0),
						`kubelet_working_pods{config="runtime_only", lifecycle="terminated", static="unknown"}`:  timelessSample(0),
					}),
					"kubelet_mirror_pods": gstruct.MatchElements(sampleLabelID, 0, gstruct.Elements{
						`kubelet_mirror_pods`: timelessSample(1),
					}),
					"kubelet_active_pods": gstruct.MatchElements(sampleLabelID, 0, gstruct.Elements{
						`kubelet_active_pods{static=""}`:     timelessSample(0),
						`kubelet_active_pods{static="true"}`: timelessSample(1),
					}),
					"kubelet_desired_pods": gstruct.MatchElements(sampleLabelID, 0, gstruct.Elements{
						`kubelet_desired_pods{static=""}`:     timelessSample(0),
						`kubelet_desired_pods{static="true"}`: timelessSample(1),
					}),
				}))

				ginkgo.By("delete the static pod")
				err := deleteStaticPod(podPath, staticPodName, ns)
				framework.ExpectNoError(err)

				// Note it is important we have a small delay here as we would like to reproduce https://issues.k8s.io/113091 which requires a failure in syncTerminatingPod()
				// This requires waiting a small period between the static pod being deleted so that syncTerminatingPod() will attempt to run
				ginkgo.By("sleeping before stopping the container runtime")
				time.Sleep(2 * time.Second)

				ginkgo.By("stop the container runtime")
				err = stopContainerRuntime()
				framework.ExpectNoError(err, "expected no error stopping the container runtime")

				ginkgo.By("waiting for the container runtime to be stopped")
				gomega.Eventually(ctx, func(ctx context.Context) error {
					_, _, err := getCRIClient()
					return err
				}, 2*time.Minute, time.Second*5).ShouldNot(gomega.Succeed())

				ginkgo.By("verifying the mirror pod is running")
				gomega.Consistently(ctx, func(ctx context.Context) error {
					return checkMirrorPodRunning(ctx, f.ClientSet, mirrorPodName, ns, staticPod)
				}, 19*time.Second, 200*time.Millisecond).Should(gomega.BeNil())

				ginkgo.By("verifying the pod is described as terminating in metrics")
				gomega.Eventually(ctx, getKubeletMetrics, 5*time.Second, time.Second).Should(gstruct.MatchKeys(gstruct.IgnoreExtras, gstruct.Keys{
					"kubelet_working_pods": gstruct.MatchElements(sampleLabelID, 0, gstruct.Elements{
						`kubelet_working_pods{config="desired", lifecycle="sync", static=""}`:                    timelessSample(0),
						`kubelet_working_pods{config="desired", lifecycle="sync", static="true"}`:                timelessSample(0),
						`kubelet_working_pods{config="orphan", lifecycle="sync", static=""}`:                     timelessSample(0),
						`kubelet_working_pods{config="orphan", lifecycle="sync", static="true"}`:                 timelessSample(0),
						`kubelet_working_pods{config="runtime_only", lifecycle="sync", static="unknown"}`:        timelessSample(0),
						`kubelet_working_pods{config="desired", lifecycle="terminating", static=""}`:             timelessSample(0),
						`kubelet_working_pods{config="desired", lifecycle="terminating", static="true"}`:         timelessSample(0),
						`kubelet_working_pods{config="orphan", lifecycle="terminating", static=""}`:              timelessSample(0),
						`kubelet_working_pods{config="orphan", lifecycle="terminating", static="true"}`:          timelessSample(1),
						`kubelet_working_pods{config="runtime_only", lifecycle="terminating", static="unknown"}`: timelessSample(0),
						`kubelet_working_pods{config="desired", lifecycle="terminated", static=""}`:              timelessSample(0),
						`kubelet_working_pods{config="desired", lifecycle="terminated", static="true"}`:          timelessSample(0),
						`kubelet_working_pods{config="orphan", lifecycle="terminated", static=""}`:               timelessSample(0),
						`kubelet_working_pods{config="orphan", lifecycle="terminated", static="true"}`:           timelessSample(0),
						`kubelet_working_pods{config="runtime_only", lifecycle="terminated", static="unknown"}`:  timelessSample(0),
					}),
					"kubelet_mirror_pods": gstruct.MatchElements(sampleLabelID, 0, gstruct.Elements{
						`kubelet_mirror_pods`: timelessSample(1),
					}),
					"kubelet_active_pods": gstruct.MatchElements(sampleLabelID, 0, gstruct.Elements{
						`kubelet_active_pods{static=""}`: timelessSample(0),
						// TODO: the pod is still running and consuming resources, it should be considered in
						// admission https://github.com/kubernetes/kubernetes/issues/104824 for static pods at
						// least, which means it should be 1
						`kubelet_active_pods{static="true"}`: timelessSample(0),
					}),
					"kubelet_desired_pods": gstruct.MatchElements(sampleLabelID, 0, gstruct.Elements{
						`kubelet_desired_pods{static=""}`:     timelessSample(0),
						`kubelet_desired_pods{static="true"}`: timelessSample(0),
					})}))

				ginkgo.By("start the container runtime")
				err = startContainerRuntime()
				framework.ExpectNoError(err, "expected no error starting the container runtime")
				ginkgo.By("waiting for the container runtime to start")
				gomega.Eventually(ctx, func(ctx context.Context) error {
					r, _, err := getCRIClient()
					if err != nil {
						return fmt.Errorf("error getting CRI client: %w", err)
					}
					status, err := r.Status(ctx, true)
					if err != nil {
						return fmt.Errorf("error checking CRI status: %w", err)
					}
					framework.Logf("Runtime started: %#v", status)
					return nil
				}, 2*time.Minute, time.Second*5).Should(gomega.Succeed())

				ginkgo.By(fmt.Sprintf("verifying that the mirror pod (%s/%s) stops running after about 30s", ns, mirrorPodName))
				// from the time the container runtime starts, it should take a maximum of:
				// 20s (grace period) + 2 sync transitions * 1s + 2s between housekeeping + 3s to detect CRI up +
				//   2s overhead
				// which we calculate here as "about 30s", so we try a bit longer than that but verify that it is
				// tightly bounded by not waiting longer (we want to catch regressions to shutdown)
				time.Sleep(30 * time.Second)
				gomega.Eventually(ctx, func(ctx context.Context) error {
					return checkMirrorPodDisappear(ctx, f.ClientSet, mirrorPodName, ns)
				}, time.Second*3, time.Second).Should(gomega.Succeed())

				ginkgo.By("verifying the pod finishes terminating and is removed from metrics")
				gomega.Eventually(ctx, getKubeletMetrics, 15*time.Second, time.Second).Should(gstruct.MatchKeys(gstruct.IgnoreExtras, gstruct.Keys{
					"kubelet_working_pods": gstruct.MatchElements(sampleLabelID, 0, gstruct.Elements{
						`kubelet_working_pods{config="desired", lifecycle="sync", static=""}`:                    timelessSample(0),
						`kubelet_working_pods{config="desired", lifecycle="sync", static="true"}`:                timelessSample(0),
						`kubelet_working_pods{config="orphan", lifecycle="sync", static=""}`:                     timelessSample(0),
						`kubelet_working_pods{config="orphan", lifecycle="sync", static="true"}`:                 timelessSample(0),
						`kubelet_working_pods{config="runtime_only", lifecycle="sync", static="unknown"}`:        timelessSample(0),
						`kubelet_working_pods{config="desired", lifecycle="terminating", static=""}`:             timelessSample(0),
						`kubelet_working_pods{config="desired", lifecycle="terminating", static="true"}`:         timelessSample(0),
						`kubelet_working_pods{config="orphan", lifecycle="terminating", static=""}`:              timelessSample(0),
						`kubelet_working_pods{config="orphan", lifecycle="terminating", static="true"}`:          timelessSample(0),
						`kubelet_working_pods{config="runtime_only", lifecycle="terminating", static="unknown"}`: timelessSample(0),
						`kubelet_working_pods{config="desired", lifecycle="terminated", static=""}`:              timelessSample(0),
						`kubelet_working_pods{config="desired", lifecycle="terminated", static="true"}`:          timelessSample(0),
						`kubelet_working_pods{config="orphan", lifecycle="terminated", static=""}`:               timelessSample(0),
						`kubelet_working_pods{config="orphan", lifecycle="terminated", static="true"}`:           timelessSample(0),
						`kubelet_working_pods{config="runtime_only", lifecycle="terminated", static="unknown"}`:  timelessSample(0),
					}),
					"kubelet_mirror_pods": gstruct.MatchElements(sampleLabelID, 0, gstruct.Elements{
						`kubelet_mirror_pods`: timelessSample(0),
					}),
					"kubelet_active_pods": gstruct.MatchElements(sampleLabelID, 0, gstruct.Elements{
						`kubelet_active_pods{static=""}`:     timelessSample(0),
						`kubelet_active_pods{static="true"}`: timelessSample(0),
					}),
					"kubelet_desired_pods": gstruct.MatchElements(sampleLabelID, 0, gstruct.Elements{
						`kubelet_desired_pods{static=""}`:     timelessSample(0),
						`kubelet_desired_pods{static="true"}`: timelessSample(0),
					}),
				}))
			})

			ginkgo.AfterEach(func(ctx context.Context) {
				ginkgo.By("starting the container runtime")
				err := startContainerRuntime()
				framework.ExpectNoError(err, "expected no error starting the container runtime")
				ginkgo.By("waiting for the container runtime to start")
				gomega.Eventually(ctx, func(ctx context.Context) error {
					_, _, err := getCRIClient()
					if err != nil {
						return fmt.Errorf("error getting cri client: %v", err)
					}
					return nil
				}, 2*time.Minute, time.Second*5).Should(gomega.Succeed())
			})
		})

		ginkgo.AfterEach(func(ctx context.Context) {
			ginkgo.By("delete the static pod")
			err := deleteStaticPod(podPath, staticPodName, ns)
			if !os.IsNotExist(err) {
				framework.ExpectNoError(err)
			}

			ginkgo.By("wait for the mirror pod to disappear")
			gomega.Eventually(ctx, func(ctx context.Context) error {
				return checkMirrorPodDisappear(ctx, f.ClientSet, mirrorPodName, ns)
			}, 2*time.Minute, time.Second*4).Should(gomega.BeNil())
		})
	})
})

func createStaticPodWithGracePeriod(dir, name, namespace string) *v1.Pod {
	return &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1.PodSpec{
			TerminationGracePeriodSeconds: utilpointer.Int64(20),
			Containers: []v1.Container{
				{
					Name:    "m-test",
					Image:   imageutils.GetE2EImage(imageutils.BusyBox),
					Command: []string{"sh", "-c"},
					Args: []string{`
							_term() {
								echo "Caught SIGTERM signal!"
								sleep 100
							}
							trap _term SIGTERM
							sleep 1000
							`,
					},
				},
			},
		},
	}
}

func checkMirrorPodRunningWithUID(ctx context.Context, cl clientset.Interface, name, namespace string, oUID types.UID) error {
	pod, err := cl.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("expected the mirror pod %q to appear: %w", name, err)
	}
	if pod.UID != oUID {
		return fmt.Errorf("expected the uid of mirror pod %q to be same, got %q", name, pod.UID)
	}
	if pod.Status.Phase != v1.PodRunning {
		return fmt.Errorf("expected the mirror pod %q to be running, got %q", name, pod.Status.Phase)
	}
	return nil
}

func sampleLabelID(element interface{}) string {
	el := element.(*model.Sample)
	return el.Metric.String()
}
