//=======================================================================
// Copyright (c) 2022 Aarna Networks, Inc.
// All rights reserved.
// ======================================================================
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//           http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// ========================================================================

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	v1alpha1 "edgemetricscollector/pkg/apis/edgemetricscollector/v1alpha1"
	scheme "edgemetricscollector/pkg/generated/clientset/versioned/scheme"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MetricsCollectorsGetter has a method to return a MetricsCollectorInterface.
// A group's client should implement this interface.
type MetricsCollectorsGetter interface {
	MetricsCollectors(namespace string) MetricsCollectorInterface
}

// MetricsCollectorInterface has methods to work with MetricsCollector resources.
type MetricsCollectorInterface interface {
	Create(ctx context.Context, metricsCollector *v1alpha1.MetricsCollector, opts v1.CreateOptions) (*v1alpha1.MetricsCollector, error)
	Update(ctx context.Context, metricsCollector *v1alpha1.MetricsCollector, opts v1.UpdateOptions) (*v1alpha1.MetricsCollector, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.MetricsCollector, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.MetricsCollectorList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.MetricsCollector, err error)
	MetricsCollectorExpansion
}

// metricsCollectors implements MetricsCollectorInterface
type metricsCollectors struct {
	client rest.Interface
	ns     string
}

// newMetricsCollectors returns a MetricsCollectors
func newMetricsCollectors(c *EdgemetricscollectorV1alpha1Client, namespace string) *metricsCollectors {
	return &metricsCollectors{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the metricsCollector, and returns the corresponding metricsCollector object, and an error if there is any.
func (c *metricsCollectors) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.MetricsCollector, err error) {
	result = &v1alpha1.MetricsCollector{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("metricscollectors").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MetricsCollectors that match those selectors.
func (c *metricsCollectors) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.MetricsCollectorList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.MetricsCollectorList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("metricscollectors").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested metricsCollectors.
func (c *metricsCollectors) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("metricscollectors").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a metricsCollector and creates it.  Returns the server's representation of the metricsCollector, and an error, if there is any.
func (c *metricsCollectors) Create(ctx context.Context, metricsCollector *v1alpha1.MetricsCollector, opts v1.CreateOptions) (result *v1alpha1.MetricsCollector, err error) {
	result = &v1alpha1.MetricsCollector{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("metricscollectors").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(metricsCollector).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a metricsCollector and updates it. Returns the server's representation of the metricsCollector, and an error, if there is any.
func (c *metricsCollectors) Update(ctx context.Context, metricsCollector *v1alpha1.MetricsCollector, opts v1.UpdateOptions) (result *v1alpha1.MetricsCollector, err error) {
	result = &v1alpha1.MetricsCollector{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("metricscollectors").
		Name(metricsCollector.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(metricsCollector).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the metricsCollector and deletes it. Returns an error if one occurs.
func (c *metricsCollectors) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("metricscollectors").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *metricsCollectors) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("metricscollectors").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched metricsCollector.
func (c *metricsCollectors) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.MetricsCollector, err error) {
	result = &v1alpha1.MetricsCollector{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("metricscollectors").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}