/*
Copyright 2015 The Kubernetes Authors All rights reserved.

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

package controller

import (
	client "k8s.io/kubernetes/pkg/client/unversioned"
)

type fakeHost struct {
	kubeClient client.Interface
}

func NewFakeHost(kubeClient client.Interface) *fakeHost {
	return &fakeHost{kubeClient}
}

func (h *fakeHost) GetKubeClient() client.Interface {
	return h.kubeClient
}

var _ Plugin = &FakeControllerPlugin{}

type FakeControllerPlugin struct {
	PluginName string
	Host       Host
}

func (p *FakeControllerPlugin) Init(host Host) {
	p.Host = host
}

func (p *FakeControllerPlugin) Name() string {
	return p.PluginName
}

func (p *FakeControllerPlugin) NewController() (Controller, error) {
	return &FakeController{}, nil
}

type FakeController struct {
	Error   error
	Running bool
}

func (c *FakeController) Run() error {
	if c.Error != nil {
		return c.Error
	}
	c.Running = true
	return nil
}

func (c *FakeController) Stop() error {
	if c.Error != nil {
		return c.Error
	}
	c.Running = false
	return nil
}

func (c *FakeController) Status() (bool, error) {
	return c.Running, c.Error
}
