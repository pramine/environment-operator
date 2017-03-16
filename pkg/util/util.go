package util

import (
	"fmt"
	"os"

	"github.com/pearsontechnology/environment-operator/pkg/bitesize"
)

// func trimBlueGreenFromName(orig string) string {
// 	return strings.TrimSuffix(strings.TrimSuffix(orig, "-blue"), "-green")
// }
//
// func trimBlueGreenFromHost(orig string) string {
// 	split := strings.Split(orig, ".")
// 	split[0] = trimBlueGreenFromName(split[0])
// 	return strings.Join(split, ".")
// }
// func collectHealthCheck(probe *v1.Probe) *bitesize.HealthCheck {
// 	return &bitesize.HealthCheck{}
// }

// Registry returns docker registry setting
func Registry() string {
	return os.Getenv("DOCKER_REGISTRY")
}

// Project returns project's name. TODO: Should be loaded from namespace labels..
func Project() string {
	return os.Getenv("PROJECT")
}

// ApplicationImage returns full image name gievn bitesize.Service object
func ApplicationImage(svc *bitesize.Service) (string, error) {
	return fmt.Sprintf("%s/%s/%s:%s",
		Registry(),
		Project(),
		svc.Application,
		svc.Version,
	), nil
}
