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

// DisruptionPolicyApplyConfiguration represents an declarative configuration of the DisruptionPolicy type for use
// with apply.
type DisruptionPolicyApplyConfiguration struct {
	Policy                     *v1.PDBDisruptionPolicy `json:"policy,omitempty"`
	PriorityGreaterThanOrEqual *int32                  `json:"priorityGreaterThanOrEqual,omitempty"`
}

// DisruptionPolicyApplyConfiguration constructs an declarative configuration of the DisruptionPolicy type for use with
// apply.
func DisruptionPolicy() *DisruptionPolicyApplyConfiguration {
	return &DisruptionPolicyApplyConfiguration{}
}

// WithPolicy sets the Policy field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Policy field is set to the value of the last call.
func (b *DisruptionPolicyApplyConfiguration) WithPolicy(value v1.PDBDisruptionPolicy) *DisruptionPolicyApplyConfiguration {
	b.Policy = &value
	return b
}

// WithPriorityGreaterThanOrEqual sets the PriorityGreaterThanOrEqual field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PriorityGreaterThanOrEqual field is set to the value of the last call.
func (b *DisruptionPolicyApplyConfiguration) WithPriorityGreaterThanOrEqual(value int32) *DisruptionPolicyApplyConfiguration {
	b.PriorityGreaterThanOrEqual = &value
	return b
}
