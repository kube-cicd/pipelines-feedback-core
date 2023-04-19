/*
Copyright Damian Kęska.

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
	"context"

	v1alpha1 "github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakePFConfigs implements PFConfigInterface
type FakePFConfigs struct {
	Fake *FakePipelinesfeedbackV1alpha1
	ns   string
}

var pfconfigsResource = v1alpha1.SchemeGroupVersion.WithResource("pfconfigs")

var pfconfigsKind = v1alpha1.SchemeGroupVersion.WithKind("PFConfig")

// Get takes name of the pFConfig, and returns the corresponding pFConfig object, and an error if there is any.
func (c *FakePFConfigs) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.PFConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(pfconfigsResource, c.ns, name), &v1alpha1.PFConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PFConfig), err
}

// List takes label and field selectors, and returns the list of PFConfigs that match those selectors.
func (c *FakePFConfigs) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.PFConfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(pfconfigsResource, pfconfigsKind, c.ns, opts), &v1alpha1.PFConfigList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.PFConfigList{ListMeta: obj.(*v1alpha1.PFConfigList).ListMeta}
	for _, item := range obj.(*v1alpha1.PFConfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested pFConfigs.
func (c *FakePFConfigs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(pfconfigsResource, c.ns, opts))

}

// Create takes the representation of a pFConfig and creates it.  Returns the server's representation of the pFConfig, and an error, if there is any.
func (c *FakePFConfigs) Create(ctx context.Context, pFConfig *v1alpha1.PFConfig, opts v1.CreateOptions) (result *v1alpha1.PFConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(pfconfigsResource, c.ns, pFConfig), &v1alpha1.PFConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PFConfig), err
}

// Update takes the representation of a pFConfig and updates it. Returns the server's representation of the pFConfig, and an error, if there is any.
func (c *FakePFConfigs) Update(ctx context.Context, pFConfig *v1alpha1.PFConfig, opts v1.UpdateOptions) (result *v1alpha1.PFConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(pfconfigsResource, c.ns, pFConfig), &v1alpha1.PFConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PFConfig), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakePFConfigs) UpdateStatus(ctx context.Context, pFConfig *v1alpha1.PFConfig, opts v1.UpdateOptions) (*v1alpha1.PFConfig, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(pfconfigsResource, "status", c.ns, pFConfig), &v1alpha1.PFConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PFConfig), err
}

// Delete takes name of the pFConfig and deletes it. Returns an error if one occurs.
func (c *FakePFConfigs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(pfconfigsResource, c.ns, name, opts), &v1alpha1.PFConfig{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakePFConfigs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(pfconfigsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.PFConfigList{})
	return err
}

// Patch applies the patch and returns the patched pFConfig.
func (c *FakePFConfigs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.PFConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(pfconfigsResource, c.ns, name, pt, data, subresources...), &v1alpha1.PFConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PFConfig), err
}
