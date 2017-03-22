package kubernetes

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	api_v1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/rest"

	log "github.com/Sirupsen/logrus"
	"github.com/pearsontechnology/environment-operator/pkg/bitesize"
	"github.com/pearsontechnology/environment-operator/pkg/translator"
	"github.com/pearsontechnology/environment-operator/pkg/util"
	"github.com/pearsontechnology/environment-operator/pkg/util/k8s"
)

// Wrapper wraps low level kubernetes api requests to an object easier
// to interact with
type Wrapper struct {
	kubernetes.Interface
}

// NewWrapper returns default in-cluster kubernetes client
func NewWrapper() (*Wrapper, error) {
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return &Wrapper{Interface: clientset}, nil
}

func listOptions() api_v1.ListOptions {
	return api_v1.ListOptions{
		LabelSelector: "creator=pipeline",
	}
}

// Services returns the list of environment-operator managed services within
// given namespace
func (w *Wrapper) Services(ns string) ([]api_v1.Service, error) {
	list, err := w.Core().Services(ns).List(listOptions())

	if err != nil {
		return nil, err
	}

	return list.Items, nil

}

// Ingresses returns the list of environment-operator managed ingresses within
// given namespace
func (w *Wrapper) Ingresses(ns string) ([]v1beta1.Ingress, error) {

	list, err := w.Extensions().Ingresses(ns).List(listOptions())
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

// Deployments returns the list of environment-operator managed deployments within
// given namespace
func (w *Wrapper) Deployments(ns string) ([]v1beta1.Deployment, error) {

	list, err := w.Extensions().Deployments(ns).List(listOptions())
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

// Deployment returns Deployment in the namespace by given name
func (w *Wrapper) Deployment(ns, name string) (*v1beta1.Deployment, error) {
	return w.Extensions().Deployments(ns).Get(name)
}

// PersistentVolumeClaims returns the list of environment-operator managed
// persistent volume claims  within given namespace
func (w *Wrapper) PersistentVolumeClaims(ns string) ([]api_v1.PersistentVolumeClaim, error) {
	list, err := w.Core().PersistentVolumeClaims(ns).List(listOptions())
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

// PersistentVolumeClaimsForDeployment returns the list of persistent volume claims
// attached to the specific deployment, identified by name
func (w *Wrapper) PersistentVolumeClaimsForDeployment(ns, deployment string) ([]api_v1.PersistentVolumeClaim, error) {
	labels := fmt.Sprintf("creator=pipeline,deployment=%s", deployment)
	opts := api_v1.ListOptions{
		LabelSelector: labels,
	}
	list, err := w.Core().PersistentVolumeClaims(ns).List(opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

// VolumesForDeployment returns list of bitesize formatted volumes for a specific
// deployment, identified by name
func (w *Wrapper) VolumesForDeployment(namespace string, name string) ([]bitesize.Volume, error) {
	var retval []bitesize.Volume
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

		q := claim.Spec.Resources.Requests["storage"]

		vol := bitesize.Volume{
			Size:  q.String(),
			Path:  vmount.MountPath,
			Modes: getAccessModesAsString(claim.Spec.AccessModes),
		}
		retval = append(retval, vol)
	}
	return retval, nil
}

// EnvVarsForDeployment returns of bitesize formatted env vars from deployment
// name. XXX: figure out how to returns secrets in this list as well.
func (w *Wrapper) EnvVarsForDeployment(namespace, name string) ([]bitesize.EnvVar, error) {
	var retval []bitesize.EnvVar
	d, err := w.Deployment(namespace, name)
	if err != nil {
		return nil, err
	}

	for _, e := range d.Spec.Template.Spec.Containers[0].Env {
		v := bitesize.EnvVar{
			Name:  e.Name,
			Value: e.Value,
		}
		retval = append(retval, v)
	}
	return retval, nil
}

// HealthCheckForDeployment returns bitesize formatted bitesizeliveness object
// given deployment name
func (w *Wrapper) HealthCheckForDeployment(namespace, name string) (*bitesize.HealthCheck, error) {
	var retval *bitesize.HealthCheck

	d, err := w.deploymentFromName(namespace, name)
	if err != nil {
		return nil, err
	}
	probe := d.Spec.Template.Spec.Containers[0].LivenessProbe
	if probe != nil && probe.Exec != nil {

		retval = &bitesize.HealthCheck{
			Command:      probe.Exec.Command,
			InitialDelay: int(probe.InitialDelaySeconds),
			Timeout:      int(probe.TimeoutSeconds),
		}
	}

	return retval, nil
}

// NamespaceInfo loads namespace information from Kubernetes given namespace name.
// Used to retrieve essential namespace labels (e.g. project name, environment name)
func (w *Wrapper) NamespaceInfo(ns string) (*api_v1.Namespace, error) {
	return w.Core().Namespaces().Get(ns)
}

func (w *Wrapper) deploymentFromName(namespace, name string) (v1beta1.Deployment, error) {
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

func (w *Wrapper) volumeMountFromName(d v1beta1.Deployment, name string) (api_v1.VolumeMount, error) {
	if len(d.Spec.Template.Spec.Containers) == 0 {
		return api_v1.VolumeMount{}, fmt.Errorf("No volume mount %s found", name)
	}

	for _, vmount := range d.Spec.Template.Spec.Containers[0].VolumeMounts {
		if vmount.Name == name {
			return vmount, nil
		}
	}
	return api_v1.VolumeMount{}, fmt.Errorf("No volume mount %s found", name)
}

// ApplyEnvironment executes kubectl apply against ingresses, services, deployments
// etc.
func (w *Wrapper) ApplyEnvironment(e *bitesize.Environment) {
	var err error

	for _, service := range e.Services {

		mapper := &translator.KubeMapper{
			BiteService: &service,
			Namespace:   e.Namespace,
		}

		if service.Type == "" {
			if err = w.updateService(mapper); err != nil {
				log.Error(err)
			}

			if err = w.updateDeployment(mapper); err != nil {
				log.Error(err)
			}

			if err = w.updatePersistentVolumeClaims(mapper); err != nil {
				log.Error(err)
			}

			if service.ExternalURL != "" {

				if err = w.updateIngress(mapper); err != nil {
					log.Error(err)
				}
			}

		} else {
			if err = w.updateThirdPartyResource(mapper); err != nil {
				log.Error(err)
			}
		}
		// 	var tprconfig *rest.Config
		// tprconfig = config
		// configureClient(tprconfig)

		// tprclient, err := rest.RESTClientFor(tprconfig)
		// if err != nil {
		// 	panic(err)
		// }
		// err = tprclient.Post().
		// 	Resource("examples").
		// 	Namespace(api.NamespaceDefault).
		// 	Body(example).
		// 	Do().Into(&result)

	}

}

func (w *Wrapper) updateService(m *translator.KubeMapper) error {
	var err error
	var existingSvc *api_v1.Service

	svc, err := m.Service()
	if err != nil {
		return err
	}

	if existingSvc, err = w.Core().Services(m.Namespace).Get(svc.Name); err == nil {
		log.Debugf("Updating service %s", svc.Name)
		version := existingSvc.GetResourceVersion()
		svc.ResourceVersion = version
		svc.Spec.ClusterIP = existingSvc.Spec.ClusterIP
		_, err = w.Core().Services(m.Namespace).Update(svc)
	} else {
		log.Debugf("Creating service %s", svc.Name)
		_, err = w.Core().Services(m.Namespace).Create(svc)
	}
	return err
}

func (w *Wrapper) updateDeployment(m *translator.KubeMapper) error {
	var err error

	deployment, err := m.Deployment()
	if err != nil {
		return err
	}

	client := k8s.Deployment{
		Namespace: m.Namespace,
		Interface: w.Interface,
	}

	if client.Exist(deployment.Name) {
		log.Debugf("Updating deployment %s", deployment.Name)
		err = client.Update(deployment)
	} else if m.BiteService.Version != "" {
		log.Debugf("Creating deployment %s", deployment.Name)
		image, _ := util.ApplicationImage(m.BiteService)

		deployment.Spec.Template.Spec.Containers[0].Image = image
		err = client.Create(deployment)
	}

	return err
}

func (w *Wrapper) updateThirdPartyResource(m *translator.KubeMapper) error {
	var err error
	// var result ext.PrsnExternalResource
	// var rsc ext.PrsnExternalResource

	tpr, err := m.ThirdPartyResource()
	if err != nil {
		log.Error(err)
	}

	iface, err := NewTPRClient()
	if err != nil {
		log.Error(err)
		return err
	}

	client := k8s.ThirdPartyResource{
		Interface: iface,
		Namespace: m.Namespace,
		Type:      tpr.Kind,
	}

	return client.Apply(tpr)
}

func (w *Wrapper) updateIngress(m *translator.KubeMapper) error {
	var err error

	ingress, err := m.Ingress()
	if err != nil {
		return err
	}

	client := k8s.Ingress{
		Namespace: m.Namespace,
		Interface: w.Interface,
	}

	return client.Apply(ingress)
}

func (w *Wrapper) updatePersistentVolumeClaims(m *translator.KubeMapper) error {
	var err error

	claims, err := m.PersistentVolumeClaims()
	if err != nil {
		return err
	}

	client := k8s.PersistentVolumeClaim{
		Interface: w.Interface,
		Namespace: m.Namespace,
	}

	for _, claim := range claims {

		if err = client.Apply(&claim); err != nil {
			return err
		}
	}

	return nil
}
