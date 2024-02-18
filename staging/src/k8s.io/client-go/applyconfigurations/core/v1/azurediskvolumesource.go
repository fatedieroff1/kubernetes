/*
Copyright The Kubernetes Authors.

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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	v1 "k8s.io/api/core/v1"
)

// AzureDiskVolumeSourceApplyConfiguration represents an declarative configuration of the AzureDiskVolumeSource type for use
// with apply.
type AzureDiskVolumeSourceApplyConfiguration struct {
	DiskName    *string                      `json:"diskName,omitempty"`
	DiskURI     *string                      `json:"diskURI,omitempty"`
	CachingMode *v1.AzureDataDiskCachingMode `json:"cachingMode,omitempty"`
	FSType      *string                      `json:"fsType,omitempty"`
	ReadOnly    *bool                        `json:"readOnly,omitempty"`
	Kind        *v1.AzureDataDiskKind        `json:"kind,omitempty"`
}

// AzureDiskVolumeSourceApplyConfiguration constructs an declarative configuration of the AzureDiskVolumeSource type for use with
// apply.
func AzureDiskVolumeSource() *AzureDiskVolumeSourceApplyConfiguration {
	return &AzureDiskVolumeSourceApplyConfiguration{}
}

// WithDiskName sets the DiskName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DiskName field is set to the value of the last call.
func (b *AzureDiskVolumeSourceApplyConfiguration) WithDiskName(value string) *AzureDiskVolumeSourceApplyConfiguration {
	b.DiskName = &value
	return b
}

// WithDiskURI sets the DiskURI field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DiskURI field is set to the value of the last call.
func (b *AzureDiskVolumeSourceApplyConfiguration) WithDiskURI(value string) *AzureDiskVolumeSourceApplyConfiguration {
	b.DiskURI = &value
	return b
}

// WithCachingMode sets the CachingMode field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CachingMode field is set to the value of the last call.
func (b *AzureDiskVolumeSourceApplyConfiguration) WithCachingMode(value v1.AzureDataDiskCachingMode) *AzureDiskVolumeSourceApplyConfiguration {
	b.CachingMode = &value
	return b
}

// WithFSType sets the FSType field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the FSType field is set to the value of the last call.
func (b *AzureDiskVolumeSourceApplyConfiguration) WithFSType(value string) *AzureDiskVolumeSourceApplyConfiguration {
	b.FSType = &value
	return b
}

// WithReadOnly sets the ReadOnly field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ReadOnly field is set to the value of the last call.
func (b *AzureDiskVolumeSourceApplyConfiguration) WithReadOnly(value bool) *AzureDiskVolumeSourceApplyConfiguration {
	b.ReadOnly = &value
	return b
}

// WithKind sets the Kind field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Kind field is set to the value of the last call.
func (b *AzureDiskVolumeSourceApplyConfiguration) WithKind(value v1.AzureDataDiskKind) *AzureDiskVolumeSourceApplyConfiguration {
	b.Kind = &value
	return b
}
