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
package fake

import (
	caus_rss_uni_stuttgart_de_v1 "github.com/klinakuf/caus/pkg/apis/caus.rss.uni-stuttgart.de/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeElasticities implements ElasticityInterface
type FakeElasticities struct {
	Fake *FakeCausV1
	ns   string
}

var elasticitiesResource = schema.GroupVersionResource{Group: "caus", Version: "v1", Resource: "elasticities"}

var elasticitiesKind = schema.GroupVersionKind{Group: "caus", Version: "v1", Kind: "Elasticity"}

// Get takes name of the elasticity, and returns the corresponding elasticity object, and an error if there is any.
func (c *FakeElasticities) Get(name string, options v1.GetOptions) (result *caus_rss_uni_stuttgart_de_v1.Elasticity, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(elasticitiesResource, c.ns, name), &caus_rss_uni_stuttgart_de_v1.Elasticity{})

	if obj == nil {
		return nil, err
	}
	return obj.(*caus_rss_uni_stuttgart_de_v1.Elasticity), err
}

// List takes label and field selectors, and returns the list of Elasticities that match those selectors.
func (c *FakeElasticities) List(opts v1.ListOptions) (result *caus_rss_uni_stuttgart_de_v1.ElasticityList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(elasticitiesResource, elasticitiesKind, c.ns, opts), &caus_rss_uni_stuttgart_de_v1.ElasticityList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &caus_rss_uni_stuttgart_de_v1.ElasticityList{}
	for _, item := range obj.(*caus_rss_uni_stuttgart_de_v1.ElasticityList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested elasticities.
func (c *FakeElasticities) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(elasticitiesResource, c.ns, opts))

}

// Create takes the representation of a elasticity and creates it.  Returns the server's representation of the elasticity, and an error, if there is any.
func (c *FakeElasticities) Create(elasticity *caus_rss_uni_stuttgart_de_v1.Elasticity) (result *caus_rss_uni_stuttgart_de_v1.Elasticity, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(elasticitiesResource, c.ns, elasticity), &caus_rss_uni_stuttgart_de_v1.Elasticity{})

	if obj == nil {
		return nil, err
	}
	return obj.(*caus_rss_uni_stuttgart_de_v1.Elasticity), err
}

// Update takes the representation of a elasticity and updates it. Returns the server's representation of the elasticity, and an error, if there is any.
func (c *FakeElasticities) Update(elasticity *caus_rss_uni_stuttgart_de_v1.Elasticity) (result *caus_rss_uni_stuttgart_de_v1.Elasticity, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(elasticitiesResource, c.ns, elasticity), &caus_rss_uni_stuttgart_de_v1.Elasticity{})

	if obj == nil {
		return nil, err
	}
	return obj.(*caus_rss_uni_stuttgart_de_v1.Elasticity), err
}

// Delete takes name of the elasticity and deletes it. Returns an error if one occurs.
func (c *FakeElasticities) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(elasticitiesResource, c.ns, name), &caus_rss_uni_stuttgart_de_v1.Elasticity{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeElasticities) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(elasticitiesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &caus_rss_uni_stuttgart_de_v1.ElasticityList{})
	return err
}

// Patch applies the patch and returns the patched elasticity.
func (c *FakeElasticities) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *caus_rss_uni_stuttgart_de_v1.Elasticity, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(elasticitiesResource, c.ns, name, data, subresources...), &caus_rss_uni_stuttgart_de_v1.Elasticity{})

	if obj == nil {
		return nil, err
	}
	return obj.(*caus_rss_uni_stuttgart_de_v1.Elasticity), err
}
