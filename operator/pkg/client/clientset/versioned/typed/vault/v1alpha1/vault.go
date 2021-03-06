// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/banzaicloud/bank-vaults/operator/pkg/apis/vault/v1alpha1"
	scheme "github.com/banzaicloud/bank-vaults/operator/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// VaultsGetter has a method to return a VaultInterface.
// A group's client should implement this interface.
type VaultsGetter interface {
	Vaults(namespace string) VaultInterface
}

// VaultInterface has methods to work with Vault resources.
type VaultInterface interface {
	Create(*v1alpha1.Vault) (*v1alpha1.Vault, error)
	Update(*v1alpha1.Vault) (*v1alpha1.Vault, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Vault, error)
	List(opts v1.ListOptions) (*v1alpha1.VaultList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Vault, err error)
	VaultExpansion
}

// vaults implements VaultInterface
type vaults struct {
	client rest.Interface
	ns     string
}

// newVaults returns a Vaults
func newVaults(c *VaultV1alpha1Client, namespace string) *vaults {
	return &vaults{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the vault, and returns the corresponding vault object, and an error if there is any.
func (c *vaults) Get(name string, options v1.GetOptions) (result *v1alpha1.Vault, err error) {
	result = &v1alpha1.Vault{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("vaults").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Vaults that match those selectors.
func (c *vaults) List(opts v1.ListOptions) (result *v1alpha1.VaultList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.VaultList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("vaults").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested vaults.
func (c *vaults) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("vaults").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a vault and creates it.  Returns the server's representation of the vault, and an error, if there is any.
func (c *vaults) Create(vault *v1alpha1.Vault) (result *v1alpha1.Vault, err error) {
	result = &v1alpha1.Vault{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("vaults").
		Body(vault).
		Do().
		Into(result)
	return
}

// Update takes the representation of a vault and updates it. Returns the server's representation of the vault, and an error, if there is any.
func (c *vaults) Update(vault *v1alpha1.Vault) (result *v1alpha1.Vault, err error) {
	result = &v1alpha1.Vault{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("vaults").
		Name(vault.Name).
		Body(vault).
		Do().
		Into(result)
	return
}

// Delete takes name of the vault and deletes it. Returns an error if one occurs.
func (c *vaults) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("vaults").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *vaults) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("vaults").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched vault.
func (c *vaults) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Vault, err error) {
	result = &v1alpha1.Vault{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("vaults").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
