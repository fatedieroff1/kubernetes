// +build !ignore_autogenerated

/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

// This file was autogenerated by deepcopy-gen. Do not edit it manually!

package extensions

import (
	api "k8s.io/kubernetes/pkg/api"
	resource "k8s.io/kubernetes/pkg/api/resource"
	unversioned "k8s.io/kubernetes/pkg/api/unversioned"
	conversion "k8s.io/kubernetes/pkg/conversion"
	runtime "k8s.io/kubernetes/pkg/runtime"
	intstr "k8s.io/kubernetes/pkg/util/intstr"
)

func init() {
	if err := api.Scheme.AddGeneratedDeepCopyFuncs(
		DeepCopy_extensions_APIVersion,
		DeepCopy_extensions_CPUTargetUtilization,
		DeepCopy_extensions_CustomMetricCurrentStatus,
		DeepCopy_extensions_CustomMetricCurrentStatusList,
		DeepCopy_extensions_CustomMetricTarget,
		DeepCopy_extensions_CustomMetricTargetList,
		DeepCopy_extensions_DaemonSet,
		DeepCopy_extensions_DaemonSetList,
		DeepCopy_extensions_DaemonSetSpec,
		DeepCopy_extensions_DaemonSetStatus,
		DeepCopy_extensions_Deployment,
		DeepCopy_extensions_DeploymentList,
		DeepCopy_extensions_DeploymentRollback,
		DeepCopy_extensions_DeploymentSpec,
		DeepCopy_extensions_DeploymentStatus,
		DeepCopy_extensions_DeploymentStrategy,
		DeepCopy_extensions_HTTPIngressPath,
		DeepCopy_extensions_HTTPIngressRuleValue,
		DeepCopy_extensions_HorizontalPodAutoscaler,
		DeepCopy_extensions_HorizontalPodAutoscalerList,
		DeepCopy_extensions_HorizontalPodAutoscalerSpec,
		DeepCopy_extensions_HorizontalPodAutoscalerStatus,
		DeepCopy_extensions_HostPortRange,
		DeepCopy_extensions_IDRange,
		DeepCopy_extensions_Ingress,
		DeepCopy_extensions_IngressBackend,
		DeepCopy_extensions_IngressList,
		DeepCopy_extensions_IngressRule,
		DeepCopy_extensions_IngressRuleValue,
		DeepCopy_extensions_IngressSpec,
		DeepCopy_extensions_IngressStatus,
		DeepCopy_extensions_IngressTLS,
		DeepCopy_extensions_Job,
		DeepCopy_extensions_JobCondition,
		DeepCopy_extensions_JobList,
		DeepCopy_extensions_JobSpec,
		DeepCopy_extensions_JobStatus,
		DeepCopy_extensions_Parameter,
		DeepCopy_extensions_PodSecurityPolicy,
		DeepCopy_extensions_PodSecurityPolicyList,
		DeepCopy_extensions_PodSecurityPolicySpec,
		DeepCopy_extensions_ReplicaSet,
		DeepCopy_extensions_ReplicaSetList,
		DeepCopy_extensions_ReplicaSetSpec,
		DeepCopy_extensions_ReplicaSetStatus,
		DeepCopy_extensions_ReplicationControllerDummy,
		DeepCopy_extensions_RollbackConfig,
		DeepCopy_extensions_RollingUpdateDeployment,
		DeepCopy_extensions_RunAsUserStrategyOptions,
		DeepCopy_extensions_SELinuxStrategyOptions,
		DeepCopy_extensions_Scale,
		DeepCopy_extensions_ScaleSpec,
		DeepCopy_extensions_ScaleStatus,
		DeepCopy_extensions_SubresourceReference,
		DeepCopy_extensions_Template,
		DeepCopy_extensions_TemplateList,
		DeepCopy_extensions_ThirdPartyResource,
		DeepCopy_extensions_ThirdPartyResourceData,
		DeepCopy_extensions_ThirdPartyResourceDataList,
		DeepCopy_extensions_ThirdPartyResourceList,
	); err != nil {
		// if one of the deep copy functions is malformed, detect it immediately.
		panic(err)
	}
}

func DeepCopy_extensions_APIVersion(in APIVersion, out *APIVersion, c *conversion.Cloner) error {
	out.Name = in.Name
	return nil
}

func DeepCopy_extensions_CPUTargetUtilization(in CPUTargetUtilization, out *CPUTargetUtilization, c *conversion.Cloner) error {
	out.TargetPercentage = in.TargetPercentage
	return nil
}

func DeepCopy_extensions_CustomMetricCurrentStatus(in CustomMetricCurrentStatus, out *CustomMetricCurrentStatus, c *conversion.Cloner) error {
	out.Name = in.Name
	if err := resource.DeepCopy_resource_Quantity(in.CurrentValue, &out.CurrentValue, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_CustomMetricCurrentStatusList(in CustomMetricCurrentStatusList, out *CustomMetricCurrentStatusList, c *conversion.Cloner) error {
	if in.Items != nil {
		in, out := in.Items, &out.Items
		*out = make([]CustomMetricCurrentStatus, len(in))
		for i := range in {
			if err := DeepCopy_extensions_CustomMetricCurrentStatus(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func DeepCopy_extensions_CustomMetricTarget(in CustomMetricTarget, out *CustomMetricTarget, c *conversion.Cloner) error {
	out.Name = in.Name
	if err := resource.DeepCopy_resource_Quantity(in.TargetValue, &out.TargetValue, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_CustomMetricTargetList(in CustomMetricTargetList, out *CustomMetricTargetList, c *conversion.Cloner) error {
	if in.Items != nil {
		in, out := in.Items, &out.Items
		*out = make([]CustomMetricTarget, len(in))
		for i := range in {
			if err := DeepCopy_extensions_CustomMetricTarget(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func DeepCopy_extensions_DaemonSet(in DaemonSet, out *DaemonSet, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := api.DeepCopy_api_ObjectMeta(in.ObjectMeta, &out.ObjectMeta, c); err != nil {
		return err
	}
	if err := DeepCopy_extensions_DaemonSetSpec(in.Spec, &out.Spec, c); err != nil {
		return err
	}
	if err := DeepCopy_extensions_DaemonSetStatus(in.Status, &out.Status, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_DaemonSetList(in DaemonSetList, out *DaemonSetList, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := unversioned.DeepCopy_unversioned_ListMeta(in.ListMeta, &out.ListMeta, c); err != nil {
		return err
	}
	if in.Items != nil {
		in, out := in.Items, &out.Items
		*out = make([]DaemonSet, len(in))
		for i := range in {
			if err := DeepCopy_extensions_DaemonSet(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func DeepCopy_extensions_DaemonSetSpec(in DaemonSetSpec, out *DaemonSetSpec, c *conversion.Cloner) error {
	if in.Selector != nil {
		in, out := in.Selector, &out.Selector
		*out = new(unversioned.LabelSelector)
		if err := unversioned.DeepCopy_unversioned_LabelSelector(*in, *out, c); err != nil {
			return err
		}
	} else {
		out.Selector = nil
	}
	if err := api.DeepCopy_api_PodTemplateSpec(in.Template, &out.Template, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_DaemonSetStatus(in DaemonSetStatus, out *DaemonSetStatus, c *conversion.Cloner) error {
	out.CurrentNumberScheduled = in.CurrentNumberScheduled
	out.NumberMisscheduled = in.NumberMisscheduled
	out.DesiredNumberScheduled = in.DesiredNumberScheduled
	return nil
}

func DeepCopy_extensions_Deployment(in Deployment, out *Deployment, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := api.DeepCopy_api_ObjectMeta(in.ObjectMeta, &out.ObjectMeta, c); err != nil {
		return err
	}
	if err := DeepCopy_extensions_DeploymentSpec(in.Spec, &out.Spec, c); err != nil {
		return err
	}
	if err := DeepCopy_extensions_DeploymentStatus(in.Status, &out.Status, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_DeploymentList(in DeploymentList, out *DeploymentList, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := unversioned.DeepCopy_unversioned_ListMeta(in.ListMeta, &out.ListMeta, c); err != nil {
		return err
	}
	if in.Items != nil {
		in, out := in.Items, &out.Items
		*out = make([]Deployment, len(in))
		for i := range in {
			if err := DeepCopy_extensions_Deployment(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func DeepCopy_extensions_DeploymentRollback(in DeploymentRollback, out *DeploymentRollback, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	out.Name = in.Name
	if in.UpdatedAnnotations != nil {
		in, out := in.UpdatedAnnotations, &out.UpdatedAnnotations
		*out = make(map[string]string)
		for key, val := range in {
			(*out)[key] = val
		}
	} else {
		out.UpdatedAnnotations = nil
	}
	if err := DeepCopy_extensions_RollbackConfig(in.RollbackTo, &out.RollbackTo, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_DeploymentSpec(in DeploymentSpec, out *DeploymentSpec, c *conversion.Cloner) error {
	out.Replicas = in.Replicas
	if in.Selector != nil {
		in, out := in.Selector, &out.Selector
		*out = new(unversioned.LabelSelector)
		if err := unversioned.DeepCopy_unversioned_LabelSelector(*in, *out, c); err != nil {
			return err
		}
	} else {
		out.Selector = nil
	}
	if err := api.DeepCopy_api_PodTemplateSpec(in.Template, &out.Template, c); err != nil {
		return err
	}
	if err := DeepCopy_extensions_DeploymentStrategy(in.Strategy, &out.Strategy, c); err != nil {
		return err
	}
	out.MinReadySeconds = in.MinReadySeconds
	if in.RevisionHistoryLimit != nil {
		in, out := in.RevisionHistoryLimit, &out.RevisionHistoryLimit
		*out = new(int)
		**out = *in
	} else {
		out.RevisionHistoryLimit = nil
	}
	out.Paused = in.Paused
	if in.RollbackTo != nil {
		in, out := in.RollbackTo, &out.RollbackTo
		*out = new(RollbackConfig)
		if err := DeepCopy_extensions_RollbackConfig(*in, *out, c); err != nil {
			return err
		}
	} else {
		out.RollbackTo = nil
	}
	return nil
}

func DeepCopy_extensions_DeploymentStatus(in DeploymentStatus, out *DeploymentStatus, c *conversion.Cloner) error {
	out.ObservedGeneration = in.ObservedGeneration
	out.Replicas = in.Replicas
	out.UpdatedReplicas = in.UpdatedReplicas
	out.AvailableReplicas = in.AvailableReplicas
	out.UnavailableReplicas = in.UnavailableReplicas
	return nil
}

func DeepCopy_extensions_DeploymentStrategy(in DeploymentStrategy, out *DeploymentStrategy, c *conversion.Cloner) error {
	out.Type = in.Type
	if in.RollingUpdate != nil {
		in, out := in.RollingUpdate, &out.RollingUpdate
		*out = new(RollingUpdateDeployment)
		if err := DeepCopy_extensions_RollingUpdateDeployment(*in, *out, c); err != nil {
			return err
		}
	} else {
		out.RollingUpdate = nil
	}
	return nil
}

func DeepCopy_extensions_HTTPIngressPath(in HTTPIngressPath, out *HTTPIngressPath, c *conversion.Cloner) error {
	out.Path = in.Path
	if err := DeepCopy_extensions_IngressBackend(in.Backend, &out.Backend, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_HTTPIngressRuleValue(in HTTPIngressRuleValue, out *HTTPIngressRuleValue, c *conversion.Cloner) error {
	if in.Paths != nil {
		in, out := in.Paths, &out.Paths
		*out = make([]HTTPIngressPath, len(in))
		for i := range in {
			if err := DeepCopy_extensions_HTTPIngressPath(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.Paths = nil
	}
	return nil
}

func DeepCopy_extensions_HorizontalPodAutoscaler(in HorizontalPodAutoscaler, out *HorizontalPodAutoscaler, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := api.DeepCopy_api_ObjectMeta(in.ObjectMeta, &out.ObjectMeta, c); err != nil {
		return err
	}
	if err := DeepCopy_extensions_HorizontalPodAutoscalerSpec(in.Spec, &out.Spec, c); err != nil {
		return err
	}
	if err := DeepCopy_extensions_HorizontalPodAutoscalerStatus(in.Status, &out.Status, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_HorizontalPodAutoscalerList(in HorizontalPodAutoscalerList, out *HorizontalPodAutoscalerList, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := unversioned.DeepCopy_unversioned_ListMeta(in.ListMeta, &out.ListMeta, c); err != nil {
		return err
	}
	if in.Items != nil {
		in, out := in.Items, &out.Items
		*out = make([]HorizontalPodAutoscaler, len(in))
		for i := range in {
			if err := DeepCopy_extensions_HorizontalPodAutoscaler(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func DeepCopy_extensions_HorizontalPodAutoscalerSpec(in HorizontalPodAutoscalerSpec, out *HorizontalPodAutoscalerSpec, c *conversion.Cloner) error {
	if err := DeepCopy_extensions_SubresourceReference(in.ScaleRef, &out.ScaleRef, c); err != nil {
		return err
	}
	if in.MinReplicas != nil {
		in, out := in.MinReplicas, &out.MinReplicas
		*out = new(int)
		**out = *in
	} else {
		out.MinReplicas = nil
	}
	out.MaxReplicas = in.MaxReplicas
	if in.CPUUtilization != nil {
		in, out := in.CPUUtilization, &out.CPUUtilization
		*out = new(CPUTargetUtilization)
		if err := DeepCopy_extensions_CPUTargetUtilization(*in, *out, c); err != nil {
			return err
		}
	} else {
		out.CPUUtilization = nil
	}
	return nil
}

func DeepCopy_extensions_HorizontalPodAutoscalerStatus(in HorizontalPodAutoscalerStatus, out *HorizontalPodAutoscalerStatus, c *conversion.Cloner) error {
	if in.ObservedGeneration != nil {
		in, out := in.ObservedGeneration, &out.ObservedGeneration
		*out = new(int64)
		**out = *in
	} else {
		out.ObservedGeneration = nil
	}
	if in.LastScaleTime != nil {
		in, out := in.LastScaleTime, &out.LastScaleTime
		*out = new(unversioned.Time)
		if err := unversioned.DeepCopy_unversioned_Time(*in, *out, c); err != nil {
			return err
		}
	} else {
		out.LastScaleTime = nil
	}
	out.CurrentReplicas = in.CurrentReplicas
	out.DesiredReplicas = in.DesiredReplicas
	if in.CurrentCPUUtilizationPercentage != nil {
		in, out := in.CurrentCPUUtilizationPercentage, &out.CurrentCPUUtilizationPercentage
		*out = new(int)
		**out = *in
	} else {
		out.CurrentCPUUtilizationPercentage = nil
	}
	return nil
}

func DeepCopy_extensions_HostPortRange(in HostPortRange, out *HostPortRange, c *conversion.Cloner) error {
	out.Min = in.Min
	out.Max = in.Max
	return nil
}

func DeepCopy_extensions_IDRange(in IDRange, out *IDRange, c *conversion.Cloner) error {
	out.Min = in.Min
	out.Max = in.Max
	return nil
}

func DeepCopy_extensions_Ingress(in Ingress, out *Ingress, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := api.DeepCopy_api_ObjectMeta(in.ObjectMeta, &out.ObjectMeta, c); err != nil {
		return err
	}
	if err := DeepCopy_extensions_IngressSpec(in.Spec, &out.Spec, c); err != nil {
		return err
	}
	if err := DeepCopy_extensions_IngressStatus(in.Status, &out.Status, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_IngressBackend(in IngressBackend, out *IngressBackend, c *conversion.Cloner) error {
	out.ServiceName = in.ServiceName
	if err := intstr.DeepCopy_intstr_IntOrString(in.ServicePort, &out.ServicePort, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_IngressList(in IngressList, out *IngressList, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := unversioned.DeepCopy_unversioned_ListMeta(in.ListMeta, &out.ListMeta, c); err != nil {
		return err
	}
	if in.Items != nil {
		in, out := in.Items, &out.Items
		*out = make([]Ingress, len(in))
		for i := range in {
			if err := DeepCopy_extensions_Ingress(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func DeepCopy_extensions_IngressRule(in IngressRule, out *IngressRule, c *conversion.Cloner) error {
	out.Host = in.Host
	if err := DeepCopy_extensions_IngressRuleValue(in.IngressRuleValue, &out.IngressRuleValue, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_IngressRuleValue(in IngressRuleValue, out *IngressRuleValue, c *conversion.Cloner) error {
	if in.HTTP != nil {
		in, out := in.HTTP, &out.HTTP
		*out = new(HTTPIngressRuleValue)
		if err := DeepCopy_extensions_HTTPIngressRuleValue(*in, *out, c); err != nil {
			return err
		}
	} else {
		out.HTTP = nil
	}
	return nil
}

func DeepCopy_extensions_IngressSpec(in IngressSpec, out *IngressSpec, c *conversion.Cloner) error {
	if in.Backend != nil {
		in, out := in.Backend, &out.Backend
		*out = new(IngressBackend)
		if err := DeepCopy_extensions_IngressBackend(*in, *out, c); err != nil {
			return err
		}
	} else {
		out.Backend = nil
	}
	if in.TLS != nil {
		in, out := in.TLS, &out.TLS
		*out = make([]IngressTLS, len(in))
		for i := range in {
			if err := DeepCopy_extensions_IngressTLS(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.TLS = nil
	}
	if in.Rules != nil {
		in, out := in.Rules, &out.Rules
		*out = make([]IngressRule, len(in))
		for i := range in {
			if err := DeepCopy_extensions_IngressRule(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.Rules = nil
	}
	return nil
}

func DeepCopy_extensions_IngressStatus(in IngressStatus, out *IngressStatus, c *conversion.Cloner) error {
	if err := api.DeepCopy_api_LoadBalancerStatus(in.LoadBalancer, &out.LoadBalancer, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_IngressTLS(in IngressTLS, out *IngressTLS, c *conversion.Cloner) error {
	if in.Hosts != nil {
		in, out := in.Hosts, &out.Hosts
		*out = make([]string, len(in))
		copy(*out, in)
	} else {
		out.Hosts = nil
	}
	out.SecretName = in.SecretName
	return nil
}

func DeepCopy_extensions_Job(in Job, out *Job, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := api.DeepCopy_api_ObjectMeta(in.ObjectMeta, &out.ObjectMeta, c); err != nil {
		return err
	}
	if err := DeepCopy_extensions_JobSpec(in.Spec, &out.Spec, c); err != nil {
		return err
	}
	if err := DeepCopy_extensions_JobStatus(in.Status, &out.Status, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_JobCondition(in JobCondition, out *JobCondition, c *conversion.Cloner) error {
	out.Type = in.Type
	out.Status = in.Status
	if err := unversioned.DeepCopy_unversioned_Time(in.LastProbeTime, &out.LastProbeTime, c); err != nil {
		return err
	}
	if err := unversioned.DeepCopy_unversioned_Time(in.LastTransitionTime, &out.LastTransitionTime, c); err != nil {
		return err
	}
	out.Reason = in.Reason
	out.Message = in.Message
	return nil
}

func DeepCopy_extensions_JobList(in JobList, out *JobList, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := unversioned.DeepCopy_unversioned_ListMeta(in.ListMeta, &out.ListMeta, c); err != nil {
		return err
	}
	if in.Items != nil {
		in, out := in.Items, &out.Items
		*out = make([]Job, len(in))
		for i := range in {
			if err := DeepCopy_extensions_Job(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func DeepCopy_extensions_JobSpec(in JobSpec, out *JobSpec, c *conversion.Cloner) error {
	if in.Parallelism != nil {
		in, out := in.Parallelism, &out.Parallelism
		*out = new(int)
		**out = *in
	} else {
		out.Parallelism = nil
	}
	if in.Completions != nil {
		in, out := in.Completions, &out.Completions
		*out = new(int)
		**out = *in
	} else {
		out.Completions = nil
	}
	if in.ActiveDeadlineSeconds != nil {
		in, out := in.ActiveDeadlineSeconds, &out.ActiveDeadlineSeconds
		*out = new(int64)
		**out = *in
	} else {
		out.ActiveDeadlineSeconds = nil
	}
	if in.Selector != nil {
		in, out := in.Selector, &out.Selector
		*out = new(unversioned.LabelSelector)
		if err := unversioned.DeepCopy_unversioned_LabelSelector(*in, *out, c); err != nil {
			return err
		}
	} else {
		out.Selector = nil
	}
	if in.ManualSelector != nil {
		in, out := in.ManualSelector, &out.ManualSelector
		*out = new(bool)
		**out = *in
	} else {
		out.ManualSelector = nil
	}
	if err := api.DeepCopy_api_PodTemplateSpec(in.Template, &out.Template, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_JobStatus(in JobStatus, out *JobStatus, c *conversion.Cloner) error {
	if in.Conditions != nil {
		in, out := in.Conditions, &out.Conditions
		*out = make([]JobCondition, len(in))
		for i := range in {
			if err := DeepCopy_extensions_JobCondition(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.Conditions = nil
	}
	if in.StartTime != nil {
		in, out := in.StartTime, &out.StartTime
		*out = new(unversioned.Time)
		if err := unversioned.DeepCopy_unversioned_Time(*in, *out, c); err != nil {
			return err
		}
	} else {
		out.StartTime = nil
	}
	if in.CompletionTime != nil {
		in, out := in.CompletionTime, &out.CompletionTime
		*out = new(unversioned.Time)
		if err := unversioned.DeepCopy_unversioned_Time(*in, *out, c); err != nil {
			return err
		}
	} else {
		out.CompletionTime = nil
	}
	out.Active = in.Active
	out.Succeeded = in.Succeeded
	out.Failed = in.Failed
	return nil
}

func DeepCopy_extensions_Parameter(in Parameter, out *Parameter, c *conversion.Cloner) error {
	out.Name = in.Name
	out.DisplayName = in.DisplayName
	out.Description = in.Description
	out.Value = in.Value
	out.Generate = in.Generate
	out.From = in.From
	out.Required = in.Required
	return nil
}

func DeepCopy_extensions_PodSecurityPolicy(in PodSecurityPolicy, out *PodSecurityPolicy, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := api.DeepCopy_api_ObjectMeta(in.ObjectMeta, &out.ObjectMeta, c); err != nil {
		return err
	}
	if err := DeepCopy_extensions_PodSecurityPolicySpec(in.Spec, &out.Spec, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_PodSecurityPolicyList(in PodSecurityPolicyList, out *PodSecurityPolicyList, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := unversioned.DeepCopy_unversioned_ListMeta(in.ListMeta, &out.ListMeta, c); err != nil {
		return err
	}
	if in.Items != nil {
		in, out := in.Items, &out.Items
		*out = make([]PodSecurityPolicy, len(in))
		for i := range in {
			if err := DeepCopy_extensions_PodSecurityPolicy(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func DeepCopy_extensions_PodSecurityPolicySpec(in PodSecurityPolicySpec, out *PodSecurityPolicySpec, c *conversion.Cloner) error {
	out.Privileged = in.Privileged
	if in.Capabilities != nil {
		in, out := in.Capabilities, &out.Capabilities
		*out = make([]api.Capability, len(in))
		for i := range in {
			(*out)[i] = in[i]
		}
	} else {
		out.Capabilities = nil
	}
	if in.Volumes != nil {
		in, out := in.Volumes, &out.Volumes
		*out = make([]FSType, len(in))
		for i := range in {
			(*out)[i] = in[i]
		}
	} else {
		out.Volumes = nil
	}
	out.HostNetwork = in.HostNetwork
	if in.HostPorts != nil {
		in, out := in.HostPorts, &out.HostPorts
		*out = make([]HostPortRange, len(in))
		for i := range in {
			if err := DeepCopy_extensions_HostPortRange(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.HostPorts = nil
	}
	out.HostPID = in.HostPID
	out.HostIPC = in.HostIPC
	if err := DeepCopy_extensions_SELinuxStrategyOptions(in.SELinux, &out.SELinux, c); err != nil {
		return err
	}
	if err := DeepCopy_extensions_RunAsUserStrategyOptions(in.RunAsUser, &out.RunAsUser, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_ReplicaSet(in ReplicaSet, out *ReplicaSet, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := api.DeepCopy_api_ObjectMeta(in.ObjectMeta, &out.ObjectMeta, c); err != nil {
		return err
	}
	if err := DeepCopy_extensions_ReplicaSetSpec(in.Spec, &out.Spec, c); err != nil {
		return err
	}
	if err := DeepCopy_extensions_ReplicaSetStatus(in.Status, &out.Status, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_ReplicaSetList(in ReplicaSetList, out *ReplicaSetList, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := unversioned.DeepCopy_unversioned_ListMeta(in.ListMeta, &out.ListMeta, c); err != nil {
		return err
	}
	if in.Items != nil {
		in, out := in.Items, &out.Items
		*out = make([]ReplicaSet, len(in))
		for i := range in {
			if err := DeepCopy_extensions_ReplicaSet(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func DeepCopy_extensions_ReplicaSetSpec(in ReplicaSetSpec, out *ReplicaSetSpec, c *conversion.Cloner) error {
	out.Replicas = in.Replicas
	if in.Selector != nil {
		in, out := in.Selector, &out.Selector
		*out = new(unversioned.LabelSelector)
		if err := unversioned.DeepCopy_unversioned_LabelSelector(*in, *out, c); err != nil {
			return err
		}
	} else {
		out.Selector = nil
	}
	if err := api.DeepCopy_api_PodTemplateSpec(in.Template, &out.Template, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_ReplicaSetStatus(in ReplicaSetStatus, out *ReplicaSetStatus, c *conversion.Cloner) error {
	out.Replicas = in.Replicas
	out.FullyLabeledReplicas = in.FullyLabeledReplicas
	out.ObservedGeneration = in.ObservedGeneration
	return nil
}

func DeepCopy_extensions_ReplicationControllerDummy(in ReplicationControllerDummy, out *ReplicationControllerDummy, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_RollbackConfig(in RollbackConfig, out *RollbackConfig, c *conversion.Cloner) error {
	out.Revision = in.Revision
	return nil
}

func DeepCopy_extensions_RollingUpdateDeployment(in RollingUpdateDeployment, out *RollingUpdateDeployment, c *conversion.Cloner) error {
	if err := intstr.DeepCopy_intstr_IntOrString(in.MaxUnavailable, &out.MaxUnavailable, c); err != nil {
		return err
	}
	if err := intstr.DeepCopy_intstr_IntOrString(in.MaxSurge, &out.MaxSurge, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_RunAsUserStrategyOptions(in RunAsUserStrategyOptions, out *RunAsUserStrategyOptions, c *conversion.Cloner) error {
	out.Rule = in.Rule
	if in.Ranges != nil {
		in, out := in.Ranges, &out.Ranges
		*out = make([]IDRange, len(in))
		for i := range in {
			if err := DeepCopy_extensions_IDRange(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.Ranges = nil
	}
	return nil
}

func DeepCopy_extensions_SELinuxStrategyOptions(in SELinuxStrategyOptions, out *SELinuxStrategyOptions, c *conversion.Cloner) error {
	out.Rule = in.Rule
	if in.SELinuxOptions != nil {
		in, out := in.SELinuxOptions, &out.SELinuxOptions
		*out = new(api.SELinuxOptions)
		if err := api.DeepCopy_api_SELinuxOptions(*in, *out, c); err != nil {
			return err
		}
	} else {
		out.SELinuxOptions = nil
	}
	return nil
}

func DeepCopy_extensions_Scale(in Scale, out *Scale, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := api.DeepCopy_api_ObjectMeta(in.ObjectMeta, &out.ObjectMeta, c); err != nil {
		return err
	}
	if err := DeepCopy_extensions_ScaleSpec(in.Spec, &out.Spec, c); err != nil {
		return err
	}
	if err := DeepCopy_extensions_ScaleStatus(in.Status, &out.Status, c); err != nil {
		return err
	}
	return nil
}

func DeepCopy_extensions_ScaleSpec(in ScaleSpec, out *ScaleSpec, c *conversion.Cloner) error {
	out.Replicas = in.Replicas
	return nil
}

func DeepCopy_extensions_ScaleStatus(in ScaleStatus, out *ScaleStatus, c *conversion.Cloner) error {
	out.Replicas = in.Replicas
	if in.Selector != nil {
		in, out := in.Selector, &out.Selector
		*out = new(unversioned.LabelSelector)
		if err := unversioned.DeepCopy_unversioned_LabelSelector(*in, *out, c); err != nil {
			return err
		}
	} else {
		out.Selector = nil
	}
	return nil
}

func DeepCopy_extensions_SubresourceReference(in SubresourceReference, out *SubresourceReference, c *conversion.Cloner) error {
	out.Kind = in.Kind
	out.Name = in.Name
	out.APIVersion = in.APIVersion
	out.Subresource = in.Subresource
	return nil
}

func DeepCopy_extensions_Template(in Template, out *Template, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := api.DeepCopy_api_ObjectMeta(in.ObjectMeta, &out.ObjectMeta, c); err != nil {
		return err
	}
	if in.Parameters != nil {
		in, out := in.Parameters, &out.Parameters
		*out = make([]Parameter, len(in))
		for i := range in {
			if err := DeepCopy_extensions_Parameter(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.Parameters = nil
	}
	if in.Objects != nil {
		in, out := in.Objects, &out.Objects
		*out = make([]runtime.Object, len(in))
		for i := range in {
			if newVal, err := c.DeepCopy(in[i]); err != nil {
				return err
			} else {
				(*out)[i] = newVal.(runtime.Object)
			}
		}
	} else {
		out.Objects = nil
	}
	if in.ObjectLabels != nil {
		in, out := in.ObjectLabels, &out.ObjectLabels
		*out = make(map[string]string)
		for key, val := range in {
			(*out)[key] = val
		}
	} else {
		out.ObjectLabels = nil
	}
	return nil
}

func DeepCopy_extensions_TemplateList(in TemplateList, out *TemplateList, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := unversioned.DeepCopy_unversioned_ListMeta(in.ListMeta, &out.ListMeta, c); err != nil {
		return err
	}
	if in.Items != nil {
		in, out := in.Items, &out.Items
		*out = make([]Template, len(in))
		for i := range in {
			if err := DeepCopy_extensions_Template(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func DeepCopy_extensions_ThirdPartyResource(in ThirdPartyResource, out *ThirdPartyResource, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := api.DeepCopy_api_ObjectMeta(in.ObjectMeta, &out.ObjectMeta, c); err != nil {
		return err
	}
	out.Description = in.Description
	if in.Versions != nil {
		in, out := in.Versions, &out.Versions
		*out = make([]APIVersion, len(in))
		for i := range in {
			if err := DeepCopy_extensions_APIVersion(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.Versions = nil
	}
	return nil
}

func DeepCopy_extensions_ThirdPartyResourceData(in ThirdPartyResourceData, out *ThirdPartyResourceData, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := api.DeepCopy_api_ObjectMeta(in.ObjectMeta, &out.ObjectMeta, c); err != nil {
		return err
	}
	if in.Data != nil {
		in, out := in.Data, &out.Data
		*out = make([]byte, len(in))
		copy(*out, in)
	} else {
		out.Data = nil
	}
	return nil
}

func DeepCopy_extensions_ThirdPartyResourceDataList(in ThirdPartyResourceDataList, out *ThirdPartyResourceDataList, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := unversioned.DeepCopy_unversioned_ListMeta(in.ListMeta, &out.ListMeta, c); err != nil {
		return err
	}
	if in.Items != nil {
		in, out := in.Items, &out.Items
		*out = make([]ThirdPartyResourceData, len(in))
		for i := range in {
			if err := DeepCopy_extensions_ThirdPartyResourceData(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func DeepCopy_extensions_ThirdPartyResourceList(in ThirdPartyResourceList, out *ThirdPartyResourceList, c *conversion.Cloner) error {
	if err := unversioned.DeepCopy_unversioned_TypeMeta(in.TypeMeta, &out.TypeMeta, c); err != nil {
		return err
	}
	if err := unversioned.DeepCopy_unversioned_ListMeta(in.ListMeta, &out.ListMeta, c); err != nil {
		return err
	}
	if in.Items != nil {
		in, out := in.Items, &out.Items
		*out = make([]ThirdPartyResource, len(in))
		for i := range in {
			if err := DeepCopy_extensions_ThirdPartyResource(in[i], &(*out)[i], c); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}
