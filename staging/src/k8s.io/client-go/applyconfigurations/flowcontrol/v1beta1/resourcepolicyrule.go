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

// ResourcePolicyRuleApplyConfiguration represents an declarative configuration of the ResourcePolicyRule type for use
// with apply.
type ResourcePolicyRuleApplyConfiguration struct {
	Verbs        *[]string `json:"verbs,omitempty"`
	APIGroups    *[]string `json:"apiGroups,omitempty"`
	Resources    *[]string `json:"resources,omitempty"`
	ClusterScope *bool     `json:"clusterScope,omitempty"`
	Namespaces   *[]string `json:"namespaces,omitempty"`
}

// ResourcePolicyRuleApplyConfiguration constructs an declarative configuration of the ResourcePolicyRule type for use with
// apply.
func ResourcePolicyRule() *ResourcePolicyRuleApplyConfiguration {
	return &ResourcePolicyRuleApplyConfiguration{}
}

// SetVerbs sets the Verbs field in the declarative configuration to the given value.
func (b *ResourcePolicyRuleApplyConfiguration) SetVerbs(value []string) *ResourcePolicyRuleApplyConfiguration {
	b.Verbs = &value
	return b
}

// RemoveVerbs removes the Verbs field from the declarative configuration.
func (b *ResourcePolicyRuleApplyConfiguration) RemoveVerbs() *ResourcePolicyRuleApplyConfiguration {
	b.Verbs = nil
	return b
}

// GetVerbs gets the Verbs field from the declarative configuration.
func (b *ResourcePolicyRuleApplyConfiguration) GetVerbs() (value []string, ok bool) {
	if v := b.Verbs; v != nil {
		return *v, true
	}
	return value, false
}

// SetAPIGroups sets the APIGroups field in the declarative configuration to the given value.
func (b *ResourcePolicyRuleApplyConfiguration) SetAPIGroups(value []string) *ResourcePolicyRuleApplyConfiguration {
	b.APIGroups = &value
	return b
}

// RemoveAPIGroups removes the APIGroups field from the declarative configuration.
func (b *ResourcePolicyRuleApplyConfiguration) RemoveAPIGroups() *ResourcePolicyRuleApplyConfiguration {
	b.APIGroups = nil
	return b
}

// GetAPIGroups gets the APIGroups field from the declarative configuration.
func (b *ResourcePolicyRuleApplyConfiguration) GetAPIGroups() (value []string, ok bool) {
	if v := b.APIGroups; v != nil {
		return *v, true
	}
	return value, false
}

// SetResources sets the Resources field in the declarative configuration to the given value.
func (b *ResourcePolicyRuleApplyConfiguration) SetResources(value []string) *ResourcePolicyRuleApplyConfiguration {
	b.Resources = &value
	return b
}

// RemoveResources removes the Resources field from the declarative configuration.
func (b *ResourcePolicyRuleApplyConfiguration) RemoveResources() *ResourcePolicyRuleApplyConfiguration {
	b.Resources = nil
	return b
}

// GetResources gets the Resources field from the declarative configuration.
func (b *ResourcePolicyRuleApplyConfiguration) GetResources() (value []string, ok bool) {
	if v := b.Resources; v != nil {
		return *v, true
	}
	return value, false
}

// SetClusterScope sets the ClusterScope field in the declarative configuration to the given value.
func (b *ResourcePolicyRuleApplyConfiguration) SetClusterScope(value bool) *ResourcePolicyRuleApplyConfiguration {
	b.ClusterScope = &value
	return b
}

// RemoveClusterScope removes the ClusterScope field from the declarative configuration.
func (b *ResourcePolicyRuleApplyConfiguration) RemoveClusterScope() *ResourcePolicyRuleApplyConfiguration {
	b.ClusterScope = nil
	return b
}

// GetClusterScope gets the ClusterScope field from the declarative configuration.
func (b *ResourcePolicyRuleApplyConfiguration) GetClusterScope() (value bool, ok bool) {
	if v := b.ClusterScope; v != nil {
		return *v, true
	}
	return value, false
}

// SetNamespaces sets the Namespaces field in the declarative configuration to the given value.
func (b *ResourcePolicyRuleApplyConfiguration) SetNamespaces(value []string) *ResourcePolicyRuleApplyConfiguration {
	b.Namespaces = &value
	return b
}

// RemoveNamespaces removes the Namespaces field from the declarative configuration.
func (b *ResourcePolicyRuleApplyConfiguration) RemoveNamespaces() *ResourcePolicyRuleApplyConfiguration {
	b.Namespaces = nil
	return b
}

// GetNamespaces gets the Namespaces field from the declarative configuration.
func (b *ResourcePolicyRuleApplyConfiguration) GetNamespaces() (value []string, ok bool) {
	if v := b.Namespaces; v != nil {
		return *v, true
	}
	return value, false
}

// ResourcePolicyRuleList represents a listAlias of ResourcePolicyRuleApplyConfiguration.
type ResourcePolicyRuleList []*ResourcePolicyRuleApplyConfiguration

// ResourcePolicyRuleList represents a map of ResourcePolicyRuleApplyConfiguration.
type ResourcePolicyRuleMap map[string]ResourcePolicyRuleApplyConfiguration
