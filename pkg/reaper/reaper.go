package reaper

import (
	log "github.com/Sirupsen/logrus"
	"github.com/pearsontechnology/environment-operator/pkg/bitesize"
	"github.com/pearsontechnology/environment-operator/pkg/kubernetes"
)

// Cleanup collects all orphan services (not mentioned in cfg) and
// deletes them from the cluster
func Cleanup(cfg *bitesize.Environment, client *kubernetes.Wrapper) {

	current, err := client.LoadEnvironment(cfg.Namespace)
	if err != nil {
		log.Errorf("REAPER Error loading environment: %s", err.Error())
	}

	for _, service := range current.Services {
		if cfg.Services.FindByName(service.Name) == nil {
			client.DeleteService(service)
		}
	}
}
