package unit

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/fake"
)

func TestKubernetesWrapper(t *testing.T) {
	t.Run("w", testWrapper)
}

func testWrapper(t *testing.T) {
	fake := fake.NewSimpleClientset(
		newDeployment("foo", "bar"),
	)

	i, err := fake.Extensions().Deployments("foo").Get("bar")
	spew.Dump(i)
	if err != nil {
		t.Error(err)
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
