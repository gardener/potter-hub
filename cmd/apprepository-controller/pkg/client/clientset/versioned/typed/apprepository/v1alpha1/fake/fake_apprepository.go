/*
Copyright 2018 Bitnami.

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

package fake

import (
	v1alpha1 "github.wdf.sap.corp/kubernetes/hub/cmd/apprepository-controller/pkg/apis/apprepository/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeAppRepositories implements AppRepositoryInterface
type FakeAppRepositories struct {
	Fake *FakeKubeappsV1alpha1
	ns   string
}

var apprepositoriesResource = schema.GroupVersionResource{Group: "kubeapps.com", Version: "v1alpha1", Resource: "apprepositories"}

var apprepositoriesKind = schema.GroupVersionKind{Group: "kubeapps.com", Version: "v1alpha1", Kind: "AppRepository"}

// Get takes name of the appRepository, and returns the corresponding appRepository object, and an error if there is any.
func (c *FakeAppRepositories) Get(name string, options v1.GetOptions) (result *v1alpha1.AppRepository, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(apprepositoriesResource, c.ns, name), &v1alpha1.AppRepository{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AppRepository), err
}

// List takes label and field selectors, and returns the list of AppRepositories that match those selectors.
func (c *FakeAppRepositories) List(opts v1.ListOptions) (result *v1alpha1.AppRepositoryList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(apprepositoriesResource, apprepositoriesKind, c.ns, opts), &v1alpha1.AppRepositoryList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.AppRepositoryList{}
	for _, item := range obj.(*v1alpha1.AppRepositoryList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested appRepositories.
func (c *FakeAppRepositories) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(apprepositoriesResource, c.ns, opts))

}

// Create takes the representation of a appRepository and creates it.  Returns the server's representation of the appRepository, and an error, if there is any.
func (c *FakeAppRepositories) Create(appRepository *v1alpha1.AppRepository) (result *v1alpha1.AppRepository, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(apprepositoriesResource, c.ns, appRepository), &v1alpha1.AppRepository{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AppRepository), err
}

// Update takes the representation of a appRepository and updates it. Returns the server's representation of the appRepository, and an error, if there is any.
func (c *FakeAppRepositories) Update(appRepository *v1alpha1.AppRepository) (result *v1alpha1.AppRepository, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(apprepositoriesResource, c.ns, appRepository), &v1alpha1.AppRepository{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AppRepository), err
}

// Delete takes name of the appRepository and deletes it. Returns an error if one occurs.
func (c *FakeAppRepositories) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(apprepositoriesResource, c.ns, name), &v1alpha1.AppRepository{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAppRepositories) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(apprepositoriesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.AppRepositoryList{})
	return err
}

// Patch applies the patch and returns the patched appRepository.
func (c *FakeAppRepositories) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.AppRepository, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(apprepositoriesResource, c.ns, name, pt, data, subresources...), &v1alpha1.AppRepository{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AppRepository), err
}
