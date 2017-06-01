/*
Copyright 2016 The Kubernetes Authors.

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

package v1alpha1

// This file contains a collection of methods that can be used from go-restful to
// generate Swagger API documentation for its models. Please read this PR for more
// information on the implementation: https://github.com/emicklei/go-restful/pull/215
//
// TODOs are ignored from the parser (e.g. TODO(andronat):... || TODO:...) if and only if
// they are on one line! For multiple line or blocks that you want to ignore use ---.
// Any context after a --- is ignored.
//
// Those methods can be generated by using hack/update-generated-swagger-docs.sh

// AUTO-GENERATED FUNCTIONS START HERE
var map_PodPreset = map[string]string{
	"": "PodPreset is a policy resource that defines additional runtime requirements for a Pod.",
}

func (PodPreset) SwaggerDoc() map[string]string {
	return map_PodPreset
}

var map_PodPresetList = map[string]string{
	"":         "PodPresetList is a list of PodPreset objects.",
	"metadata": "Standard list metadata. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata",
	"items":    "Items is a list of schema objects.",
}

func (PodPresetList) SwaggerDoc() map[string]string {
	return map_PodPresetList
}

var map_PodPresetSpec = map[string]string{
	"":             "PodPresetSpec is a description of a pod injection policy.",
	"selector":     "Selector is a label query over a set of resources, in this case pods. Required.",
	"env":          "Env defines the collection of EnvVar to inject into containers.",
	"envFrom":      "EnvFrom defines the collection of EnvFromSource to inject into containers.",
	"volumes":      "Volumes defines the collection of Volume to inject into the pod.",
	"volumeMounts": "VolumeMounts defines the collection of VolumeMount to inject into containers.",
}

func (PodPresetSpec) SwaggerDoc() map[string]string {
	return map_PodPresetSpec
}

// AUTO-GENERATED FUNCTIONS END HERE
