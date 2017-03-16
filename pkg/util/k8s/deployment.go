package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

// Deployment type actions on ingresses in k8s cluster
type Deployment struct {
	kubernetes.Interface
	Namespace string
}

// Get returns deployment object from the k8s by name
func (d *Deployment) Get(name string) (*v1beta1.Deployment, error) {
	return d.
		Extensions().
		Deployments(d.Namespace).
		Get(name)
}

// Exist returns boolean value if deployment exists in k8s
func (d *Deployment) Exist(name string) bool {
	_, err := d.Get(name)
	return err == nil
}

// Apply updates or creates deployment in k8s
func (d *Deployment) Apply(deployment *v1beta1.Deployment) error {
	if d.Exist(deployment.Name) {
		return d.Update(deployment)
	} else {
		return d.Create(deployment)
	}
}

// Update updates existing deployment in k8s
func (d *Deployment) Update(deployment *v1beta1.Deployment) error {
	current, err := d.Get(deployment.Name)
	if err != nil {
		return err
	}
	deployment.ResourceVersion = current.GetResourceVersion()
	if deployment.ObjectMeta.Labels["version"] == "" {
		deployment.ObjectMeta.Labels["version"] = current.ObjectMeta.Labels["version"]
	}

	if len(current.Spec.Template.Spec.Containers) > 0 &&
		len(deployment.Spec.Template.Spec.Containers) > 0 &&
		deployment.Spec.Template.Spec.Containers[0].Image == "" {
		deployment.Spec.Template.Spec.Containers[0].Image = current.Spec.Template.Spec.Containers[0].Image
	}
	_, err = d.
		Extensions().
		Deployments(d.Namespace).
		Update(deployment)
	return err
}

// Create creates new deployment in k8s
func (d *Deployment) Create(deployment *v1beta1.Deployment) error {
	_, err := d.
		Extensions().
		Deployments(d.Namespace).
		Create(deployment)
	return err
}
