/*
Copyright 2018 The Kubernetes Authors.

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

package options

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"

	// libs that provide registration functions
	"k8s.io/component-base/logs"
	"k8s.io/component-base/version/verflag"
	"k8s.io/klog/v2"

	// ensure libs have a chance to globally register their flags
	_ "k8s.io/kubernetes/pkg/credentialprovider/azure"
	_ "k8s.io/kubernetes/pkg/credentialprovider/gcp"
)

// AddGlobalFlags explicitly registers flags that libraries (glog, verflag, etc.) register
// against the global flagsets from "flag" and "github.com/spf13/pflag".
// We do this in order to prevent unwanted flags from leaking into the Kubelet's flagset.
func AddGlobalFlags(fs *pflag.FlagSet) {
	addKlogFlags(fs)
	addCadvisorFlags(fs)
	addCredentialProviderFlags(fs)
	verflag.AddFlags(fs)
	logs.AddFlags(fs)
}

// normalize replaces underscores with hyphens
// we should always use hyphens instead of underscores when registering kubelet flags
func normalize(s string) string {
	return strings.Replace(s, "_", "-", -1)
}

// pflagRegister adds a flag to local that targets the Value associated with the Flag named globalName in global
func pflagRegister(global, local *pflag.FlagSet, globalName string) {
	if f := global.Lookup(globalName); f != nil {
		f.Name = normalize(f.Name)
		local.AddFlag(f)
	} else {
		panic(fmt.Sprintf("failed to find flag in global flagset (pflag): %s", globalName))
	}
}

// addCredentialProviderFlags adds flags from k8s.io/kubernetes/pkg/credentialprovider
func addCredentialProviderFlags(fs *pflag.FlagSet) {
	// lookup flags in global flag set and re-register the values with our flagset
	global := pflag.CommandLine
	local := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	// TODO(#58034): This is not a static file, so it's not quite as straightforward as --google-json-key.
	// We need to figure out how ACR users can dynamically provide pull credentials before we can deprecate this.
	pflagRegister(global, local, "azure-container-registry-config")

	fs.AddFlagSet(local)
}

// addKlogFlags adds flags from k8s.io/klog
func addKlogFlags(fs *pflag.FlagSet) {
	local := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	klog.InitFlags(local)
	fs.AddGoFlagSet(local)
}
