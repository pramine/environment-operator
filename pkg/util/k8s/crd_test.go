package k8s

import (
	"testing"
)

func TestCRDGet(t *testing.T) {
	// crdcli := fake.CRDClient(
	// 	&ext.PrsnExternalResource{
	// 		TypeMeta: runtime.TypeMeta{
	// 			Kind: "Mysql",
	// 		},
	// 		ObjectMeta: v1.ObjectMeta{
	// 			Name:      "test",
	// 			Namespace: "test",
	// 		},
	// 	},
	// )
	//
	// client := CustomResourceDefinition{
	// 	Interface: crdcli,
	// 	Namespace: "test",
	// }
	//
	// _, err := client.Get("test")
	// if err != nil {
	// 	t.Errorf("Got unexpected error: %s", err.Error())
	// }

}

func TestCRDExist(t *testing.T) {
}

func TestCRDApply(t *testing.T) {
}

func TestCRDDestroy(t *testing.T) {
}
