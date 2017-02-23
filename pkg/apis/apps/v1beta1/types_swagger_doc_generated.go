/*
Copyright 2016 The Kubernetes Authors.

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

package v1beta1

// This file contains a collection of methods that can be used from go-restful to
// generate Swagger API documentation for its models. Please read this PR for more
// information on the implementation: https://github.com/emicklei/go-restful/pull/215
//
// TODOs are ignored from the parser (e.g. TODO(andronat):... || TODO:...) if and only if
// they are on one line! For multiple line or blocks that you want to ignore use ---.
// Any context after a --- is ignored.
//
// Those methods can be generated by using hack/update-generated-swagger-docs.sh

// AUTO-GENERATED FUNCTIONS START HERE
var map_Deployment = map[string]string{
	"":         "Deployment enables declarative updates for Pods and ReplicaSets.",
	"metadata": "Standard object metadata.",
	"spec":     "Specification of the desired behavior of the Deployment.",
	"status":   "Most recently observed status of the Deployment.",
}

func (Deployment) SwaggerDoc() map[string]string {
	return map_Deployment
}

var map_DeploymentCondition = map[string]string{
	"":                   "DeploymentCondition describes the state of a deployment at a certain point.",
	"type":               "Type of deployment condition.",
	"status":             "Status of the condition, one of True, False, Unknown.",
	"lastUpdateTime":     "The last time this condition was updated.",
	"lastTransitionTime": "Last time the condition transitioned from one status to another.",
	"reason":             "The reason for the condition's last transition.",
	"message":            "A human readable message indicating details about the transition.",
}

func (DeploymentCondition) SwaggerDoc() map[string]string {
	return map_DeploymentCondition
}

var map_DeploymentList = map[string]string{
	"":         "DeploymentList is a list of Deployments.",
	"metadata": "Standard list metadata.",
	"items":    "Items is the list of Deployments.",
}

func (DeploymentList) SwaggerDoc() map[string]string {
	return map_DeploymentList
}

var map_DeploymentRollback = map[string]string{
	"":                   "DeploymentRollback stores the information required to rollback a deployment.",
	"name":               "Required: This must match the Name of a deployment.",
	"updatedAnnotations": "The annotations to be updated to a deployment",
	"rollbackTo":         "The config of this deployment rollback.",
}

func (DeploymentRollback) SwaggerDoc() map[string]string {
	return map_DeploymentRollback
}

var map_DeploymentSpec = map[string]string{
	"":                        "DeploymentSpec is the specification of the desired behavior of the Deployment.",
	"replicas":                "Number of desired pods. This is a pointer to distinguish between explicit zero and not specified. Defaults to 1.",
	"selector":                "Label selector for pods. Existing ReplicaSets whose pods are selected by this will be the ones affected by this deployment.",
	"template":                "Template describes the pods that will be created.",
	"strategy":                "The deployment strategy to use to replace existing pods with new ones.",
	"minReadySeconds":         "Minimum number of seconds for which a newly created pod should be ready without any of its container crashing, for it to be considered available. Defaults to 0 (pod will be considered available as soon as it is ready)",
	"revisionHistoryLimit":    "The number of old ReplicaSets to retain to allow rollback. This is a pointer to distinguish between explicit zero and not specified.",
	"paused":                  "Indicates that the deployment is paused.",
	"rollbackTo":              "The config this deployment is rolling back to. Will be cleared after rollback is done.",
	"progressDeadlineSeconds": "The maximum time in seconds for a deployment to make progress before it is considered to be failed. The deployment controller will continue to process failed deployments and a condition with a ProgressDeadlineExceeded reason will be surfaced in the deployment status. Once autoRollback is implemented, the deployment controller will automatically rollback failed deployments. Note that progress will not be estimated during the time a deployment is paused. This is not set by default.",
}

func (DeploymentSpec) SwaggerDoc() map[string]string {
	return map_DeploymentSpec
}

var map_DeploymentStatus = map[string]string{
	"":                    "DeploymentStatus is the most recently observed status of the Deployment.",
	"observedGeneration":  "The generation observed by the deployment controller.",
	"replicas":            "Total number of non-terminated pods targeted by this deployment (their labels match the selector).",
	"updatedReplicas":     "Total number of non-terminated pods targeted by this deployment that have the desired template spec.",
	"readyReplicas":       "Total number of ready pods targeted by this deployment.",
	"availableReplicas":   "Total number of available pods (ready for at least minReadySeconds) targeted by this deployment.",
	"unavailableReplicas": "Total number of unavailable pods targeted by this deployment.",
	"conditions":          "Represents the latest available observations of a deployment's current state.",
}

func (DeploymentStatus) SwaggerDoc() map[string]string {
	return map_DeploymentStatus
}

var map_DeploymentStrategy = map[string]string{
	"":              "DeploymentStrategy describes how to replace existing pods with new ones.",
	"type":          "Type of deployment. Can be \"Recreate\" or \"RollingUpdate\". Default is RollingUpdate.",
	"rollingUpdate": "Rolling update config params. Present only if DeploymentStrategyType = RollingUpdate.",
}

func (DeploymentStrategy) SwaggerDoc() map[string]string {
	return map_DeploymentStrategy
}

var map_RollbackConfig = map[string]string{
	"revision": "The revision to rollback to. If set to 0, rollbck to the last revision.",
}

func (RollbackConfig) SwaggerDoc() map[string]string {
	return map_RollbackConfig
}

var map_RollingUpdateDeployment = map[string]string{
	"":               "Spec to control the desired behavior of rolling update.",
	"maxUnavailable": "The maximum number of pods that can be unavailable during the update. Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%). Absolute number is calculated from percentage by rounding down. This can not be 0 if MaxSurge is 0. By default, a fixed value of 1 is used. Example: when this is set to 30%, the old RC can be scaled down to 70% of desired pods immediately when the rolling update starts. Once new pods are ready, old RC can be scaled down further, followed by scaling up the new RC, ensuring that the total number of pods available at all times during the update is at least 70% of desired pods.",
	"maxSurge":       "The maximum number of pods that can be scheduled above the desired number of pods. Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%). This can not be 0 if MaxUnavailable is 0. Absolute number is calculated from percentage by rounding up. By default, a value of 1 is used. Example: when this is set to 30%, the new RC can be scaled up immediately when the rolling update starts, such that the total number of old and new pods do not exceed 130% of desired pods. Once old pods have been killed, new RC can be scaled up further, ensuring that total number of pods running at any time during the update is atmost 130% of desired pods.",
}

func (RollingUpdateDeployment) SwaggerDoc() map[string]string {
	return map_RollingUpdateDeployment
}

var map_Scale = map[string]string{
	"":         "Scale represents a scaling request for a resource.",
	"metadata": "Standard object metadata; More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata.",
	"spec":     "defines the behavior of the scale. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#spec-and-status.",
	"status":   "current status of the scale. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#spec-and-status. Read-only.",
}

func (Scale) SwaggerDoc() map[string]string {
	return map_Scale
}

var map_ScaleSpec = map[string]string{
	"":         "ScaleSpec describes the attributes of a scale subresource",
	"replicas": "desired number of instances for the scaled object.",
}

func (ScaleSpec) SwaggerDoc() map[string]string {
	return map_ScaleSpec
}

var map_ScaleStatus = map[string]string{
	"":               "ScaleStatus represents the current status of a scale subresource.",
	"replicas":       "actual number of observed instances of the scaled object.",
	"selector":       "label query over pods that should match the replicas count. More info: http://kubernetes.io/docs/user-guide/labels#label-selectors",
	"targetSelector": "label selector for pods that should match the replicas count. This is a serializated version of both map-based and more expressive set-based selectors. This is done to avoid introspection in the clients. The string will be in the same format as the query-param syntax. If the target type only supports map-based selectors, both this field and map-based selector field are populated. More info: http://kubernetes.io/docs/user-guide/labels#label-selectors",
}

func (ScaleStatus) SwaggerDoc() map[string]string {
	return map_ScaleStatus
}

var map_StatefulSet = map[string]string{
	"":       "StatefulSet represents a set of pods with consistent identities. Identities are defined as:\n - Network: A single stable DNS and hostname.\n - Storage: As many VolumeClaims as requested.\nThe StatefulSet guarantees that a given network identity will always map to the same storage identity.",
	"spec":   "Spec defines the desired identities of pods in this set.",
	"status": "Status is the current status of Pods in this StatefulSet. This data may be out of date by some window of time.",
}

func (StatefulSet) SwaggerDoc() map[string]string {
	return map_StatefulSet
}

var map_StatefulSetList = map[string]string{
	"": "StatefulSetList is a collection of StatefulSets.",
}

func (StatefulSetList) SwaggerDoc() map[string]string {
	return map_StatefulSetList
}

var map_StatefulSetSpec = map[string]string{
	"":                     "A StatefulSetSpec is the specification of a StatefulSet.",
	"replicas":             "Replicas is the desired number of replicas of the given Template. These are replicas in the sense that they are instantiations of the same Template, but individual replicas also have a consistent identity. If unspecified, defaults to 1.",
	"selector":             "Selector is a label query over pods that should match the replica count. If empty, defaulted to labels on the pod template. More info: http://kubernetes.io/docs/user-guide/labels#label-selectors",
	"template":             "Template is the object that describes the pod that will be created if insufficient replicas are detected. Each pod stamped out by the StatefulSet will fulfill this Template, but have a unique identity from the rest of the StatefulSet.",
	"volumeClaimTemplates": "VolumeClaimTemplates is a list of claims that pods are allowed to reference. The StatefulSet controller is responsible for mapping network identities to claims in a way that maintains the identity of a pod. Every claim in this list must have at least one matching (by name) volumeMount in one container in the template. A claim in this list takes precedence over any volumes in the template, with the same name.",
	"serviceName":          "ServiceName is the name of the service that governs this StatefulSet. This service must exist before the StatefulSet, and is responsible for the network identity of the set. Pods get DNS/hostnames that follow the pattern: pod-specific-string.serviceName.default.svc.cluster.local where \"pod-specific-string\" is managed by the StatefulSet controller.",
}

func (StatefulSetSpec) SwaggerDoc() map[string]string {
	return map_StatefulSetSpec
}

var map_StatefulSetStatus = map[string]string{
	"":                   "StatefulSetStatus represents the current state of a StatefulSet.",
	"observedGeneration": "most recent generation observed by this StatefulSet.",
	"replicas":           "Replicas is the number of actual replicas.",
}

func (StatefulSetStatus) SwaggerDoc() map[string]string {
	return map_StatefulSetStatus
}

// AUTO-GENERATED FUNCTIONS END HERE
