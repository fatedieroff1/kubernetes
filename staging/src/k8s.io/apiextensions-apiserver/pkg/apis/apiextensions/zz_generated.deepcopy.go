//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package apiextensions

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomResourceColumnDefinition) DeepCopyInto(out *CustomResourceColumnDefinition) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomResourceColumnDefinition.
func (in *CustomResourceColumnDefinition) DeepCopy() *CustomResourceColumnDefinition {
	if in == nil {
		return nil
	}
	out := new(CustomResourceColumnDefinition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomResourceConversion) DeepCopyInto(out *CustomResourceConversion) {
	*out = *in
	if in.WebhookClientConfig != nil {
		in, out := &in.WebhookClientConfig, &out.WebhookClientConfig
		*out = new(WebhookClientConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.ConversionReviewVersions != nil {
		in, out := &in.ConversionReviewVersions, &out.ConversionReviewVersions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomResourceConversion.
func (in *CustomResourceConversion) DeepCopy() *CustomResourceConversion {
	if in == nil {
		return nil
	}
	out := new(CustomResourceConversion)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomResourceDefinition) DeepCopyInto(out *CustomResourceDefinition) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomResourceDefinition.
func (in *CustomResourceDefinition) DeepCopy() *CustomResourceDefinition {
	if in == nil {
		return nil
	}
	out := new(CustomResourceDefinition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CustomResourceDefinition) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomResourceDefinitionCondition) DeepCopyInto(out *CustomResourceDefinitionCondition) {
	*out = *in
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomResourceDefinitionCondition.
func (in *CustomResourceDefinitionCondition) DeepCopy() *CustomResourceDefinitionCondition {
	if in == nil {
		return nil
	}
	out := new(CustomResourceDefinitionCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomResourceDefinitionFeatureGate) DeepCopyInto(out *CustomResourceDefinitionFeatureGate) {
	*out = *in
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		*out = new(bool)
		**out = **in
	}
	if in.Default != nil {
		in, out := &in.Default, &out.Default
		*out = new(bool)
		**out = **in
	}
	if in.FieldPaths != nil {
		in, out := &in.FieldPaths, &out.FieldPaths
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomResourceDefinitionFeatureGate.
func (in *CustomResourceDefinitionFeatureGate) DeepCopy() *CustomResourceDefinitionFeatureGate {
	if in == nil {
		return nil
	}
	out := new(CustomResourceDefinitionFeatureGate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomResourceDefinitionFeatureGates) DeepCopyInto(out *CustomResourceDefinitionFeatureGates) {
	*out = *in
	if in.FeatureGates != nil {
		in, out := &in.FeatureGates, &out.FeatureGates
		*out = make([]CustomResourceDefinitionFeatureGate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Component != nil {
		in, out := &in.Component, &out.Component
		*out = new(string)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomResourceDefinitionFeatureGates.
func (in *CustomResourceDefinitionFeatureGates) DeepCopy() *CustomResourceDefinitionFeatureGates {
	if in == nil {
		return nil
	}
	out := new(CustomResourceDefinitionFeatureGates)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomResourceDefinitionList) DeepCopyInto(out *CustomResourceDefinitionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CustomResourceDefinition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomResourceDefinitionList.
func (in *CustomResourceDefinitionList) DeepCopy() *CustomResourceDefinitionList {
	if in == nil {
		return nil
	}
	out := new(CustomResourceDefinitionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CustomResourceDefinitionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomResourceDefinitionNames) DeepCopyInto(out *CustomResourceDefinitionNames) {
	*out = *in
	if in.ShortNames != nil {
		in, out := &in.ShortNames, &out.ShortNames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Categories != nil {
		in, out := &in.Categories, &out.Categories
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomResourceDefinitionNames.
func (in *CustomResourceDefinitionNames) DeepCopy() *CustomResourceDefinitionNames {
	if in == nil {
		return nil
	}
	out := new(CustomResourceDefinitionNames)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomResourceDefinitionSpec) DeepCopyInto(out *CustomResourceDefinitionSpec) {
	*out = *in
	in.Names.DeepCopyInto(&out.Names)
	if in.Validation != nil {
		in, out := &in.Validation, &out.Validation
		*out = new(CustomResourceValidation)
		(*in).DeepCopyInto(*out)
	}
	if in.Subresources != nil {
		in, out := &in.Subresources, &out.Subresources
		*out = new(CustomResourceSubresources)
		(*in).DeepCopyInto(*out)
	}
	if in.Versions != nil {
		in, out := &in.Versions, &out.Versions
		*out = make([]CustomResourceDefinitionVersion, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.AdditionalPrinterColumns != nil {
		in, out := &in.AdditionalPrinterColumns, &out.AdditionalPrinterColumns
		*out = make([]CustomResourceColumnDefinition, len(*in))
		copy(*out, *in)
	}
	if in.Conversion != nil {
		in, out := &in.Conversion, &out.Conversion
		*out = new(CustomResourceConversion)
		(*in).DeepCopyInto(*out)
	}
	if in.PreserveUnknownFields != nil {
		in, out := &in.PreserveUnknownFields, &out.PreserveUnknownFields
		*out = new(bool)
		**out = **in
	}
	if in.CustomFeatureGates != nil {
		in, out := &in.CustomFeatureGates, &out.CustomFeatureGates
		*out = new(CustomResourceDefinitionFeatureGates)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomResourceDefinitionSpec.
func (in *CustomResourceDefinitionSpec) DeepCopy() *CustomResourceDefinitionSpec {
	if in == nil {
		return nil
	}
	out := new(CustomResourceDefinitionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomResourceDefinitionStatus) DeepCopyInto(out *CustomResourceDefinitionStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]CustomResourceDefinitionCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.AcceptedNames.DeepCopyInto(&out.AcceptedNames)
	if in.StoredVersions != nil {
		in, out := &in.StoredVersions, &out.StoredVersions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomResourceDefinitionStatus.
func (in *CustomResourceDefinitionStatus) DeepCopy() *CustomResourceDefinitionStatus {
	if in == nil {
		return nil
	}
	out := new(CustomResourceDefinitionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomResourceDefinitionVersion) DeepCopyInto(out *CustomResourceDefinitionVersion) {
	*out = *in
	if in.DeprecationWarning != nil {
		in, out := &in.DeprecationWarning, &out.DeprecationWarning
		*out = new(string)
		**out = **in
	}
	if in.Schema != nil {
		in, out := &in.Schema, &out.Schema
		*out = new(CustomResourceValidation)
		(*in).DeepCopyInto(*out)
	}
	if in.Subresources != nil {
		in, out := &in.Subresources, &out.Subresources
		*out = new(CustomResourceSubresources)
		(*in).DeepCopyInto(*out)
	}
	if in.AdditionalPrinterColumns != nil {
		in, out := &in.AdditionalPrinterColumns, &out.AdditionalPrinterColumns
		*out = make([]CustomResourceColumnDefinition, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomResourceDefinitionVersion.
func (in *CustomResourceDefinitionVersion) DeepCopy() *CustomResourceDefinitionVersion {
	if in == nil {
		return nil
	}
	out := new(CustomResourceDefinitionVersion)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomResourceSubresourceScale) DeepCopyInto(out *CustomResourceSubresourceScale) {
	*out = *in
	if in.LabelSelectorPath != nil {
		in, out := &in.LabelSelectorPath, &out.LabelSelectorPath
		*out = new(string)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomResourceSubresourceScale.
func (in *CustomResourceSubresourceScale) DeepCopy() *CustomResourceSubresourceScale {
	if in == nil {
		return nil
	}
	out := new(CustomResourceSubresourceScale)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomResourceSubresourceStatus) DeepCopyInto(out *CustomResourceSubresourceStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomResourceSubresourceStatus.
func (in *CustomResourceSubresourceStatus) DeepCopy() *CustomResourceSubresourceStatus {
	if in == nil {
		return nil
	}
	out := new(CustomResourceSubresourceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomResourceSubresources) DeepCopyInto(out *CustomResourceSubresources) {
	*out = *in
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(CustomResourceSubresourceStatus)
		**out = **in
	}
	if in.Scale != nil {
		in, out := &in.Scale, &out.Scale
		*out = new(CustomResourceSubresourceScale)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomResourceSubresources.
func (in *CustomResourceSubresources) DeepCopy() *CustomResourceSubresources {
	if in == nil {
		return nil
	}
	out := new(CustomResourceSubresources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomResourceValidation) DeepCopyInto(out *CustomResourceValidation) {
	*out = *in
	if in.OpenAPIV3Schema != nil {
		in, out := &in.OpenAPIV3Schema, &out.OpenAPIV3Schema
		*out = (*in).DeepCopy()
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomResourceValidation.
func (in *CustomResourceValidation) DeepCopy() *CustomResourceValidation {
	if in == nil {
		return nil
	}
	out := new(CustomResourceValidation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalDocumentation) DeepCopyInto(out *ExternalDocumentation) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalDocumentation.
func (in *ExternalDocumentation) DeepCopy() *ExternalDocumentation {
	if in == nil {
		return nil
	}
	out := new(ExternalDocumentation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in JSONSchemaDefinitions) DeepCopyInto(out *JSONSchemaDefinitions) {
	{
		in := &in
		*out = make(JSONSchemaDefinitions, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JSONSchemaDefinitions.
func (in JSONSchemaDefinitions) DeepCopy() JSONSchemaDefinitions {
	if in == nil {
		return nil
	}
	out := new(JSONSchemaDefinitions)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in JSONSchemaDependencies) DeepCopyInto(out *JSONSchemaDependencies) {
	{
		in := &in
		*out = make(JSONSchemaDependencies, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JSONSchemaDependencies.
func (in JSONSchemaDependencies) DeepCopy() JSONSchemaDependencies {
	if in == nil {
		return nil
	}
	out := new(JSONSchemaDependencies)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JSONSchemaProps) DeepCopyInto(out *JSONSchemaProps) {
	clone := in.DeepCopy()
	*out = *clone
	return
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JSONSchemaPropsOrArray) DeepCopyInto(out *JSONSchemaPropsOrArray) {
	*out = *in
	if in.Schema != nil {
		in, out := &in.Schema, &out.Schema
		*out = (*in).DeepCopy()
	}
	if in.JSONSchemas != nil {
		in, out := &in.JSONSchemas, &out.JSONSchemas
		*out = make([]JSONSchemaProps, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JSONSchemaPropsOrArray.
func (in *JSONSchemaPropsOrArray) DeepCopy() *JSONSchemaPropsOrArray {
	if in == nil {
		return nil
	}
	out := new(JSONSchemaPropsOrArray)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JSONSchemaPropsOrBool) DeepCopyInto(out *JSONSchemaPropsOrBool) {
	*out = *in
	if in.Schema != nil {
		in, out := &in.Schema, &out.Schema
		*out = (*in).DeepCopy()
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JSONSchemaPropsOrBool.
func (in *JSONSchemaPropsOrBool) DeepCopy() *JSONSchemaPropsOrBool {
	if in == nil {
		return nil
	}
	out := new(JSONSchemaPropsOrBool)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JSONSchemaPropsOrStringArray) DeepCopyInto(out *JSONSchemaPropsOrStringArray) {
	*out = *in
	if in.Schema != nil {
		in, out := &in.Schema, &out.Schema
		*out = (*in).DeepCopy()
	}
	if in.Property != nil {
		in, out := &in.Property, &out.Property
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JSONSchemaPropsOrStringArray.
func (in *JSONSchemaPropsOrStringArray) DeepCopy() *JSONSchemaPropsOrStringArray {
	if in == nil {
		return nil
	}
	out := new(JSONSchemaPropsOrStringArray)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceReference) DeepCopyInto(out *ServiceReference) {
	*out = *in
	if in.Path != nil {
		in, out := &in.Path, &out.Path
		*out = new(string)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceReference.
func (in *ServiceReference) DeepCopy() *ServiceReference {
	if in == nil {
		return nil
	}
	out := new(ServiceReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValidationRule) DeepCopyInto(out *ValidationRule) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValidationRule.
func (in *ValidationRule) DeepCopy() *ValidationRule {
	if in == nil {
		return nil
	}
	out := new(ValidationRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in ValidationRules) DeepCopyInto(out *ValidationRules) {
	{
		in := &in
		*out = make(ValidationRules, len(*in))
		copy(*out, *in)
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValidationRules.
func (in ValidationRules) DeepCopy() ValidationRules {
	if in == nil {
		return nil
	}
	out := new(ValidationRules)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebhookClientConfig) DeepCopyInto(out *WebhookClientConfig) {
	*out = *in
	if in.URL != nil {
		in, out := &in.URL, &out.URL
		*out = new(string)
		**out = **in
	}
	if in.Service != nil {
		in, out := &in.Service, &out.Service
		*out = new(ServiceReference)
		(*in).DeepCopyInto(*out)
	}
	if in.CABundle != nil {
		in, out := &in.CABundle, &out.CABundle
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebhookClientConfig.
func (in *WebhookClientConfig) DeepCopy() *WebhookClientConfig {
	if in == nil {
		return nil
	}
	out := new(WebhookClientConfig)
	in.DeepCopyInto(out)
	return out
}
