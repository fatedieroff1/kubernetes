// +build !ignore_autogenerated

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

// This file was autogenerated by conversion-gen. Do not edit it manually!

package v1alpha1

import (
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	api "k8s.io/kubernetes/pkg/api"
	v1 "k8s.io/kubernetes/pkg/api/v1"
	apps "k8s.io/kubernetes/pkg/apis/apps"
	unsafe "unsafe"
)

func init() {
	SchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedConversionFuncs(
		Convert_v1alpha1_PodInjectionPolicy_To_apps_PodInjectionPolicy,
		Convert_apps_PodInjectionPolicy_To_v1alpha1_PodInjectionPolicy,
		Convert_v1alpha1_PodInjectionPolicyList_To_apps_PodInjectionPolicyList,
		Convert_apps_PodInjectionPolicyList_To_v1alpha1_PodInjectionPolicyList,
		Convert_v1alpha1_PodInjectionPolicySpec_To_apps_PodInjectionPolicySpec,
		Convert_apps_PodInjectionPolicySpec_To_v1alpha1_PodInjectionPolicySpec,
	)
}

func autoConvert_v1alpha1_PodInjectionPolicy_To_apps_PodInjectionPolicy(in *PodInjectionPolicy, out *apps.PodInjectionPolicy, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_PodInjectionPolicySpec_To_apps_PodInjectionPolicySpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

func Convert_v1alpha1_PodInjectionPolicy_To_apps_PodInjectionPolicy(in *PodInjectionPolicy, out *apps.PodInjectionPolicy, s conversion.Scope) error {
	return autoConvert_v1alpha1_PodInjectionPolicy_To_apps_PodInjectionPolicy(in, out, s)
}

func autoConvert_apps_PodInjectionPolicy_To_v1alpha1_PodInjectionPolicy(in *apps.PodInjectionPolicy, out *PodInjectionPolicy, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_apps_PodInjectionPolicySpec_To_v1alpha1_PodInjectionPolicySpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

func Convert_apps_PodInjectionPolicy_To_v1alpha1_PodInjectionPolicy(in *apps.PodInjectionPolicy, out *PodInjectionPolicy, s conversion.Scope) error {
	return autoConvert_apps_PodInjectionPolicy_To_v1alpha1_PodInjectionPolicy(in, out, s)
}

func autoConvert_v1alpha1_PodInjectionPolicyList_To_apps_PodInjectionPolicyList(in *PodInjectionPolicyList, out *apps.PodInjectionPolicyList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]apps.PodInjectionPolicy, len(*in))
		for i := range *in {
			if err := Convert_v1alpha1_PodInjectionPolicy_To_apps_PodInjectionPolicy(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func Convert_v1alpha1_PodInjectionPolicyList_To_apps_PodInjectionPolicyList(in *PodInjectionPolicyList, out *apps.PodInjectionPolicyList, s conversion.Scope) error {
	return autoConvert_v1alpha1_PodInjectionPolicyList_To_apps_PodInjectionPolicyList(in, out, s)
}

func autoConvert_apps_PodInjectionPolicyList_To_v1alpha1_PodInjectionPolicyList(in *apps.PodInjectionPolicyList, out *PodInjectionPolicyList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PodInjectionPolicy, len(*in))
		for i := range *in {
			if err := Convert_apps_PodInjectionPolicy_To_v1alpha1_PodInjectionPolicy(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

func Convert_apps_PodInjectionPolicyList_To_v1alpha1_PodInjectionPolicyList(in *apps.PodInjectionPolicyList, out *PodInjectionPolicyList, s conversion.Scope) error {
	return autoConvert_apps_PodInjectionPolicyList_To_v1alpha1_PodInjectionPolicyList(in, out, s)
}

func autoConvert_v1alpha1_PodInjectionPolicySpec_To_apps_PodInjectionPolicySpec(in *PodInjectionPolicySpec, out *apps.PodInjectionPolicySpec, s conversion.Scope) error {
	out.Selector = in.Selector
	out.Env = *(*[]api.EnvVar)(unsafe.Pointer(&in.Env))
	out.EnvFrom = *(*[]api.EnvFromSource)(unsafe.Pointer(&in.EnvFrom))
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]api.Volume, len(*in))
		for i := range *in {
			// TODO: Inefficient conversion - can we improve it?
			if err := s.Convert(&(*in)[i], &(*out)[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Volumes = nil
	}
	out.VolumeMounts = *(*[]api.VolumeMount)(unsafe.Pointer(&in.VolumeMounts))
	return nil
}

func Convert_v1alpha1_PodInjectionPolicySpec_To_apps_PodInjectionPolicySpec(in *PodInjectionPolicySpec, out *apps.PodInjectionPolicySpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_PodInjectionPolicySpec_To_apps_PodInjectionPolicySpec(in, out, s)
}

func autoConvert_apps_PodInjectionPolicySpec_To_v1alpha1_PodInjectionPolicySpec(in *apps.PodInjectionPolicySpec, out *PodInjectionPolicySpec, s conversion.Scope) error {
	out.Selector = in.Selector
	out.Env = *(*[]v1.EnvVar)(unsafe.Pointer(&in.Env))
	out.EnvFrom = *(*[]v1.EnvFromSource)(unsafe.Pointer(&in.EnvFrom))
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]v1.Volume, len(*in))
		for i := range *in {
			// TODO: Inefficient conversion - can we improve it?
			if err := s.Convert(&(*in)[i], &(*out)[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.Volumes = nil
	}
	out.VolumeMounts = *(*[]v1.VolumeMount)(unsafe.Pointer(&in.VolumeMounts))
	return nil
}

func Convert_apps_PodInjectionPolicySpec_To_v1alpha1_PodInjectionPolicySpec(in *apps.PodInjectionPolicySpec, out *PodInjectionPolicySpec, s conversion.Scope) error {
	return autoConvert_apps_PodInjectionPolicySpec_To_v1alpha1_PodInjectionPolicySpec(in, out, s)
}
