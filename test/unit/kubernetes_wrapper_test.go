package unit

import (
	"testing"

	"github.com/pearsontechnology/environment-operator/pkg/config"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/resource"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/fake"
)

func TestKubernetesWrapper(t *testing.T) {
	t.Run("test service count", testServiceCount)
	t.Run("test volumes", testVolumes)
}

func testServiceCount(t *testing.T) {

	client := fake.NewSimpleClientset(
		newService("test", "test1", "test", 80),
		newService("test", "test2", "test", 81),
	)

	wrapper := config.KubernetesWrapper{Interface: client}

	services, err := wrapper.Services("test")
	if err != nil {
		t.Errorf("Unexpected exception: %s", err.Error())
	}

	if len(services) != 2 {
		t.Errorf("Unexpected count of services")
	}
}

func testVolumes(t *testing.T) {
	capacity, _ := resource.ParseQuantity("59G")
	validLabels := map[string]string{"creator": "pipeline"}

	client := fake.NewSimpleClientset(
		&api.PersistentVolumeClaim{
			ObjectMeta: api.ObjectMeta{
				Name:      "ts",
				Namespace: "test",
			},
		},
		&api.PersistentVolumeClaim{
			ObjectMeta: validMeta("test", "test"),
			Spec: api.PersistentVolumeClaimSpec{
				AccessModes: []api.PersistentVolumeAccessMode{
					api.ReadWriteOnce,
				},
				Resources: api.ResourceRequirements{
					Requests: api.ResourceList{
						api.ResourceStorage: capacity,
					},
				},
			},
		},
		&api.PersistentVolume{
			ObjectMeta: api.ObjectMeta{
				Name:   "test",
				Labels: validLabels,
			},
		},
		&extensions.Deployment{
			ObjectMeta: validMeta("test", "test"),
			Spec: extensions.DeploymentSpec{
				Replicas: 1,
				Template: api.PodTemplateSpec{
					ObjectMeta: validMeta("test", "test"),
					Spec: api.PodSpec{
						Containers: []api.Container{
							{
								VolumeMounts: []api.VolumeMount{
									{
										Name:      "test",
										MountPath: "/tmp/blah",
										ReadOnly:  true,
									},
								},
							},
						},
						Volumes: []api.Volume{
							{
								Name: "test",
								VolumeSource: api.VolumeSource{
									PersistentVolumeClaim: &api.PersistentVolumeClaimVolumeSource{
										ClaimName: "test",
									},
								},
							},
						},
					},
				},
			},
		},
		// &api.Service{
		// 	ObjectMeta: api.ObjectMeta{
		// 		Name:      "ts",
		// 		Namespace: "test",
		// 	},
		// },
		// &extensions.Ingress{
		// 	ObjectMeta: api.ObjectMeta{
		// 		Name:      "ts",
		// 		Namespace: "test",
		// 	},
		// },
	)
	wrapper := config.KubernetesWrapper{Interface: client}

	a, err := wrapper.PersistentVolumeClaims("test")
	if err != nil {
		t.Errorf("Unexpected exception: %s", err.Error())
	}

	if len(a) != 1 {
		t.Errorf("Unexpected amount of persistent volume claims returned: %d, expected 1", len(a))
	}

	r, err := wrapper.VolumesForDeployment("test", "test")
	if err != nil {
		t.Errorf("Error retrieving volumes: %s", err.Error())
	}

	if len(r) != 1 {
		t.Errorf("Unexpected amount of volumes in VolumesForDeployment: %d", len(r))
	}

	if r[0].Size != "59G" {
		t.Errorf("Error in volume size, expected: %s, actual: %s", "59G", r[0].Size)
	}

}

func newDeployment(namespace, name string) *extensions.Deployment {
	d := extensions.Deployment{
		ObjectMeta: api.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: extensions.DeploymentSpec{
			Template: api.PodTemplateSpec{},
		},
	}
	return &d
}

func newService(namespace, serviceName, portName string, portNumber int32) *api.Service {
	labels := map[string]string{"creator": "pipeline"}

	service := api.Service{
		ObjectMeta: api.ObjectMeta{
			Name:      serviceName,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: api.ServiceSpec{
			ClusterIP: "1.1.1.1",
			Ports: []api.ServicePort{
				{Port: portNumber, Name: portName, Protocol: "TCP"},
			},
		},
	}
	return &service
}

func validMeta(name, namespace string) api.ObjectMeta {
	validLabels := map[string]string{"creator": "pipeline"}

	if namespace != "" {
		return api.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    validLabels,
		}
	}
	return api.ObjectMeta{
		Name:   name,
		Labels: validLabels,
	}
}
