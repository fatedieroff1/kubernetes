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

import "k8s.io/apimachinery/pkg/runtime"

const (
	DefaultServiceDNSDomain  = "cluster.local"
	DefaultServicesSubnet    = "10.96.0.0/12"
	DefaultKubernetesVersion = "latest"
	// This is only for clusters without internet, were the latest stable version can't be determined
	DefaultKubernetesFallbackVersion = "v1.6.0-alpha.1"
	DefaultAPIBindAddress            = "0.0.0.0"
	DefaultAPIBindPort               = 6443
	DefaultDiscoveryBindPort         = 9898
	DefaultAuthorizationMode         = "RBAC"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	RegisterDefaults(scheme)
	return scheme.AddDefaultingFuncs(
		SetDefaults_MasterConfiguration,
	)
}

func SetDefaults_MasterConfiguration(obj *MasterConfiguration) {
	if obj.KubernetesVersion == "" {
		obj.KubernetesVersion = DefaultKubernetesVersion
	}

	if obj.API.BindPort == 0 {
		obj.API.BindPort = DefaultAPIBindPort
	}

	if obj.API.BindAddress == "" {
		obj.API.BindAddress = DefaultAPIBindAddress
	}

	if obj.Networking.ServiceSubnet == "" {
		obj.Networking.ServiceSubnet = DefaultServicesSubnet
	}

	if obj.Networking.DNSDomain == "" {
		obj.Networking.DNSDomain = DefaultServiceDNSDomain
	}

	if obj.Discovery.Token == nil && obj.Discovery.File == nil && obj.Discovery.HTTPS == nil {
		obj.Discovery.Token = &TokenDiscovery{}
	}

	if obj.AuthorizationMode == "" {
		obj.AuthorizationMode = DefaultAuthorizationMode
	}
}
