/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	scheme "workspace/ligato/nsm2nsm/pkg/client/clientset/versioned/scheme"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1 "workspace/ligato/nsm2nsm/pkg/apis/sbezverk.io/v1"
)

// ServerEndpointsGetter has a method to return a ServerEndpointInterface.
// A group's client should implement this interface.
type ServerEndpointsGetter interface {
	ServerEndpoints(namespace string) ServerEndpointInterface
}

// ServerEndpointInterface has methods to work with ServerEndpoint resources.
type ServerEndpointInterface interface {
	Create(*v1.ServerEndpoint) (*v1.ServerEndpoint, error)
	Update(*v1.ServerEndpoint) (*v1.ServerEndpoint, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.ServerEndpoint, error)
	List(opts metav1.ListOptions) (*v1.ServerEndpointList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ServerEndpoint, err error)
	ServerEndpointExpansion
}

// serverEndpoints implements ServerEndpointInterface
type serverEndpoints struct {
	client rest.Interface
	ns     string
}

// newServerEndpoints returns a ServerEndpoints
func newServerEndpoints(c *SbezverkV1Client, namespace string) *serverEndpoints {
	return &serverEndpoints{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the serverEndpoint, and returns the corresponding serverEndpoint object, and an error if there is any.
func (c *serverEndpoints) Get(name string, options metav1.GetOptions) (result *v1.ServerEndpoint, err error) {
	result = &v1.ServerEndpoint{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("serverendpoints").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ServerEndpoints that match those selectors.
func (c *serverEndpoints) List(opts metav1.ListOptions) (result *v1.ServerEndpointList, err error) {
	result = &v1.ServerEndpointList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("serverendpoints").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested serverEndpoints.
func (c *serverEndpoints) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("serverendpoints").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a serverEndpoint and creates it.  Returns the server's representation of the serverEndpoint, and an error, if there is any.
func (c *serverEndpoints) Create(serverEndpoint *v1.ServerEndpoint) (result *v1.ServerEndpoint, err error) {
	result = &v1.ServerEndpoint{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("serverendpoints").
		Body(serverEndpoint).
		Do().
		Into(result)
	return
}

// Update takes the representation of a serverEndpoint and updates it. Returns the server's representation of the serverEndpoint, and an error, if there is any.
func (c *serverEndpoints) Update(serverEndpoint *v1.ServerEndpoint) (result *v1.ServerEndpoint, err error) {
	result = &v1.ServerEndpoint{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("serverendpoints").
		Name(serverEndpoint.Name).
		Body(serverEndpoint).
		Do().
		Into(result)
	return
}

// Delete takes name of the serverEndpoint and deletes it. Returns an error if one occurs.
func (c *serverEndpoints) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("serverendpoints").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *serverEndpoints) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("serverendpoints").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched serverEndpoint.
func (c *serverEndpoints) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ServerEndpoint, err error) {
	result = &v1.ServerEndpoint{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("serverendpoints").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}