package diff

import (
	"testing"
)

func TestAddAndRetrieveChange(t *testing.T) {

	newChangeMap()

	addServiceChange("testservice", "diff string")

	if !ServiceChanged("testservice") {
		t.Errorf("Expected the service should have changed, but was: %s", Changes())

	}

}
