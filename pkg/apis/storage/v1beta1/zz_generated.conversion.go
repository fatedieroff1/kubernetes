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

package v1beta1

import (
	v1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/storage/v1beta1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	core "k8s.io/kubernetes/pkg/apis/core"
	storage "k8s.io/kubernetes/pkg/apis/storage"
	unsafe "unsafe"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedConversionFuncs(
		Convert_v1beta1_StorageClass_To_storage_StorageClass,
		Convert_storage_StorageClass_To_v1beta1_StorageClass,
		Convert_v1beta1_StorageClassList_To_storage_StorageClassList,
		Convert_storage_StorageClassList_To_v1beta1_StorageClassList,
	)
}

func autoConvert_v1beta1_StorageClass_To_storage_StorageClass(in *v1beta1.StorageClass, out *storage.StorageClass, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	out.Provisioner = in.Provisioner
	out.Parameters = *(*map[string]string)(unsafe.Pointer(&in.Parameters))
	out.ReclaimPolicy = (*core.PersistentVolumeReclaimPolicy)(unsafe.Pointer(in.ReclaimPolicy))
	out.MountOptions = *(*[]string)(unsafe.Pointer(&in.MountOptions))
	out.AllowVolumeExpansion = (*bool)(unsafe.Pointer(in.AllowVolumeExpansion))
	return nil
}

// Convert_v1beta1_StorageClass_To_storage_StorageClass is an autogenerated conversion function.
func Convert_v1beta1_StorageClass_To_storage_StorageClass(in *v1beta1.StorageClass, out *storage.StorageClass, s conversion.Scope) error {
	return autoConvert_v1beta1_StorageClass_To_storage_StorageClass(in, out, s)
}

func autoConvert_storage_StorageClass_To_v1beta1_StorageClass(in *storage.StorageClass, out *v1beta1.StorageClass, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	out.Provisioner = in.Provisioner
	out.Parameters = *(*map[string]string)(unsafe.Pointer(&in.Parameters))
	out.ReclaimPolicy = (*v1.PersistentVolumeReclaimPolicy)(unsafe.Pointer(in.ReclaimPolicy))
	out.MountOptions = *(*[]string)(unsafe.Pointer(&in.MountOptions))
	out.AllowVolumeExpansion = (*bool)(unsafe.Pointer(in.AllowVolumeExpansion))
	return nil
}

// Convert_storage_StorageClass_To_v1beta1_StorageClass is an autogenerated conversion function.
func Convert_storage_StorageClass_To_v1beta1_StorageClass(in *storage.StorageClass, out *v1beta1.StorageClass, s conversion.Scope) error {
	return autoConvert_storage_StorageClass_To_v1beta1_StorageClass(in, out, s)
}

func autoConvert_v1beta1_StorageClassList_To_storage_StorageClassList(in *v1beta1.StorageClassList, out *storage.StorageClassList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]storage.StorageClass)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1beta1_StorageClassList_To_storage_StorageClassList is an autogenerated conversion function.
func Convert_v1beta1_StorageClassList_To_storage_StorageClassList(in *v1beta1.StorageClassList, out *storage.StorageClassList, s conversion.Scope) error {
	return autoConvert_v1beta1_StorageClassList_To_storage_StorageClassList(in, out, s)
}

func autoConvert_storage_StorageClassList_To_v1beta1_StorageClassList(in *storage.StorageClassList, out *v1beta1.StorageClassList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]v1beta1.StorageClass)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_storage_StorageClassList_To_v1beta1_StorageClassList is an autogenerated conversion function.
func Convert_storage_StorageClassList_To_v1beta1_StorageClassList(in *storage.StorageClassList, out *v1beta1.StorageClassList, s conversion.Scope) error {
	return autoConvert_storage_StorageClassList_To_v1beta1_StorageClassList(in, out, s)
}
