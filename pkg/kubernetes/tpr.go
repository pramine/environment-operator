package kubernetes

import (
	"fmt"
	"net/http"

	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/rest/fake"
)

// NewTPRClient returns rest client for prsn.io handling
func NewTPRClient() (*rest.RESTClient, error) {
	var restcli *rest.RESTClient
	var err error

	// const GroupName = "prsn.io"
	// var SchemeGroupVersion = unversioned.GroupVersion{Group: GroupName, Version: "v1"}
	//
	// var SchemeBuilder = runtime.NewSchemeBuilder(
	// 	func(scheme *runtime.Scheme) error {
	// 		scheme.AddKnownTypes(SchemeGroupVersion,
	// 			&api.ListOptions{},
	// 			&api.DeleteOptions{},
	// 			&api.ExportOptions{},
	// 			&PrsnExternalResource{},
	// 			&PrsnExternalResourceList{},
	// 		)
	// 		scheme.AddKnownTypeWithName(SchemeGroupVersion.WithKind(kind))
	// 		return nil
	// 	})
	//
	// SchemeBuilder.AddToScheme(api.Scheme)
	//
	// if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	config.GroupVersion = &unversioned.GroupVersion{
		Group:   "prsn.io",
		Version: "v1",
	}
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: api.Codecs}

	restcli, err = rest.RESTClientFor(config)
	// } else {
	// 	restcli = fakeClient()
	// }

	if err != nil {
		return nil, err
	}
	return restcli, nil
}

func fakeClient() *fake.RESTClient {
	return &fake.RESTClient{
		GroupName:            "prsn.io",
		NegotiatedSerializer: serializer.DirectCodecFactory{CodecFactory: api.Codecs},
		Client: fake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {
			switch p, m := req.URL.Path, req.Method; {
			case p == "/clusters" && m == http.MethodPost:
				return &http.Response{StatusCode: http.StatusCreated}, nil
			default:
				return nil, fmt.Errorf("unexpected request: %#v\n%#v", req.URL, req)
			}
		}),
	}
}
