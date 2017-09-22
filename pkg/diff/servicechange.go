package diff

import (
	log "github.com/Sirupsen/logrus"
)

var changeMap map[string]string

func newChangeMap() {
	changeMap = make(map[string]string)
}

func addServiceChange(svc, diff string) {
	changeMap[svc] = diff
}

func ServiceChanged(serviceName string) bool {
	val, serviceChangeExists := changeMap[serviceName]
	if serviceChangeExists {
		log.Debugf("Applying changes to the '%s' service:\n %s", serviceName, val)
		return true
	}
	return false
}

func Changes() map[string]string {
	return changeMap
}
