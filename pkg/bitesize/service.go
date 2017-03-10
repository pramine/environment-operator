package bitesize

import (
	"fmt"

	validator "gopkg.in/validator.v2"
)

// Service represents a single service and it's configuration,
// running in environment
type Service struct {
	Name         string              `yaml:"name" validate:"nonzero"`
	ExternalURL  string              `yaml:"external_url,omitempty" validate:"regexp=^([a-zA-Z\\.\\-]+$)*"`
	Port         int                 `yaml:"port,omitempty" validate:"max=65535"`
	Ssl          string              `yaml:"ssl,omitempty" validate:"regexp=^(true|false)*$"`
	Version      string              `yaml:"version,omitempty"`
	Application  string              `yaml:"application,omitempty"`
	Replicas     int                 `yaml:"replicas,omitempty"`
	Deployment   *DeploymentSettings `yaml:"deployment,omitempty"`
	HealthCheck  *HealthCheck        `yaml:"health_check,omitempty"`
	EnvVars      []EnvVar            `yaml:"env,omitempty"`
	Volumes      []Volume            `yaml:"volumes,omitempty"`
	Options      map[string]string   `yaml:"options,omitempty"`
	HTTPSOnly    string
	HTTPSBackend string
	Type         string `yaml:"type,omitempty"`
	// XXX          map[string]interface{} `yaml:",inline"`
}

// Services implement sort.Interface
type Services []Service

// UnmarshalYAML implements the yaml.Unmarshaler interface for BitesizeService.
func (e *Service) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error
	ee := &Service{}
	type plain Service
	if err = unmarshal((*plain)(ee)); err != nil {
		return fmt.Errorf("service.%s", err.Error())
	}

	*e = *ee
	if err = validator.Validate(e); err != nil {
		return fmt.Errorf("service.%s", err.Error())
	}
	// Set defaults
	if e.Type == "" && e.Port == 0 {
		e.Port = 80
	}

	if e.Replicas == 0 {
		e.Replicas = 1
	}

	return nil
}

func (slice Services) Len() int {
	return len(slice)
}

func (slice Services) Less(i, j int) bool {
	return slice[i].Name < slice[j].Name
}

func (slice Services) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// FindByName returns service with a name match
func (slice Services) FindByName(name string) *Service {
	for _, s := range slice {
		if s.Name == name {
			return &s
		}
	}
	return nil
}
