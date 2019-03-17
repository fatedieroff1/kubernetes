/*
Copyright 2019 The Kubernetes Authors.

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

// This file is generated. DO NOT EDIT.

package kflag

import(
	"github.com/spf13/pflag"
)

// Float64Value contains the scratch space for a registered float64 flag.
// Values can be applied from this scratch space to a target using the Set or Apply methods.
type Float64Value struct {
	name string
	value float64
	fs *pflag.FlagSet
}

// Float64Var registers a flag for type float64 against the FlagSet, and returns a struct
// of type Float64Value that contains the scratch space the flag will be parsed into.
func (fs *FlagSet) Float64Var(name string, def float64, usage string) *Float64Value {
	v := &Float64Value{
		name: name,
		fs: fs.fs,
	}
	fs.fs.Float64Var(&v.value, name, def, usage)
	return v
}

// Set copies the float64 value to the target if the flag was detected.
func (v *Float64Value) Set(target *float64) {
	if v.fs.Changed(v.name) {
		*target = v.value
	}
}

// Apply calls the user-provided apply function with the float64 value if the flag was detected.
func (v *Float64Value) Apply(apply func(value float64)) {
	if v.fs.Changed(v.name) {
		apply(v.value)
	}
}
