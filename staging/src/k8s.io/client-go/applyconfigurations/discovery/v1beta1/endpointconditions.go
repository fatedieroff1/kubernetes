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

// EndpointConditionsApplyConfiguration represents an declarative configuration of the EndpointConditions type for use
// with apply.
type EndpointConditionsApplyConfiguration struct {
	Ready       *bool `json:"ready,omitempty"`
	Serving     *bool `json:"serving,omitempty"`
	Terminating *bool `json:"terminating,omitempty"`
}

// EndpointConditionsApplyConfiguration constructs an declarative configuration of the EndpointConditions type for use with
// apply.
func EndpointConditions() *EndpointConditionsApplyConfiguration {
	return &EndpointConditionsApplyConfiguration{}
}

// SetReady sets the Ready field in the declarative configuration to the given value.
func (b *EndpointConditionsApplyConfiguration) SetReady(value bool) *EndpointConditionsApplyConfiguration {
	b.Ready = &value
	return b
}

// RemoveReady removes the Ready field from the declarative configuration.
func (b *EndpointConditionsApplyConfiguration) RemoveReady() *EndpointConditionsApplyConfiguration {
	b.Ready = nil
	return b
}

// GetReady gets the Ready field from the declarative configuration.
func (b *EndpointConditionsApplyConfiguration) GetReady() (value bool, ok bool) {
	if v := b.Ready; v != nil {
		return *v, true
	}
	return value, false
}

// SetServing sets the Serving field in the declarative configuration to the given value.
func (b *EndpointConditionsApplyConfiguration) SetServing(value bool) *EndpointConditionsApplyConfiguration {
	b.Serving = &value
	return b
}

// RemoveServing removes the Serving field from the declarative configuration.
func (b *EndpointConditionsApplyConfiguration) RemoveServing() *EndpointConditionsApplyConfiguration {
	b.Serving = nil
	return b
}

// GetServing gets the Serving field from the declarative configuration.
func (b *EndpointConditionsApplyConfiguration) GetServing() (value bool, ok bool) {
	if v := b.Serving; v != nil {
		return *v, true
	}
	return value, false
}

// SetTerminating sets the Terminating field in the declarative configuration to the given value.
func (b *EndpointConditionsApplyConfiguration) SetTerminating(value bool) *EndpointConditionsApplyConfiguration {
	b.Terminating = &value
	return b
}

// RemoveTerminating removes the Terminating field from the declarative configuration.
func (b *EndpointConditionsApplyConfiguration) RemoveTerminating() *EndpointConditionsApplyConfiguration {
	b.Terminating = nil
	return b
}

// GetTerminating gets the Terminating field from the declarative configuration.
func (b *EndpointConditionsApplyConfiguration) GetTerminating() (value bool, ok bool) {
	if v := b.Terminating; v != nil {
		return *v, true
	}
	return value, false
}

// EndpointConditionsList represents a listAlias of EndpointConditionsApplyConfiguration.
type EndpointConditionsList []*EndpointConditionsApplyConfiguration

// EndpointConditionsList represents a map of EndpointConditionsApplyConfiguration.
type EndpointConditionsMap map[string]EndpointConditionsApplyConfiguration
