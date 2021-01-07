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

package v1beta1

import (
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

// RollingUpdateDaemonSetApplyConfiguration represents an declarative configuration of the RollingUpdateDaemonSet type for use
// with apply.
type RollingUpdateDaemonSetApplyConfiguration struct {
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`
	MaxSurge       *intstr.IntOrString `json:"maxSurge,omitempty"`
}

// RollingUpdateDaemonSetApplyConfiguration constructs an declarative configuration of the RollingUpdateDaemonSet type for use with
// apply.
func RollingUpdateDaemonSet() *RollingUpdateDaemonSetApplyConfiguration {
	return &RollingUpdateDaemonSetApplyConfiguration{}
}

// SetMaxUnavailable sets the MaxUnavailable field in the declarative configuration to the given value.
func (b *RollingUpdateDaemonSetApplyConfiguration) SetMaxUnavailable(value intstr.IntOrString) *RollingUpdateDaemonSetApplyConfiguration {
	b.MaxUnavailable = &value
	return b
}

// RemoveMaxUnavailable removes the MaxUnavailable field from the declarative configuration.
func (b *RollingUpdateDaemonSetApplyConfiguration) RemoveMaxUnavailable() *RollingUpdateDaemonSetApplyConfiguration {
	b.MaxUnavailable = nil
	return b
}

// GetMaxUnavailable gets the MaxUnavailable field from the declarative configuration.
func (b *RollingUpdateDaemonSetApplyConfiguration) GetMaxUnavailable() (value intstr.IntOrString, ok bool) {
	if v := b.MaxUnavailable; v != nil {
		return *v, true
	}
	return value, false
}

// SetMaxSurge sets the MaxSurge field in the declarative configuration to the given value.
func (b *RollingUpdateDaemonSetApplyConfiguration) SetMaxSurge(value intstr.IntOrString) *RollingUpdateDaemonSetApplyConfiguration {
	b.MaxSurge = &value
	return b
}

// RemoveMaxSurge removes the MaxSurge field from the declarative configuration.
func (b *RollingUpdateDaemonSetApplyConfiguration) RemoveMaxSurge() *RollingUpdateDaemonSetApplyConfiguration {
	b.MaxSurge = nil
	return b
}

// GetMaxSurge gets the MaxSurge field from the declarative configuration.
func (b *RollingUpdateDaemonSetApplyConfiguration) GetMaxSurge() (value intstr.IntOrString, ok bool) {
	if v := b.MaxSurge; v != nil {
		return *v, true
	}
	return value, false
}

// RollingUpdateDaemonSetList represents a listAlias of RollingUpdateDaemonSetApplyConfiguration.
type RollingUpdateDaemonSetList []*RollingUpdateDaemonSetApplyConfiguration

// RollingUpdateDaemonSetList represents a map of RollingUpdateDaemonSetApplyConfiguration.
type RollingUpdateDaemonSetMap map[string]RollingUpdateDaemonSetApplyConfiguration
