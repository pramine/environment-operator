package k8s

import "k8s.io/client-go/kubernetes"
import "k8s.io/client-go/pkg/api/v1"

// Service type actions on pvcs in k8s cluster
type Service struct {
	kubernetes.Interface
	Namespace string
}

// Get returns service object from the k8s by name
func (client *Service) Get(name string) (interface{}, error) {
	return client.Core().Services(client.Namespace).Get(name)
}

// Exist returns boolean value if pvc exists in k8s
func (client *Service) Exist(name string) bool {
	_, err := client.Get(name)
	return err == nil
}

// Apply updates or creates service in k8s
func (client *Service) Apply(i interface{}) error {
	resource := i.(*v1.Service)
	if client.Exist(resource.Name) {
		return client.Update(resource)
	}
	return client.Create(resource)
}

// Create creates new service in k8s
func (client *Service) Create(i interface{}) error {
	resource := i.(*v1.Service)
	_, err := client.
		Core().
		Services(client.Namespace).
		Create(resource)
	return err
}

// Update updates existing service in k8s
func (client *Service) Update(i interface{}) error {
	resource := i.(*v1.Service)
	ci, err := client.Get(resource.Name)
	if err != nil {
		return err
	}
	current := ci.(*v1.Service)
	resource.ResourceVersion = current.GetResourceVersion()

	_, err = client.
		Core().
		Services(client.Namespace).
		Update(resource)
	return err
}

// Destroy deletes service from the k8 cluster
func (client *Service) Destroy(name string) error {
	options := &v1.DeleteOptions{}
	return client.Core().Services(client.Namespace).Delete(name, options)
}
