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

// GroupSubjectApplyConfiguration represents an declarative configuration of the GroupSubject type for use
// with apply.
type GroupSubjectApplyConfiguration struct {
	Name *string `json:"name,omitempty"`
}

// GroupSubjectApplyConfiguration constructs an declarative configuration of the GroupSubject type for use with
// apply.
func GroupSubject() *GroupSubjectApplyConfiguration {
	return &GroupSubjectApplyConfiguration{}
}

// SetName sets the Name field in the declarative configuration to the given value.
func (b *GroupSubjectApplyConfiguration) SetName(value string) *GroupSubjectApplyConfiguration {
	b.Name = &value
	return b
}

// RemoveName removes the Name field from the declarative configuration.
func (b *GroupSubjectApplyConfiguration) RemoveName() *GroupSubjectApplyConfiguration {
	b.Name = nil
	return b
}

// GetName gets the Name field from the declarative configuration.
func (b *GroupSubjectApplyConfiguration) GetName() (value string, ok bool) {
	if v := b.Name; v != nil {
		return *v, true
	}
	return value, false
}

// GroupSubjectList represents a listAlias of GroupSubjectApplyConfiguration.
type GroupSubjectList []*GroupSubjectApplyConfiguration

// GroupSubjectList represents a map of GroupSubjectApplyConfiguration.
type GroupSubjectMap map[string]GroupSubjectApplyConfiguration
