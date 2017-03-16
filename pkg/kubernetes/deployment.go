package kubernetes

import "github.com/pearsontechnology/environment-operator/pkg/translator"

type Deployment struct {
	Mapper  *translator.KubeMapper
	Wrapper *Wrapper
}

func (d *Deployment) Exist() bool {
	deployment, _ := d.Mapper.Deployment()
	_, err := d.Wrapper.
		Extensions().
		Deployments(d.Mapper.Namespace).
		Get(deployment.Name)

	return err == nil
}

func (d *Deployment) Apply() error {
	if d.Exist() {
		return d.Update()
	} else {
		return d.Create()
	}
}

func (d *Deployment) Update() error {
	deployment, _ := d.Mapper.Deployment()
	current, err := d.Wrapper.Extensions().Deployments(d.Mapper.Namespace).Get(deployment.Name)
	if err != nil {
		return err
	}
	deployment.ResourceVersion = current.GetResourceVersion()
	deployment.ObjectMeta.Labels["version"] = current.ObjectMeta.Labels["version"]

	deployment.Spec.Template.Spec.Containers[0].Image = current.Spec.Template.Spec.Containers[0].Image
	_, err = d.Wrapper.
		Extensions().
		Deployments(d.Mapper.Namespace).
		Update(deployment)
	return err
}

func (d *Deployment) Create() error {
	deployment, _ := d.Mapper.Deployment()
	_, err := d.Wrapper.
		Extensions().
		Deployments(d.Mapper.Namespace).
		Create(deployment)
	return err
}
