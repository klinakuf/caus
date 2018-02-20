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

// This file was automatically generated by informer-gen

package v1

import (
	caus_rss_uni_stuttgart_de_v1 "github.com/klinakuf/caus/pkg/apis/caus.rss.uni-stuttgart.de/v1"
	versioned "github.com/klinakuf/caus/pkg/client/clientset/versioned"
	internalinterfaces "github.com/klinakuf/caus/pkg/client/informers/externalversions/internalinterfaces"
	v1 "github.com/klinakuf/caus/pkg/client/listers/caus/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	time "time"
)

// ElasticityInformer provides access to a shared informer and lister for
// Elasticities.
type ElasticityInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.ElasticityLister
}

type elasticityInformer struct {
	factory internalinterfaces.SharedInformerFactory
}

// NewElasticityInformer constructs a new informer for Elasticity type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewElasticityInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				return client.CausV1().Elasticities(namespace).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				return client.CausV1().Elasticities(namespace).Watch(options)
			},
		},
		&caus_rss_uni_stuttgart_de_v1.Elasticity{},
		resyncPeriod,
		indexers,
	)
}

func defaultElasticityInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewElasticityInformer(client, meta_v1.NamespaceAll, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
}

func (f *elasticityInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&caus_rss_uni_stuttgart_de_v1.Elasticity{}, defaultElasticityInformer)
}

func (f *elasticityInformer) Lister() v1.ElasticityLister {
	return v1.NewElasticityLister(f.Informer().GetIndexer())
}
