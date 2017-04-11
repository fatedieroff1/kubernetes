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

package clientset

import (
	glog "github.com/golang/glog"
	innsmouthv1 "k8s.io/apiserver-builder/example/pkg/client/clientset_generated/clientset/typed/innsmouth/v1"
	miskatonicv1beta1 "k8s.io/apiserver-builder/example/pkg/client/clientset_generated/clientset/typed/miskatonic/v1beta1"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	InnsmouthV1() innsmouthv1.InnsmouthV1Interface
	// Deprecated: please explicitly pick a version if possible.
	Innsmouth() innsmouthv1.InnsmouthV1Interface
	MiskatonicV1beta1() miskatonicv1beta1.MiskatonicV1beta1Interface
	// Deprecated: please explicitly pick a version if possible.
	Miskatonic() miskatonicv1beta1.MiskatonicV1beta1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	*innsmouthv1.InnsmouthV1Client
	*miskatonicv1beta1.MiskatonicV1beta1Client
}

// InnsmouthV1 retrieves the InnsmouthV1Client
func (c *Clientset) InnsmouthV1() innsmouthv1.InnsmouthV1Interface {
	if c == nil {
		return nil
	}
	return c.InnsmouthV1Client
}

// Deprecated: Innsmouth retrieves the default version of InnsmouthClient.
// Please explicitly pick a version.
func (c *Clientset) Innsmouth() innsmouthv1.InnsmouthV1Interface {
	if c == nil {
		return nil
	}
	return c.InnsmouthV1Client
}

// MiskatonicV1beta1 retrieves the MiskatonicV1beta1Client
func (c *Clientset) MiskatonicV1beta1() miskatonicv1beta1.MiskatonicV1beta1Interface {
	if c == nil {
		return nil
	}
	return c.MiskatonicV1beta1Client
}

// Deprecated: Miskatonic retrieves the default version of MiskatonicClient.
// Please explicitly pick a version.
func (c *Clientset) Miskatonic() miskatonicv1beta1.MiskatonicV1beta1Interface {
	if c == nil {
		return nil
	}
	return c.MiskatonicV1beta1Client
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.InnsmouthV1Client, err = innsmouthv1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.MiskatonicV1beta1Client, err = miskatonicv1beta1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		glog.Errorf("failed to create the DiscoveryClient: %v", err)
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.InnsmouthV1Client = innsmouthv1.NewForConfigOrDie(c)
	cs.MiskatonicV1beta1Client = miskatonicv1beta1.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.InnsmouthV1Client = innsmouthv1.New(c)
	cs.MiskatonicV1beta1Client = miskatonicv1beta1.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
