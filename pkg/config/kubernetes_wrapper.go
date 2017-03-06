package config

import (
	"fmt"
	"strings"

	"k8s.io/client-go/kubernetes"
	api_v1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/rest"
)

// KubernetesWrapper wraps low level kubernetes api requests to an object easier
// to interact with
type KubernetesWrapper struct {
	kubernetes.Interface
}

// NewKubernetesWrapper returns default in-cluster kubernetes client
func NewKubernetesWrapper() (*KubernetesWrapper, error) {
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return &KubernetesWrapper{Interface: clientset}, nil
}

func listOptions() api_v1.ListOptions {
	return api_v1.ListOptions{
		LabelSelector: "creator=pipeline",
	}
}

// Services returns the list of environment-operator managed services within
// given namespace
func (w *KubernetesWrapper) Services(ns string) ([]api_v1.Service, error) {
	list, err := w.Core().Services(ns).List(listOptions())
	if err != nil {
		return nil, err
	}

	return list.Items, nil

}

// Ingresses returns the list of environment-operator managed ingresses within
// given namespace
func (w *KubernetesWrapper) Ingresses(ns string) ([]v1beta1.Ingress, error) {

	list, err := w.Extensions().Ingresses(ns).List(listOptions())
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

// Deployments returns the list of environment-operator managed deployments within
// given namespace
func (w *KubernetesWrapper) Deployments(ns string) ([]v1beta1.Deployment, error) {

	list, err := w.Extensions().Deployments(ns).List(listOptions())
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

// PersistentVolumeClaims returns the list of environment-operator managed
// persistent volume claims  within given namespace
func (w *KubernetesWrapper) PersistentVolumeClaims(ns string) ([]api_v1.PersistentVolumeClaim, error) {
	list, err := w.Core().PersistentVolumeClaims(ns).List(listOptions())
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

// LoadNamespace fills in KubernetesWrapper with services, ingresses, deployments
// and volumes
// func (w *KubernetesWrapper) LoadNamespace() error {
// 	tlabels := map[string]string{"creator": "pipeline"}
// 	selector := labels.SelectorFromSet(labels.Set(tlabels))
//
// 	opts := api.ListOptions{
// 		LabelSelector: selector,
// 	}
//
// 	ingressList, err := w.Client.Extensions().Ingresses(w.Namespace).List(opts)
// 	if err != nil {
// 		return err
// 	}
//
// 	w.Ingresses = ingressList.Items
//
// 	deploymentList, err := w.Client.Extensions().Deployments(w.Namespace).List(opts)
// 	if err != nil {
// 		return err
// 	}
// 	w.Deployments = deploymentList.Items
//
// 	serviceList, err := w.Client.Core().Services(w.Namespace).List(opts)
// 	spew.Dump(serviceList)
// 	if err != nil {
// 		return err
// 	}
// 	w.Services = serviceList.Items
//
// 	volumeList, err := w.Client.Core().PersistentVolumeClaims(w.Namespace).List(opts)
// 	if err != nil {
// 		return err
// 	}
// 	w.Volumes = volumeList.Items
//
// 	return nil
// }

// VolumesForDeployment returns list of bitesize formatted volumes for a specific
// deployment, identified by name
func (w *KubernetesWrapper) VolumesForDeployment(namespace string, name string) ([]BitesizeVolume, error) {
	var retval []BitesizeVolume
	d, err := w.deploymentFromName(namespace, name)
	if err != nil {
		return nil, err
	}
	for _, v := range d.Spec.Template.Spec.Volumes {
		vmount, err := w.volumeMountFromName(d, v.Name)
		if err != nil {
			return nil, err
		}

		claim, err := w.volumeClaimFromName(namespace, v.Name)
		if err != nil {
			return nil, err
		}

		modes := []string{}
		for m := range claim.Spec.AccessModes {
			modes = append(modes, string(m))
		}

		q := claim.Spec.Resources.Requests["storage"]

		vol := BitesizeVolume{
			Size:  q.String(),
			Path:  vmount.MountPath,
			Modes: strings.Join(modes, ","),
		}
		retval = append(retval, vol)
	}
	return retval, nil
}

// EnvVarsForDeployment returns of bitesize formatted env vars from deployment
// name. XXX: figure out how to returns secrets in this list as well.
func (w *KubernetesWrapper) EnvVarsForDeployment(namespace, name string) ([]BitesizeEnvVar, error) {
	var retval []BitesizeEnvVar
	d, err := w.deploymentFromName(namespace, name)
	if err != nil {
		return nil, err
	}

	for _, e := range d.Spec.Template.Spec.Containers[0].Env {
		v := BitesizeEnvVar{
			Name:  e.Name,
			Value: e.Value,
		}
		retval = append(retval, v)
	}
	return retval, nil
}

// HealthCheckForDeployment returns bitesize formatted bitesizeliveness object
// given deployment name
func (w *KubernetesWrapper) HealthCheckForDeployment(namespace, name string) (*BitesizeLiveness, error) {
	var retval *BitesizeLiveness

	d, err := w.deploymentFromName(namespace, name)
	if err != nil {
		return nil, err
	}
	probe := d.Spec.Template.Spec.Containers[0].LivenessProbe
	if probe != nil && probe.Exec != nil {

		retval = &BitesizeLiveness{
			Command:      probe.Exec.Command,
			InitialDelay: int(probe.InitialDelaySeconds),
			Timeout:      int(probe.TimeoutSeconds),
		}
	}

	return retval, nil
}

// NamespaceInfo loads namespace information from Kubernetes given namespace name.
// Used to retrieve essential namespace labels (e.g. project name, environment name)
func (w *KubernetesWrapper) NamespaceInfo(ns string) (*api_v1.Namespace, error) {
	return w.Core().Namespaces().Get(ns)
}

func (w *KubernetesWrapper) deploymentFromName(namespace, name string) (v1beta1.Deployment, error) {
	deployments, err := w.Deployments(namespace)
	if err != nil {
		return v1beta1.Deployment{}, err
	}

	for _, d := range deployments {
		if d.Name == name {
			return d, nil
		}
	}
	return v1beta1.Deployment{}, fmt.Errorf("No deployment %s found", name)
}

func (w *KubernetesWrapper) volumeMountFromName(d v1beta1.Deployment, name string) (api_v1.VolumeMount, error) {
	for _, vmount := range d.Spec.Template.Spec.Containers[0].VolumeMounts {
		if vmount.Name == name {
			return vmount, nil
		}
	}
	return api_v1.VolumeMount{}, fmt.Errorf("No volume mount %s found", name)
}

func (w *KubernetesWrapper) volumeClaimFromName(namespace, name string) (api_v1.PersistentVolumeClaim, error) {
	claims, err := w.PersistentVolumeClaims(namespace)
	if err != nil {
		return api_v1.PersistentVolumeClaim{}, err
	}

	for _, claim := range claims {
		if claim.Name == name {
			return claim, nil
		}
	}
	return api_v1.PersistentVolumeClaim{}, fmt.Errorf("Persistent volume claim %s not found", name)
}
