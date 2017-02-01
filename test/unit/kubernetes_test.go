package unit

import (
	"testing"

	"github.com/pearsontechnology/environment-operator/pkg/config"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/resource"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/fake"
)

func testFullBitesizeEnvironment(t *testing.T) {
	capacity, _ := resource.ParseQuantity("59G")
	validLabels := map[string]string{"creator": "pipeline"}
	nsLabels := map[string]string{"environment": "Development"}

	client := fake.NewSimpleClientset(
		&api.Namespace{
			ObjectMeta: api.ObjectMeta{
				Name:   "test",
				Labels: nsLabels,
			},
		},
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
								Env: []api.EnvVar{
									{
										Name:  "test",
										Value: "1",
									},
									{
										Name:  "test2",
										Value: "2",
									},
									{
										Name:  "test3",
										Value: "3",
									},
								},

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
		&api.Service{
			ObjectMeta: api.ObjectMeta{
				Name:      "ts",
				Namespace: "test",
			},
		},
		&api.Service{
			ObjectMeta: validMeta("test", "test"),
			Spec: api.ServiceSpec{
				Ports: []api.ServicePort{
					{
						Name:     "whatevs",
						Protocol: "TCP",
						Port:     80,
					},
				},
			},
		},
		&api.Service{
			ObjectMeta: validMeta("test", "test2"),
		},
		&extensions.Ingress{
			ObjectMeta: api.ObjectMeta{
				Name:      "ts",
				Namespace: "test",
			},
		},
		&extensions.Ingress{
			ObjectMeta: validMeta("test", "test"),
			Spec: extensions.IngressSpec{
				Rules: []extensions.IngressRule{
					{
						Host: "www.test.com",
					},
				},
			},
		},
	)
	bitesize, err := config.LoadFromClient(client, "test")
	if err != nil {
		t.Error(err)
	}
	if bitesize == nil {
		t.Error("Bitesize object is nil")
	}

	if len(bitesize.Environments) != 1 {
		t.Errorf("Unexpected environment count: %d, expected: 1", len(bitesize.Environments))
	}

	environment := bitesize.Environments[0]

	if environment.Name != "Development" {
		t.Errorf("Unexpected environment name: %s", environment.Name)
	}

	if len(environment.Services) != 2 {
		t.Errorf("Unexpected service count: %d, expected: 2", len(environment.Services))
	}

	svc := environment.Services[0]
	if svc.Name != "test" {
		t.Errorf("Unexpected service name: %s, expected: test", svc.Name)
	}

	// TODO: test ingresses, env variables, replica count

	if svc.ExternalURL != "www.test.com" {
		t.Errorf("Unexpected external URL: %s, expected: www.test.com", svc.ExternalURL)
	}

	if len(svc.EnvVars) != 3 {
		t.Errorf("Unexpected environment variable count: %d, expected: 3", len(svc.EnvVars))
	}
}
