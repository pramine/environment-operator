package diff

import (
	"testing"

	"github.com/pearsontechnology/environment-operator/pkg/bitesize"
)

func TestDiffEmpty(t *testing.T) {
	a := bitesize.Environment{}
	b := bitesize.Environment{}

	if Compare(a, b); len(Changes()) != 0 {
		t.Errorf("Expected diff to be empty, got: %s", Changes())
	}

}

func TestIgnoreTestFields(t *testing.T) {
	a := bitesize.Environment{Name: "E", Tests: []bitesize.Test{}}
	b := bitesize.Environment{Name: "A", Tests: []bitesize.Test{
		{Name: "a"},
	}}

	if Compare(a, b); len(Changes()) != 0 {
		t.Errorf("Expected diff to be empty, got: %s", Changes())
	}
}

func TestIgnoreDeploymentFields(t *testing.T) {
	a := bitesize.Environment{Deployment: &bitesize.DeploymentSettings{}}
	b := bitesize.Environment{Deployment: &bitesize.DeploymentSettings{
		Method: "bluegreen",
	}}

	if Compare(a, b); len(Changes()) != 0 {
		t.Errorf("Expected diff to be empty, got: %s", Changes())
	}
}

func TestIgnoreStatusFields(t *testing.T) {
	a := bitesize.Environment{
		Services: bitesize.Services{
			{
				Name: "a",
				Status: bitesize.ServiceStatus{
					AvailableReplicas: 3,
				},
			},
		},
	}

	b := bitesize.Environment{
		Services: bitesize.Services{
			{
				Name: "a",
				Status: bitesize.ServiceStatus{
					AvailableReplicas: 1,
				},
			},
		},
	}

	if Compare(a, b); len(Changes()) != 0 {
		t.Errorf("Expected diff to be empty, got: %s", Changes())
	}
}

func TestDiffNames(t *testing.T) {
	a := bitesize.Environment{Name: "asd"}
	b := bitesize.Environment{Name: "asdf"}

	if Compare(a, b); len(Changes()) != 0 {
		t.Error("Expected diff, got the same")
	}
}

func TestDiffVersionsSame(t *testing.T) {
	var saTests = []struct {
		versionA string
		versionB string
		expected int
	}{
		{"1", "1", 0},
		{"1", "2", 1},
		{"", "1", 0}, // assume the same if environments.bitesize does not have version
		{"1", "", 1}, // assume diff  if cluster is not deployed
	}

	for _, tst := range saTests {
		a := bitesize.Environment{
			Name: "a", Services: []bitesize.Service{{Name: "a", Version: tst.versionA}},
		}
		b := bitesize.Environment{
			Name: "a", Services: []bitesize.Service{{Name: "a", Version: tst.versionB}},
		}

		if Compare(a, b); len(Changes()) != tst.expected {
			t.Errorf(
				"Unexpected version compare(%s,%s) should be %t\n%s\n A %+v\n B %+v",
				tst.versionA, tst.versionB, tst.expected, Changes(), a.Services, b.Services,
			)
		}
	}
}
