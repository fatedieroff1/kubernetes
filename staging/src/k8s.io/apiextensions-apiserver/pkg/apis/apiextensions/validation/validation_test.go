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

package validation

import (
	"testing"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type validationMatch struct {
	path      *field.Path
	errorType field.ErrorType
}

func required(path ...string) validationMatch {
	return validationMatch{path: field.NewPath(path[0], path[1:]...), errorType: field.ErrorTypeRequired}
}
func invalid(path ...string) validationMatch {
	return validationMatch{path: field.NewPath(path[0], path[1:]...), errorType: field.ErrorTypeInvalid}
}
func unsupported(path ...string) validationMatch {
	return validationMatch{path: field.NewPath(path[0], path[1:]...), errorType: field.ErrorTypeNotSupported}
}
func immutable(path ...string) validationMatch {
	return validationMatch{path: field.NewPath(path[0], path[1:]...), errorType: field.ErrorTypeInvalid}
}
func forbidden(path ...string) validationMatch {
	return validationMatch{path: field.NewPath(path[0], path[1:]...), errorType: field.ErrorTypeForbidden}
}

func (v validationMatch) matches(err *field.Error) bool {
	return err.Type == v.errorType && err.Field == v.path.String()
}

func TestValidateCustomResourceDefinition(t *testing.T) {
	singleVersionList := []apiextensions.CustomResourceDefinitionVersion{
		{
			Name:    "version",
			Served:  true,
			Storage: true,
		},
	}
	tests := []struct {
		name     string
		resource *apiextensions.CustomResourceDefinition
		errors   []validationMatch
	}{
		{
			name: "no_storage_version",
			resource: &apiextensions.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{Name: "plural.group.com"},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group: "group.com",
					Scope: apiextensions.ResourceScope("Cluster"),
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "Plural",
						ListKind: "PluralList",
					},
					Versions: []apiextensions.CustomResourceDefinitionVersion{
						{
							Name:    "version",
							Served:  true,
							Storage: false,
						},
						{
							Name:    "version2",
							Served:  true,
							Storage: false,
						},
					},
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					StoredVersions: []string{"version"},
				},
			},
			errors: []validationMatch{
				invalid("spec", "versions"),
			},
		},
		{
			name: "multiple_storage_version",
			resource: &apiextensions.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{Name: "plural.group.com"},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group: "group.com",
					Scope: apiextensions.ResourceScope("Cluster"),
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "Plural",
						ListKind: "PluralList",
					},
					Versions: []apiextensions.CustomResourceDefinitionVersion{
						{
							Name:    "version",
							Served:  true,
							Storage: true,
						},
						{
							Name:    "version2",
							Served:  true,
							Storage: true,
						},
					},
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					StoredVersions: []string{"version"},
				},
			},
			errors: []validationMatch{
				invalid("spec", "versions"),
				invalid("status", "storedVersions"),
			},
		},
		{
			name: "missing_storage_version_in_stored_versions",
			resource: &apiextensions.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{Name: "plural.group.com"},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group: "group.com",
					Scope: apiextensions.ResourceScope("Cluster"),
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "Plural",
						ListKind: "PluralList",
					},
					Versions: []apiextensions.CustomResourceDefinitionVersion{
						{
							Name:    "version",
							Served:  true,
							Storage: false,
						},
						{
							Name:    "version2",
							Served:  true,
							Storage: true,
						},
					},
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					StoredVersions: []string{"version"},
				},
			},
			errors: []validationMatch{
				invalid("status", "storedVersions"),
			},
		},
		{
			name: "empty_stored_version",
			resource: &apiextensions.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{Name: "plural.group.com"},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group: "group.com",
					Scope: apiextensions.ResourceScope("Cluster"),
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "Plural",
						ListKind: "PluralList",
					},
					Versions: []apiextensions.CustomResourceDefinitionVersion{
						{
							Name:    "version",
							Served:  true,
							Storage: true,
						},
					},
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					StoredVersions: []string{},
				},
			},
			errors: []validationMatch{
				invalid("status", "storedVersions"),
			},
		},
		{
			name: "mismatched name",
			resource: &apiextensions.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{Name: "plural.not.group.com"},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group: "group.com",
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural: "plural",
					},
				},
			},
			errors: []validationMatch{
				invalid("status", "storedVersions"),
				invalid("metadata", "name"),
				invalid("spec", "versions"),
				required("spec", "scope"),
				required("spec", "names", "singular"),
				required("spec", "names", "kind"),
				required("spec", "names", "listKind"),
			},
		},
		{
			name: "missing values",
			resource: &apiextensions.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{Name: "plural.group.com"},
			},
			errors: []validationMatch{
				invalid("status", "storedVersions"),
				invalid("metadata", "name"),
				invalid("spec", "versions"),
				required("spec", "group"),
				required("spec", "scope"),
				required("spec", "names", "plural"),
				required("spec", "names", "singular"),
				required("spec", "names", "kind"),
				required("spec", "names", "listKind"),
			},
		},
		{
			name: "bad names 01",
			resource: &apiextensions.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{Name: "plural.group"},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group:   "group",
					Version: "ve()*rsion",
					Scope:   apiextensions.ResourceScope("foo"),
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:   "pl()*ural",
						Singular: "value()*a",
						Kind:     "value()*a",
						ListKind: "value()*a",
					},
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					AcceptedNames: apiextensions.CustomResourceDefinitionNames{
						Plural:   "pl()*ural",
						Singular: "value()*a",
						Kind:     "value()*a",
						ListKind: "value()*a",
					},
				},
			},
			errors: []validationMatch{
				invalid("status", "storedVersions"),
				invalid("metadata", "name"),
				invalid("spec", "group"),
				unsupported("spec", "scope"),
				invalid("spec", "names", "plural"),
				invalid("spec", "names", "singular"),
				invalid("spec", "names", "kind"),
				invalid("spec", "names", "listKind"), // invalid format
				invalid("spec", "names", "listKind"), // kind == listKind
				invalid("status", "acceptedNames", "plural"),
				invalid("status", "acceptedNames", "singular"),
				invalid("status", "acceptedNames", "kind"),
				invalid("status", "acceptedNames", "listKind"), // invalid format
				invalid("status", "acceptedNames", "listKind"), // kind == listKind
				invalid("spec", "versions"),
				invalid("spec", "version"),
			},
		},
		{
			name: "bad names 02",
			resource: &apiextensions.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{Name: "plural.group"},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group:    "group.c(*&om",
					Version:  "version",
					Versions: singleVersionList,
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "matching",
						ListKind: "matching",
					},
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					AcceptedNames: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "matching",
						ListKind: "matching",
					},
					StoredVersions: []string{"version"},
				},
			},
			errors: []validationMatch{
				invalid("metadata", "name"),
				invalid("spec", "group"),
				required("spec", "scope"),
				invalid("spec", "names", "listKind"),
				invalid("status", "acceptedNames", "listKind"),
			},
		},
		{
			name: "additionalProperties and properties forbidden",
			resource: &apiextensions.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{Name: "plural.group.com"},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group:    "group.com",
					Version:  "version",
					Versions: singleVersionList,
					Scope:    apiextensions.NamespaceScoped,
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "Plural",
						ListKind: "PluralList",
					},
					Validation: &apiextensions.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
							Properties: map[string]apiextensions.JSONSchemaProps{
								"foo": {},
							},
							AdditionalProperties: &apiextensions.JSONSchemaPropsOrBool{Allows: false},
						},
					},
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					StoredVersions: []string{"version"},
				},
			},
			errors: []validationMatch{
				forbidden("spec", "validation", "openAPIV3Schema", "additionalProperties"),
			},
		},
		{
			name: "additionalProperties without properties allowed (map[string]string)",
			resource: &apiextensions.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{Name: "plural.group.com"},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group:    "group.com",
					Version:  "version",
					Versions: singleVersionList,
					Scope:    apiextensions.NamespaceScoped,
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "Plural",
						ListKind: "PluralList",
					},
					Validation: &apiextensions.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
							AdditionalProperties: &apiextensions.JSONSchemaPropsOrBool{
								Allows: true,
								Schema: &apiextensions.JSONSchemaProps{
									Type: "string",
								},
							},
						},
					},
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					StoredVersions: []string{"version"},
				},
			},
			errors: []validationMatch{},
		},
	}

	for _, tc := range tests {
		errs := ValidateCustomResourceDefinition(tc.resource)
		seenErrs := make([]bool, len(errs))

		for _, expectedError := range tc.errors {
			found := false
			for i, err := range errs {
				if expectedError.matches(err) && !seenErrs[i] {
					found = true
					seenErrs[i] = true
					break
				}
			}

			if !found {
				t.Errorf("%s: expected %v at %v, got %v", tc.name, expectedError.errorType, expectedError.path.String(), errs)
			}
		}

		for i, seen := range seenErrs {
			if !seen {
				t.Errorf("%s: unexpected error: %v", tc.name, errs[i])
			}
		}
	}
}

func TestValidateCustomResourceDefinitionUpdate(t *testing.T) {
	tests := []struct {
		name     string
		old      *apiextensions.CustomResourceDefinition
		resource *apiextensions.CustomResourceDefinition
		errors   []validationMatch
	}{
		{
			name: "unchanged",
			old: &apiextensions.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "plural.group.com",
					ResourceVersion: "42",
				},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group:   "group.com",
					Version: "version",
					Versions: []apiextensions.CustomResourceDefinitionVersion{
						{
							Name:    "version",
							Served:  true,
							Storage: true,
						},
					},
					Scope: apiextensions.ResourceScope("Cluster"),
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "kind",
						ListKind: "listkind",
					},
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					AcceptedNames: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "kind",
						ListKind: "listkind",
					},
				},
			},
			resource: &apiextensions.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "plural.group.com",
					ResourceVersion: "42",
				},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group:   "group.com",
					Version: "version",
					Versions: []apiextensions.CustomResourceDefinitionVersion{
						{
							Name:    "version",
							Served:  true,
							Storage: true,
						},
					},
					Scope: apiextensions.ResourceScope("Cluster"),
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "kind",
						ListKind: "listkind",
					},
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					AcceptedNames: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "kind",
						ListKind: "listkind",
					},
					StoredVersions: []string{"version"},
				},
			},
			errors: []validationMatch{},
		},
		{
			name: "unchanged-established",
			old: &apiextensions.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "plural.group.com",
					ResourceVersion: "42",
				},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group:   "group.com",
					Version: "version",
					Versions: []apiextensions.CustomResourceDefinitionVersion{
						{
							Name:    "version",
							Served:  true,
							Storage: true,
						},
					},
					Scope: apiextensions.ResourceScope("Cluster"),
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "kind",
						ListKind: "listkind",
					},
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					AcceptedNames: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "kind",
						ListKind: "listkind",
					},
					Conditions: []apiextensions.CustomResourceDefinitionCondition{
						{Type: apiextensions.Established, Status: apiextensions.ConditionTrue},
					},
				},
			},
			resource: &apiextensions.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "plural.group.com",
					ResourceVersion: "42",
				},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group:   "group.com",
					Version: "version",
					Versions: []apiextensions.CustomResourceDefinitionVersion{
						{
							Name:    "version",
							Served:  true,
							Storage: true,
						},
					},
					Scope: apiextensions.ResourceScope("Cluster"),
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "kind",
						ListKind: "listkind",
					},
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					AcceptedNames: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "kind",
						ListKind: "listkind",
					},
					StoredVersions: []string{"version"},
				},
			},
			errors: []validationMatch{},
		},
		{
			name: "version-deleted",
			old: &apiextensions.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "plural.group.com",
					ResourceVersion: "42",
				},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group:   "group.com",
					Version: "version",
					Versions: []apiextensions.CustomResourceDefinitionVersion{
						{
							Name:    "version",
							Served:  true,
							Storage: true,
						},
						{
							Name:    "version2",
							Served:  true,
							Storage: false,
						},
					},
					Scope: apiextensions.ResourceScope("Cluster"),
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "kind",
						ListKind: "listkind",
					},
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					AcceptedNames: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "kind",
						ListKind: "listkind",
					},
					StoredVersions: []string{"version", "version2"},
					Conditions: []apiextensions.CustomResourceDefinitionCondition{
						{Type: apiextensions.Established, Status: apiextensions.ConditionTrue},
					},
				},
			},
			resource: &apiextensions.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "plural.group.com",
					ResourceVersion: "42",
				},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group:   "group.com",
					Version: "version",
					Versions: []apiextensions.CustomResourceDefinitionVersion{
						{
							Name:    "version",
							Served:  true,
							Storage: true,
						},
					},
					Scope: apiextensions.ResourceScope("Cluster"),
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "kind",
						ListKind: "listkind",
					},
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					AcceptedNames: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "kind",
						ListKind: "listkind",
					},
					StoredVersions: []string{"version", "version2"},
				},
			},
			errors: []validationMatch{
				invalid("status", "storedVersions[1]"),
			},
		},
		{
			name: "changes",
			old: &apiextensions.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "plural.group.com",
					ResourceVersion: "42",
				},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group:   "group.com",
					Version: "version",
					Versions: []apiextensions.CustomResourceDefinitionVersion{
						{
							Name:    "version",
							Served:  true,
							Storage: true,
						},
					},
					Scope: apiextensions.ResourceScope("Cluster"),
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "kind",
						ListKind: "listkind",
					},
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					AcceptedNames: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "kind",
						ListKind: "listkind",
					},
					Conditions: []apiextensions.CustomResourceDefinitionCondition{
						{Type: apiextensions.Established, Status: apiextensions.ConditionFalse},
					},
				},
			},
			resource: &apiextensions.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "plural.group.com",
					ResourceVersion: "42",
				},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group:   "abc.com",
					Version: "version2",
					Versions: []apiextensions.CustomResourceDefinitionVersion{
						{
							Name:    "version2",
							Served:  true,
							Storage: true,
						},
					},
					Scope: apiextensions.ResourceScope("Namespaced"),
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural2",
						Singular: "singular2",
						Kind:     "kind2",
						ListKind: "listkind2",
					},
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					AcceptedNames: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural2",
						Singular: "singular2",
						Kind:     "kind2",
						ListKind: "listkind2",
					},
					StoredVersions: []string{"version2"},
				},
			},
			errors: []validationMatch{
				immutable("spec", "group"),
				immutable("spec", "names", "plural"),
			},
		},
		{
			name: "changes-established",
			old: &apiextensions.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "plural.group.com",
					ResourceVersion: "42",
				},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group:   "group.com",
					Version: "version",
					Versions: []apiextensions.CustomResourceDefinitionVersion{
						{
							Name:    "version",
							Served:  true,
							Storage: true,
						},
					},
					Scope: apiextensions.ResourceScope("Cluster"),
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "kind",
						ListKind: "listkind",
					},
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					AcceptedNames: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural",
						Singular: "singular",
						Kind:     "kind",
						ListKind: "listkind",
					},
					Conditions: []apiextensions.CustomResourceDefinitionCondition{
						{Type: apiextensions.Established, Status: apiextensions.ConditionTrue},
					},
				},
			},
			resource: &apiextensions.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "plural.group.com",
					ResourceVersion: "42",
				},
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Group:   "abc.com",
					Version: "version2",
					Versions: []apiextensions.CustomResourceDefinitionVersion{
						{
							Name:    "version2",
							Served:  true,
							Storage: true,
						},
					},
					Scope: apiextensions.ResourceScope("Namespaced"),
					Names: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural2",
						Singular: "singular2",
						Kind:     "kind2",
						ListKind: "listkind2",
					},
				},
				Status: apiextensions.CustomResourceDefinitionStatus{
					AcceptedNames: apiextensions.CustomResourceDefinitionNames{
						Plural:   "plural2",
						Singular: "singular2",
						Kind:     "kind2",
						ListKind: "listkind2",
					},
					StoredVersions: []string{"version2"},
				},
			},
			errors: []validationMatch{
				immutable("spec", "group"),
				immutable("spec", "scope"),
				immutable("spec", "names", "kind"),
				immutable("spec", "names", "plural"),
			},
		},
	}

	for _, tc := range tests {
		errs := ValidateCustomResourceDefinitionUpdate(tc.resource, tc.old)
		seenErrs := make([]bool, len(errs))

		for _, expectedError := range tc.errors {
			found := false
			for i, err := range errs {
				if expectedError.matches(err) && !seenErrs[i] {
					found = true
					seenErrs[i] = true
					break
				}
			}

			if !found {
				t.Errorf("%s: expected %v at %v, got %v", tc.name, expectedError.errorType, expectedError.path.String(), errs)
			}
		}

		for i, seen := range seenErrs {
			if !seen {
				t.Errorf("%s: unexpected error: %v", tc.name, errs[i])
			}
		}
	}
}

func TestValidateCustomResourceDefinitionValidation(t *testing.T) {
	tests := []struct {
		name          string
		input         apiextensions.CustomResourceValidation
		statusEnabled bool
		wantError     bool
	}{
		{
			name:      "empty",
			input:     apiextensions.CustomResourceValidation{},
			wantError: false,
		},
		{
			name:          "empty with status",
			input:         apiextensions.CustomResourceValidation{},
			statusEnabled: true,
			wantError:     false,
		},
		{
			name: "root type without status",
			input: apiextensions.CustomResourceValidation{
				OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
					Type: "object",
				},
			},
			statusEnabled: false,
			wantError:     false,
		},
		{
			name: "root type with status",
			input: apiextensions.CustomResourceValidation{
				OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
					Type: "object",
				},
			},
			statusEnabled: true,
			wantError:     true,
		},
		{
			name: "properties, required and description with status",
			input: apiextensions.CustomResourceValidation{
				OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
					Properties: map[string]apiextensions.JSONSchemaProps{
						"spec":   {},
						"status": {},
					},
					Required:    []string{"spec", "status"},
					Description: "This is a description",
				},
			},
			statusEnabled: true,
			wantError:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateCustomResourceDefinitionValidation(&tt.input, tt.statusEnabled, field.NewPath("spec", "validation"))
			if !tt.wantError && len(got) > 0 {
				t.Errorf("Expected no error, but got: %v", got)
			} else if tt.wantError && len(got) == 0 {
				t.Error("Expected error, but got none")
			}
		})
	}
}

func TestValidateCustomResourceDefinitionValidationPattern(t *testing.T) {
	tests := []struct {
		name      string
		input     apiextensions.CustomResourceValidation
		wantError bool
	}{
		{
			name: "invalid regex",
			input: apiextensions.CustomResourceValidation{
				OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
					Pattern: `+`,
				},
			},
			wantError: true,
		},
		{
			name: "valid regex",
			input: apiextensions.CustomResourceValidation{
				OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
					Pattern: `^k[0-8]s$`,
				},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateCustomResourceDefinitionValidation(&tt.input, false, field.NewPath("spec", "validation"))
			if !tt.wantError && len(got) > 0 {
				t.Errorf("Expected no error, but got: %v", got)
			} else if tt.wantError && len(got) == 0 {
				t.Error("Expected error, but got none")
			}
		})
	}
}
