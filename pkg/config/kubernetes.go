package config

import (
	"strings"

	"k8s.io/kubernetes/pkg/api/v1"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
)

// LoadFromClient returns BitesizeEnvironment object loaded from Kubernetes API
func LoadFromClient(client clientset.Interface, namespace string) (*EnvironmentsBitesize, error) {
	// var err error
	wrapper := &KubernetesWrapper{client}

	serviceMap := make(map[string]*BitesizeService)

	ns, _ := wrapper.NamespaceInfo(namespace)
	environmentName := ns.Labels["environment"]

	services, _ := wrapper.Services(namespace)
	for _, kubeService := range services {
		// name := trimBlueGreenFromName(kubeService.Name)

		name := kubeService.Name
		kubePort := 0
		if len(kubeService.Spec.Ports) > 0 {
			kubePort = int(kubeService.Spec.Ports[0].Port)
		}
		serviceMap[name] = &BitesizeService{
			Name: name,
			Port: kubePort,
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
		var externalURL string

		if len(kubeIngress.Spec.Rules) > 0 {
			externalURL = kubeIngress.Spec.Rules[0].Host
		}

		serviceMap[name].ExternalURL = externalURL
	}

	var serviceList []BitesizeService

	for _, v := range serviceMap {
		serviceList = append(serviceList, *v)
	}

	bitesizeConfig := EnvironmentsBitesize{
		Environments: []BitesizeEnvironment{
			{
				Name:     environmentName,
				Services: serviceList,
			},
		},
	}

	// spew.Dump(serviceMap)
	return &bitesizeConfig, nil
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
