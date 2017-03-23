package kubernetes

import (
	"fmt"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"

	"github.com/pearsontechnology/environment-operator/pkg/bitesize"
	"github.com/pearsontechnology/environment-operator/pkg/diff"
)

// ApplyIfChanged compares bitesize Environment passed as an argument to
// the current client environment. If there are any changes, c is applied
// to the current config
func (wrapper *Wrapper) ApplyIfChanged(newConfig *bitesize.Environment) error {
	log.Infof("Loading namespace: %s", newConfig.Namespace)
	currentConfig, _ := wrapper.LoadEnvironment(newConfig.Namespace)

	changes := diff.Compare(*newConfig, *currentConfig)
	if changes != "" {
		log.Infof("Changes: %s", changes)
		wrapper.ApplyEnvironment(newConfig)
	}
	return nil
}

// LoadEnvironment returns BitesizeEnvironment object loaded from Kubernetes API
func (wrapper *Wrapper) LoadEnvironment(namespace string) (*bitesize.Environment, error) {
	// var err error

	serviceMap := make(map[string]*bitesize.Service)

	ns, err := wrapper.NamespaceInfo(namespace)
	if err != nil {
		return nil, fmt.Errorf("Namespace %s not found", namespace)
	}
	environmentName := ns.ObjectMeta.Labels["environment"]

	services, err := wrapper.Services(namespace)
	if err != nil {
		log.Errorf("Error loading kubernetes services: %s", err.Error())
	}
	for _, kubeService := range services {
		// name := trimBlueGreenFromName(kubeService.Name)

		name := kubeService.Name
		kubePort := 0
		if len(kubeService.Spec.Ports) > 0 {
			kubePort = int(kubeService.Spec.Ports[0].Port)
		}
		// Replicas: represent default value. Overriden if deployments found
		serviceMap[name] = &bitesize.Service{
			Name:     name,
			Port:     kubePort,
			Replicas: 1,
		}
	}

	deployments, err := wrapper.Deployments(namespace)
	if err != nil {
		log.Errorf("Error loading kubernetes deployments: %s", err.Error())
	}
	for _, kubeDeployment := range deployments {
		// name := trimBlueGreenFromName(kubeDeployment.Name)
		name := kubeDeployment.Name

		if serviceMap[name] == nil {
			serviceMap[name] = &bitesize.Service{
				Name: name,
			}
		}
		if kubeDeployment.Spec.Replicas != nil {
			serviceMap[name].Replicas = int(*kubeDeployment.Spec.Replicas)
		} else {
			serviceMap[name].Replicas = 1
		}
		serviceMap[name].Ssl = getLabel(kubeDeployment, "ssl") // kubeDeployment.Labels["ssl"]
		serviceMap[name].Version = getLabel(kubeDeployment, "version")
		serviceMap[name].Application = getLabel(kubeDeployment, "application")
		serviceMap[name].HTTPSOnly = getLabel(kubeDeployment, "httpsOnly")
		serviceMap[name].HTTPSBackend = getLabel(kubeDeployment, "httpsBackend")
		serviceMap[name].Volumes, _ = wrapper.VolumesForDeployment(namespace, kubeDeployment.Name)
		serviceMap[name].EnvVars, _ = wrapper.EnvVarsForDeployment(namespace, kubeDeployment.Name)
		serviceMap[name].HealthCheck, _ = wrapper.HealthCheckForDeployment(namespace, kubeDeployment.Name)
	}

	ingresses, err := wrapper.Ingresses(namespace)
	if err != nil {
		log.Errorf("Error loading kubernetes ingresses: %s", err.Error())
	}

	for _, kubeIngress := range ingresses {
		// name := trimBlueGreenFromName(kubeIngress.Name)
		name := kubeIngress.Name
		if serviceMap[name] == nil {
			serviceMap[name] = &bitesize.Service{
				Name: name,
			}
		}
		var externalURL string

		if len(kubeIngress.Spec.Rules) > 0 {
			externalURL = kubeIngress.Spec.Rules[0].Host
		}

		serviceMap[name].ExternalURL = externalURL
	}

	claims, _ := wrapper.PersistentVolumeClaims(namespace)
	for _, claim := range claims {
		sName := claim.ObjectMeta.Labels["deployment"]

		if sName != "" {
			vol := bitesize.Volume{
				Path:  claim.ObjectMeta.Labels["mount_path"],
				Modes: getAccessModesAsString(claim.Spec.AccessModes),
				Size:  claim.ObjectMeta.Labels["size"],
				Name:  claim.ObjectMeta.Name,
			}

			serviceMap[sName].Volumes = append(serviceMap[sName].Volumes, vol)

		}

	}

	var serviceList bitesize.Services

	for _, v := range serviceMap {
		serviceList = append(serviceList, *v)
	}

	sort.Sort(serviceList)

	bitesizeConfig := bitesize.Environment{
		Name:      environmentName,
		Namespace: namespace,
		Services:  serviceList,
	}

	// spew.Dump(serviceMap)
	return &bitesizeConfig, nil
}

func getLabel(resource v1beta1.Deployment, label string) string {
	if (len(resource.ObjectMeta.Labels) > 0) &&
		(resource.ObjectMeta.Labels[label] != "") {
		return resource.ObjectMeta.Labels[label]
	}
	return ""
}

func getAccessModesAsString(modes []v1.PersistentVolumeAccessMode) string {

	modesStr := []string{}
	if containsAccessMode(modes, v1.ReadWriteOnce) {
		modesStr = append(modesStr, "ReadWriteOnce")
	}
	if containsAccessMode(modes, v1.ReadOnlyMany) {
		modesStr = append(modesStr, "ReadOnlyMany")
	}
	if containsAccessMode(modes, v1.ReadWriteMany) {
		modesStr = append(modesStr, "ReadWriteMany")
	}
	return strings.Join(modesStr, ",")
}

func containsAccessMode(modes []v1.PersistentVolumeAccessMode, mode v1.PersistentVolumeAccessMode) bool {
	for _, m := range modes {
		if m == mode {
			return true
		}
	}
	return false
}
