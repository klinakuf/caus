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
package v1

import (
	v1 "github.com/klinakuf/caus/pkg/apis/caus.rss.uni-stuttgart.de/v1"
	scheme "github.com/klinakuf/caus/pkg/client/clientset/versioned/scheme"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ElasticitiesGetter has a method to return a ElasticityInterface.
// A group's client should implement this interface.
type ElasticitiesGetter interface {
	Elasticities(namespace string) ElasticityInterface
}

// ElasticityInterface has methods to work with Elasticity resources.
type ElasticityInterface interface {
	Create(*v1.Elasticity) (*v1.Elasticity, error)
	Update(*v1.Elasticity) (*v1.Elasticity, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.Elasticity, error)
	List(opts meta_v1.ListOptions) (*v1.ElasticityList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Elasticity, err error)
	ElasticityExpansion
}

// elasticities implements ElasticityInterface
type elasticities struct {
	client rest.Interface
	ns     string
}

// newElasticities returns a Elasticities
func newElasticities(c *CausV1Client, namespace string) *elasticities {
	return &elasticities{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the elasticity, and returns the corresponding elasticity object, and an error if there is any.
func (c *elasticities) Get(name string, options meta_v1.GetOptions) (result *v1.Elasticity, err error) {
	result = &v1.Elasticity{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("elasticities").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Elasticities that match those selectors.
func (c *elasticities) List(opts meta_v1.ListOptions) (result *v1.ElasticityList, err error) {
	result = &v1.ElasticityList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("elasticities").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested elasticities.
func (c *elasticities) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("elasticities").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a elasticity and creates it.  Returns the server's representation of the elasticity, and an error, if there is any.
func (c *elasticities) Create(elasticity *v1.Elasticity) (result *v1.Elasticity, err error) {
	result = &v1.Elasticity{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("elasticities").
		Body(elasticity).
		Do().
		Into(result)
	return
}

// Update takes the representation of a elasticity and updates it. Returns the server's representation of the elasticity, and an error, if there is any.
func (c *elasticities) Update(elasticity *v1.Elasticity) (result *v1.Elasticity, err error) {
	result = &v1.Elasticity{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("elasticities").
		Name(elasticity.Name).
		Body(elasticity).
		Do().
		Into(result)
	return
}

// Delete takes name of the elasticity and deletes it. Returns an error if one occurs.
func (c *elasticities) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("elasticities").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *elasticities) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("elasticities").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched elasticity.
func (c *elasticities) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Elasticity, err error) {
	result = &v1.Elasticity{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("elasticities").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
