package translator

import (
	"k8s.io/client-go/pkg/api/resource"
	"k8s.io/client-go/pkg/api/v1"
)

// In this file we define database specific containers, volumes and annotations used to generate statefulset objects

// Couchbase

// CbContainers generates a slice containing the couchbase containers required in a couchbase statefulset
func (w *KubeMapper) CbContainers() ([]v1.Container, error) {
	mounts, err := w.volumeMounts()
	if err != nil {
		return nil, err
	}
	res, err := w.resources()
	if err != nil {
		return nil, err
	}

	cb := []v1.Container{
		{
			Name:            "couchbase",
			Image:           "couchbase/server:enterprise-" + w.BiteService.Version,
			ImagePullPolicy: v1.PullIfNotPresent,
			Ports: []v1.ContainerPort{
				{ContainerPort: int32(8091), Name: "cb-admin"},
				{ContainerPort: int32(8092), Name: "cb-views"},
				{ContainerPort: int32(8093), Name: "cb-queries"},
				{ContainerPort: int32(8094), Name: "cb-search"},
				{ContainerPort: int32(9100), Name: "cb-int-ind-ad"},
				{ContainerPort: int32(9101), Name: "cb-int-ind-sc"},
				{ContainerPort: int32(9102), Name: "cb-int-ind-ht"},
				{ContainerPort: int32(9103), Name: "cb-int-ind-in"},
				{ContainerPort: int32(9104), Name: "cb-int-ind-ca"},
				{ContainerPort: int32(9105), Name: "cb-int-ind-ma"},
				{ContainerPort: int32(9998), Name: "cb-int-rest"},
				{ContainerPort: int32(9999), Name: "cb-int-gsi"},
				{ContainerPort: int32(11207), Name: "cb-memc-ssl"},
				{ContainerPort: int32(11209), Name: "cb-int-bu"},
				{ContainerPort: int32(11210), Name: "cb-moxi"},
				{ContainerPort: int32(11211), Name: "cb-memc"},
				{ContainerPort: int32(11214), Name: "cb-ssl-xdr1"},
				{ContainerPort: int32(11215), Name: "cb-ssl-xdr2"},
				{ContainerPort: int32(18091), Name: "cb-admin-ssl"},
				{ContainerPort: int32(18092), Name: "cb-views-ssl"},
				{ContainerPort: int32(18093), Name: "cb-queries-ssl"},
				{ContainerPort: int32(4369), Name: "empd"},
			},
			VolumeMounts: mounts,
			Resources:    res,
		},
		{
			Name:            "couchbase-sidecar",
			Image:           "pearsontechnology/couchbase-sidecar:latest",
			ImagePullPolicy: v1.PullAlways,
			Env: []v1.EnvVar{
				{Name: "SIDECAR_SERVICE", Value: w.BiteService.Name},
				{Name: "SIDECAR_MASTERNODE", Value: "cb-0"},
				{Name: "SIDECAR_HOST", ValueFrom: &v1.EnvVarSource{
					FieldRef: &v1.ObjectFieldSelector{FieldPath: "metadata.name"}},
				},
				{Name: "SIDECAR_NAMESPACE", ValueFrom: &v1.EnvVarSource{
					FieldRef: &v1.ObjectFieldSelector{FieldPath: "metadata.namespace"}},
				},
			},
			VolumeMounts: []v1.VolumeMount{
				{
					MountPath: "/opt/couchbase/var/lib/couchbase/backups",
					Name:      w.BiteService.Name + "-backups",
				},
			},
		},
	}
	return cb, err
}

// CbVolumeClaimTemplates generates a slice of persistent volume claims required in a couchbase statefulset
func (w *KubeMapper) CbVolumeClaimTemplates() ([]v1.PersistentVolumeClaim, error) {
	// add pvcs from service definition
	pvcs, _ := w.PersistentVolumeClaims()
	// add backups pvc
	b := v1.PersistentVolumeClaim{
		ObjectMeta: v1.ObjectMeta{
			Name:      w.BiteService.Name + "-backups",
			Namespace: w.Namespace,
			Annotations: map[string]string{
				"volume.beta.kubernetes.io/storage-class": "aws-ebs",
			},
			Labels: map[string]string{
				"creator":    "pipeline",
				"deployment": w.BiteService.Name,
				"mount_path": "/opt/couchbase/var/lib/couchbase/backups",
				"size":       "100Gi",
			},
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{"ReadWriteMany"},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceName(v1.ResourceStorage): resource.MustParse("100Gi"),
				},
			},
		},
	}
	pvcs = append(pvcs, b)
	return pvcs, nil
}

//
