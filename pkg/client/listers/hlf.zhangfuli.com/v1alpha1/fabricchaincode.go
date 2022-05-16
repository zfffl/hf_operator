/*
Copyright 2021 zhangfuli.

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
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "HFOperator/api/hlf.zhangfuli.com/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// FabricChaincodeLister helps list FabricChaincodes.
type FabricChaincodeLister interface {
	// List lists all FabricChaincodes in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.FabricChaincode, err error)
	// FabricChaincodes returns an object that can list and get FabricChaincodes.
	FabricChaincodes(namespace string) FabricChaincodeNamespaceLister
	FabricChaincodeListerExpansion
}

// fabricChaincodeLister implements the FabricChaincodeLister interface.
type fabricChaincodeLister struct {
	indexer cache.Indexer
}

// NewFabricChaincodeLister returns a new FabricChaincodeLister.
func NewFabricChaincodeLister(indexer cache.Indexer) FabricChaincodeLister {
	return &fabricChaincodeLister{indexer: indexer}
}

// List lists all FabricChaincodes in the indexer.
func (s *fabricChaincodeLister) List(selector labels.Selector) (ret []*v1alpha1.FabricChaincode, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FabricChaincode))
	})
	return ret, err
}

// FabricChaincodes returns an object that can list and get FabricChaincodes.
func (s *fabricChaincodeLister) FabricChaincodes(namespace string) FabricChaincodeNamespaceLister {
	return fabricChaincodeNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// FabricChaincodeNamespaceLister helps list and get FabricChaincodes.
type FabricChaincodeNamespaceLister interface {
	// List lists all FabricChaincodes in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.FabricChaincode, err error)
	// Get retrieves the FabricChaincode from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.FabricChaincode, error)
	FabricChaincodeNamespaceListerExpansion
}

// fabricChaincodeNamespaceLister implements the FabricChaincodeNamespaceLister
// interface.
type fabricChaincodeNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all FabricChaincodes in the indexer for a given namespace.
func (s fabricChaincodeNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.FabricChaincode, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FabricChaincode))
	})
	return ret, err
}

// Get retrieves the FabricChaincode from the indexer for a given namespace and name.
func (s fabricChaincodeNamespaceLister) Get(name string) (*v1alpha1.FabricChaincode, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("fabricchaincode"), name)
	}
	return obj.(*v1alpha1.FabricChaincode), nil
}