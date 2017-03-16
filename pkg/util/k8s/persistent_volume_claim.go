package k8s

import "k8s.io/client-go/kubernetes"
import "k8s.io/client-go/pkg/api/v1"

// PersistentVolumeClaim type actions on pvcs in k8s cluster
type PersistentVolumeClaim struct {
	kubernetes.Interface
	Namespace string
}

func (client *PersistentVolumeClaim) Get(name string) (*v1.PersistentVolumeClaim, error) {
	return nil, nil
}

func (client *PersistentVolumeClaim) Exist(name string) bool {
	return true
}

func (client *PersistentVolumeClaim) Create(resource *v1.PersistentVolumeClaim) error {
	return nil
}

func (client *PersistentVolumeClaim) Update(resource *v1.PersistentVolumeClaim) error {
	return nil
}

func (client *PersistentVolumeClaim) Apply(resource *v1.PersistentVolumeClaim) error {
	return nil
}
