// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "gitlab.com/project-emco/core/emco-base/src/monitor/pkg/apis/k8splugin/v1alpha1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeResourceBundleStates implements ResourceBundleStateInterface
type FakeResourceBundleStates struct {
	Fake *FakeK8spluginV1alpha1
	ns   string
}

var resourcebundlestatesResource = schema.GroupVersionResource{Group: "k8splugin.io", Version: "v1alpha1", Resource: "resourcebundlestates"}

var resourcebundlestatesKind = schema.GroupVersionKind{Group: "k8splugin.io", Version: "v1alpha1", Kind: "ResourceBundleState"}

// Get takes name of the resourceBundleState, and returns the corresponding resourceBundleState object, and an error if there is any.
func (c *FakeResourceBundleStates) Get(name string, options v1.GetOptions) (result *v1alpha1.ResourceBundleState, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(resourcebundlestatesResource, c.ns, name), &v1alpha1.ResourceBundleState{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ResourceBundleState), err
}

// List takes label and field selectors, and returns the list of ResourceBundleStates that match those selectors.
func (c *FakeResourceBundleStates) List(opts v1.ListOptions) (result *v1alpha1.ResourceBundleStateList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(resourcebundlestatesResource, resourcebundlestatesKind, c.ns, opts), &v1alpha1.ResourceBundleStateList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ResourceBundleStateList{ListMeta: obj.(*v1alpha1.ResourceBundleStateList).ListMeta}
	for _, item := range obj.(*v1alpha1.ResourceBundleStateList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested resourceBundleStates.
func (c *FakeResourceBundleStates) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(resourcebundlestatesResource, c.ns, opts))

}

// Create takes the representation of a resourceBundleState and creates it.  Returns the server's representation of the resourceBundleState, and an error, if there is any.
func (c *FakeResourceBundleStates) Create(resourceBundleState *v1alpha1.ResourceBundleState) (result *v1alpha1.ResourceBundleState, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(resourcebundlestatesResource, c.ns, resourceBundleState), &v1alpha1.ResourceBundleState{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ResourceBundleState), err
}

// Update takes the representation of a resourceBundleState and updates it. Returns the server's representation of the resourceBundleState, and an error, if there is any.
func (c *FakeResourceBundleStates) Update(resourceBundleState *v1alpha1.ResourceBundleState) (result *v1alpha1.ResourceBundleState, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(resourcebundlestatesResource, c.ns, resourceBundleState), &v1alpha1.ResourceBundleState{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ResourceBundleState), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeResourceBundleStates) UpdateStatus(resourceBundleState *v1alpha1.ResourceBundleState) (*v1alpha1.ResourceBundleState, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(resourcebundlestatesResource, "status", c.ns, resourceBundleState), &v1alpha1.ResourceBundleState{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ResourceBundleState), err
}

// Delete takes name of the resourceBundleState and deletes it. Returns an error if one occurs.
func (c *FakeResourceBundleStates) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(resourcebundlestatesResource, c.ns, name), &v1alpha1.ResourceBundleState{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeResourceBundleStates) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(resourcebundlestatesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.ResourceBundleStateList{})
	return err
}

// Patch applies the patch and returns the patched resourceBundleState.
func (c *FakeResourceBundleStates) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ResourceBundleState, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(resourcebundlestatesResource, c.ns, name, pt, data, subresources...), &v1alpha1.ResourceBundleState{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ResourceBundleState), err
}
