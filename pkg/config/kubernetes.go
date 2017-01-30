package config

import (
	"strings"

	"k8s.io/kubernetes/pkg/api/v1"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
)

// LoadFromClient returns BitesizeEnvironment object loaded from Kubernetes API
func LoadFromClient(client clientset.Interface) (*EnvironmentsBitesize, error) {
	// var err error

	namespace := "sample-app-dev"
	wrapper := &KubernetesWrapper{client}

	serviceMap := make(map[string]*BitesizeService)

	services, _ := wrapper.Services(namespace)
	for _, kubeService := range services {
		// name := trimBlueGreenFromName(kubeService.Name)
		name := kubeService.Name
		serviceMap[name] = &BitesizeService{
			Name: name,
			Port: int(kubeService.Spec.Ports[0].Port),
		}
	}

	deployments, _ := wrapper.Deployments(namespace)
	for _, kubeDeployment := range deployments {
		// name := trimBlueGreenFromName(kubeDeployment.Name)
		name := kubeDeployment.Name

		if serviceMap[name] == nil {
			serviceMap[name] = &BitesizeService{}
		}

		// volumeClaims := client.Core().PersistentVolumeClaims(kubeDeployment.Namespace).List()

		serviceMap[name].Replicas = int(kubeDeployment.Spec.Replicas)
		serviceMap[name].Ssl = kubeDeployment.Labels["ssl"]
		serviceMap[name].Version = kubeDeployment.Labels["version"]
		serviceMap[name].Application = kubeDeployment.Labels["application"]
		serviceMap[name].HTTPSOnly = kubeDeployment.Labels["httpsOnly"]
		serviceMap[name].HTTPSBackend = kubeDeployment.Labels["httpsBackend"]
		serviceMap[name].Volumes, _ = wrapper.VolumesForDeployment(namespace, kubeDeployment.Name)
		serviceMap[name].EnvVars, _ = wrapper.EnvVarsForDeployment(namespace, kubeDeployment.Name)
		serviceMap[name].HealthCheck, _ = wrapper.HealthCheckForDeployment(namespace, kubeDeployment.Name)
	}

	ingresses, _ := wrapper.Ingresses(namespace)
	for _, kubeIngress := range ingresses {
		// name := trimBlueGreenFromName(kubeIngress.Name)
		name := kubeIngress.Name
		if serviceMap[name] == nil {
			serviceMap[name] = &BitesizeService{}
		}

		// serviceMap[name].ExternalURL = trimBlueGreenFromHost(kubeIngress.Spec.Rules[0].Host)
		serviceMap[name].ExternalURL = kubeIngress.Spec.Rules[0].Host

		// fmt.Printf("%+v", kubeIngress)
	}

	// spew.Dump(serviceMap)
	return nil, nil
}

func trimBlueGreenFromName(orig string) string {
	return strings.TrimSuffix(strings.TrimSuffix(orig, "-blue"), "-green")
}

func trimBlueGreenFromHost(orig string) string {
	split := strings.Split(orig, ".")
	split[0] = trimBlueGreenFromName(split[0])
	return strings.Join(split, ".")
}
func collectHealthCheck(probe *v1.Probe) *BitesizeLiveness {
	return &BitesizeLiveness{}
}
