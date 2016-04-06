/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

package app

import (
	"github.com/golang/glog"

	"k8s.io/kubernetes/federation/apis/federation"
	"k8s.io/kubernetes/federation/cmd/federated-apiserver/app/options"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/rest"
	"k8s.io/kubernetes/pkg/apimachinery/registered"
	"k8s.io/kubernetes/pkg/genericapiserver"
	"k8s.io/kubernetes/pkg/registry/generic"

	_ "k8s.io/kubernetes/federation/apis/federation/install"
	clusteretcd "k8s.io/kubernetes/federation/registry/cluster/etcd"
	subreplicasetetcd "k8s.io/kubernetes/federation/registry/subreplicaset/etcd"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

func createRESTOptionsOrDie(s *options.APIServer, g *genericapiserver.GenericAPIServer, f genericapiserver.StorageFactory, resource unversioned.GroupResource) generic.RESTOptions {
	storage, err := f.New(resource)
	if err != nil {
		glog.Fatalf("Unable to find storage destination for %v, due to %v", resource, err.Error())
	}
	return generic.RESTOptions{
		Storage:                 storage,
		Decorator:               g.StorageDecorator(),
		DeleteCollectionWorkers: s.DeleteCollectionWorkers,
	}
}

func installFederationAPIs(s *options.APIServer, g *genericapiserver.GenericAPIServer, f genericapiserver.StorageFactory) {
	clusterStorage, clusterStatusStorage := clusteretcd.NewREST(createRESTOptionsOrDie(s, g, f, federation.Resource("clusters")))
	subReplicaSetStorage, subReplicaSetStatusStorage := subreplicasetetcd.NewREST(createRESTOptionsOrDie(s, g, f, federation.Resource("subreplicasets")))
	federationResources := map[string]rest.Storage{
		"clusters":              clusterStorage,
		"clusters/status":       clusterStatusStorage,
		"subreplicasets":        subReplicaSetStorage,
		"subreplicasets/status": subReplicaSetStatusStorage,
	}
	federationGroupMeta := registered.GroupOrDie(federation.GroupName)
	apiGroupInfo := genericapiserver.APIGroupInfo{
		GroupMeta: *federationGroupMeta,
		VersionedResourcesStorageMap: map[string]map[string]rest.Storage{
			"v1alpha1": federationResources,
		},
		OptionsExternalVersion: &registered.GroupOrDie(api.GroupName).GroupVersion,
		Scheme:                 api.Scheme,
		ParameterCodec:         api.ParameterCodec,
		NegotiatedSerializer:   api.Codecs,
	}
	if err := g.InstallAPIGroup(&apiGroupInfo); err != nil {
		glog.Fatalf("Error in registering group versions: %v", err)
	}
}
