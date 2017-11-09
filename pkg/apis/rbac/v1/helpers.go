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

package v1

import (
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PolicyRuleBuilder let's us attach methods.  A no-no for API types.
// We use it to construct rules in code.  It's more compact than trying to write them
// out in a literal and allows us to perform some basic checking during construction
type PolicyRuleBuilder struct {
	PolicyRule rbacv1.PolicyRule `protobuf:"bytes,1,opt,name=policyRule"`
}

func NewRule(verbs ...string) *PolicyRuleBuilder {
	return &PolicyRuleBuilder{
		PolicyRule: rbacv1.PolicyRule{Verbs: verbs},
	}
}

func (r *PolicyRuleBuilder) Groups(groups ...string) *PolicyRuleBuilder {
	r.PolicyRule.APIGroups = append(r.PolicyRule.APIGroups, groups...)
	return r
}

func (r *PolicyRuleBuilder) Resources(resources ...string) *PolicyRuleBuilder {
	r.PolicyRule.Resources = append(r.PolicyRule.Resources, resources...)
	return r
}

func (r *PolicyRuleBuilder) Names(names ...string) *PolicyRuleBuilder {
	r.PolicyRule.ResourceNames = append(r.PolicyRule.ResourceNames, names...)
	return r
}

func (r *PolicyRuleBuilder) URLs(urls ...string) *PolicyRuleBuilder {
	r.PolicyRule.NonResourceURLs = append(r.PolicyRule.NonResourceURLs, urls...)
	return r
}

func (r *PolicyRuleBuilder) RuleOrDie() rbacv1.PolicyRule {
	ret, err := r.Rule()
	if err != nil {
		panic(err)
	}
	return ret
}

func (r *PolicyRuleBuilder) Rule() (rbacv1.PolicyRule, error) {
	if len(r.PolicyRule.Verbs) == 0 {
		return rbacv1.PolicyRule{}, fmt.Errorf("verbs are required: %#v", r.PolicyRule)
	}

	switch {
	case len(r.PolicyRule.NonResourceURLs) > 0:
		if len(r.PolicyRule.APIGroups) != 0 || len(r.PolicyRule.Resources) != 0 || len(r.PolicyRule.ResourceNames) != 0 {
			return rbacv1.PolicyRule{}, fmt.Errorf("non-resource rule may not have apiGroups, resources, or resourceNames: %#v", r.PolicyRule)
		}
	case len(r.PolicyRule.Resources) > 0:
		if len(r.PolicyRule.NonResourceURLs) != 0 {
			return rbacv1.PolicyRule{}, fmt.Errorf("resource rule may not have nonResourceURLs: %#v", r.PolicyRule)
		}
		if len(r.PolicyRule.APIGroups) == 0 {
			// this a common bug
			return rbacv1.PolicyRule{}, fmt.Errorf("resource rule must have apiGroups: %#v", r.PolicyRule)
		}
	default:
		return rbacv1.PolicyRule{}, fmt.Errorf("a rule must have either nonResourceURLs or resources: %#v", r.PolicyRule)
	}

	return r.PolicyRule, nil
}

// ClusterRoleBindingBuilder let's us attach methods.  A no-no for API types.
// We use it to construct bindings in code.  It's more compact than trying to write them
// out in a literal.
type ClusterRoleBindingBuilder struct {
	ClusterRoleBinding rbacv1.ClusterRoleBinding `protobuf:"bytes,1,opt,name=clusterRoleBinding"`
}

func NewClusterBinding(clusterRoleName string) *ClusterRoleBindingBuilder {
	return &ClusterRoleBindingBuilder{
		ClusterRoleBinding: rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{Name: clusterRoleName},
			RoleRef: rbacv1.RoleRef{
				APIGroup: rbacv1.GroupName,
				Kind:     "ClusterRole",
				Name:     clusterRoleName,
			},
		},
	}
}

func (r *ClusterRoleBindingBuilder) Groups(groups ...string) *ClusterRoleBindingBuilder {
	for _, group := range groups {
		r.ClusterRoleBinding.Subjects = append(r.ClusterRoleBinding.Subjects, rbacv1.Subject{Kind: rbacv1.GroupKind, Name: group})
	}
	return r
}

func (r *ClusterRoleBindingBuilder) Users(users ...string) *ClusterRoleBindingBuilder {
	for _, user := range users {
		r.ClusterRoleBinding.Subjects = append(r.ClusterRoleBinding.Subjects, rbacv1.Subject{Kind: rbacv1.UserKind, Name: user})
	}
	return r
}

func (r *ClusterRoleBindingBuilder) SAs(namespace string, serviceAccountNames ...string) *ClusterRoleBindingBuilder {
	for _, saName := range serviceAccountNames {
		r.ClusterRoleBinding.Subjects = append(r.ClusterRoleBinding.Subjects, rbacv1.Subject{Kind: rbacv1.ServiceAccountKind, Namespace: namespace, Name: saName})
	}
	return r
}

func (r *ClusterRoleBindingBuilder) BindingOrDie() rbacv1.ClusterRoleBinding {
	ret, err := r.Binding()
	if err != nil {
		panic(err)
	}
	return ret
}

func (r *ClusterRoleBindingBuilder) Binding() (rbacv1.ClusterRoleBinding, error) {
	if len(r.ClusterRoleBinding.Subjects) == 0 {
		return rbacv1.ClusterRoleBinding{}, fmt.Errorf("subjects are required: %#v", r.ClusterRoleBinding)
	}

	return r.ClusterRoleBinding, nil
}
