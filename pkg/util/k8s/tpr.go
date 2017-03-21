package k8s

import (
	extensions "github.com/pearsontechnology/environment-operator/pkg/k8_extensions"
	"k8s.io/client-go/rest"
)

// ThirdPartyResource represents TPR API on the cluster
type ThirdPartyResource struct {
	rest.Interface

	Namespace string
	Type      string
}

// Get retrieves PrsnExternalResource from the k8s using name
func (client *ThirdPartyResource) Get(name string) (*extensions.PrsnExternalResource, error) {
	var rsc *extensions.PrsnExternalResource

	err := client.Interface.Get().
		Resource(client.Type).
		Namespace(client.Namespace).
		Name(name).
		Do().Into(rsc)
	return rsc, err
}

// Exist checks if named resource exist in k8s cluster
func (client *ThirdPartyResource) Exist(name string) bool {
	_, err := client.Get(name)
	return err == nil
}

// Apply creates or updates PrsnExternalResource in k8s
func (client *ThirdPartyResource) Apply(resource *extensions.PrsnExternalResource) error {
	if client.Exist(resource.ObjectMeta.Name) {
		return client.Update(resource)
	}
	return client.Create(resource)
}

// Create creates given tpr in
func (client *ThirdPartyResource) Create(resource *extensions.PrsnExternalResource) error {
	return nil
}

// Update updates existing resource in k8s
func (client *ThirdPartyResource) Update(resource *extensions.PrsnExternalResource) error {
	var result *extensions.PrsnExternalResource
	return client.Interface.Put().
		Resource(client.Type).
		Namespace(client.Namespace).
		Body(resource).
		Do().Into(result)
}

// Destroy deletes named resource
func (client *ThirdPartyResource) Destroy(name string) error {
	var result *extensions.PrsnExternalResource
	return client.Interface.Delete().
		Resource(client.Type).
		Namespace(client.Namespace).
		Name(name).Do().Into(result)
}
