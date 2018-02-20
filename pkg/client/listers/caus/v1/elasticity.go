/*
Copyright 2018 klinaku@informatik.uni-stuttgart

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
	v1 "github.com/klinakuf/caus/pkg/apis/caus.rss.uni-stuttgart.de/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ElasticityLister helps list Elasticities.
type ElasticityLister interface {
	// List lists all Elasticities in the indexer.
	List(selector labels.Selector) (ret []*v1.Elasticity, err error)
	// Elasticities returns an object that can list and get Elasticities.
	Elasticities(namespace string) ElasticityNamespaceLister
	ElasticityListerExpansion
}

// elasticityLister implements the ElasticityLister interface.
type elasticityLister struct {
	indexer cache.Indexer
}

// NewElasticityLister returns a new ElasticityLister.
func NewElasticityLister(indexer cache.Indexer) ElasticityLister {
	return &elasticityLister{indexer: indexer}
}

// List lists all Elasticities in the indexer.
func (s *elasticityLister) List(selector labels.Selector) (ret []*v1.Elasticity, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Elasticity))
	})
	return ret, err
}

// Elasticities returns an object that can list and get Elasticities.
func (s *elasticityLister) Elasticities(namespace string) ElasticityNamespaceLister {
	return elasticityNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ElasticityNamespaceLister helps list and get Elasticities.
type ElasticityNamespaceLister interface {
	// List lists all Elasticities in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.Elasticity, err error)
	// Get retrieves the Elasticity from the indexer for a given namespace and name.
	Get(name string) (*v1.Elasticity, error)
	ElasticityNamespaceListerExpansion
}

// elasticityNamespaceLister implements the ElasticityNamespaceLister
// interface.
type elasticityNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Elasticities in the indexer for a given namespace.
func (s elasticityNamespaceLister) List(selector labels.Selector) (ret []*v1.Elasticity, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Elasticity))
	})
	return ret, err
}

// Get retrieves the Elasticity from the indexer for a given namespace and name.
func (s elasticityNamespaceLister) Get(name string) (*v1.Elasticity, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("elasticity"), name)
	}
	return obj.(*v1.Elasticity), nil
}
