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

// This file was automatically generated by lister-gen

package v1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	v1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

// PersistentVolumeLister helps list PersistentVolumes.
type PersistentVolumeLister interface {
	// List lists all PersistentVolumes in the indexer.
	List(selector labels.Selector) (ret []*v1.PersistentVolume, err error)
	// Get retrieves the PersistentVolume from the index for a given name.
	Get(name string) (*v1.PersistentVolume, error)
	PersistentVolumeListerExpansion
}

// persistentVolumeLister implements the PersistentVolumeLister interface.
type persistentVolumeLister struct {
	indexer cache.Indexer
}

// NewPersistentVolumeLister returns a new PersistentVolumeLister.
func NewPersistentVolumeLister(indexer cache.Indexer) PersistentVolumeLister {
	return &persistentVolumeLister{indexer: indexer}
}

// List lists all PersistentVolumes in the indexer.
func (s *persistentVolumeLister) List(selector labels.Selector) (ret []*v1.PersistentVolume, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.PersistentVolume))
	})
	return ret, err
}

// Get retrieves the PersistentVolume from the index for a given name.
func (s *persistentVolumeLister) Get(name string) (*v1.PersistentVolume, error) {
	key := &v1.PersistentVolume{ObjectMeta: meta_v1.ObjectMeta{Name: name}}
	obj, exists, err := s.indexer.Get(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("persistentvolume"), name)
	}
	return obj.(*v1.PersistentVolume), nil
}
