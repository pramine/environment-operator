package kubernetes

import (
	"fmt"

	"k8s.io/client-go/pkg/api/v1"
)

func (w *Wrapper) volumeClaimFromName(namespace, name string) (v1.PersistentVolumeClaim, error) {
	claims, err := w.PersistentVolumeClaims(namespace)
	if err != nil {
		return v1.PersistentVolumeClaim{}, err
	}

	for _, claim := range claims {
		if claim.Name == name {
			return claim, nil
		}
	}
	return v1.PersistentVolumeClaim{}, fmt.Errorf("Persistent volume claim %s not found", name)
}
