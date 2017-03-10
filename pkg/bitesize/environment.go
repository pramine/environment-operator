package bitesize

import (
	"fmt"
	"sort"

	validator "gopkg.in/validator.v2"
)

// Environment represents full managed environment,
// including services, HTTP endpoints and deployments. It can
// be either built from environments.bitesize configuration file
// or Kubernetes cluster
type Environment struct {
	Name       string              `yaml:"name" validate:"nonzero"`
	Namespace  string              `yaml:"namespace,omitempty" validate:"regexp=^[a-zA-Z\\-]*$"` // This field should be optional now
	Deployment *DeploymentSettings `yaml:"deployment,omitempty"`
	Services   Services            `yaml:"services"`
	Tests      []Test              `yaml:"tests,omitempty"`
	// XXX        map[string]interface{} `yaml:",inline"`
}

// Environments is a custom type to implement sort.Interface
type Environments []Environment

func (slice Environments) Len() int {
	return len(slice)
}

func (slice Environments) Less(i, j int) bool {
	return slice[i].Name < slice[j].Name
}

func (slice Environments) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for BitesizeEnvironment.
func (e *Environment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error
	ee := &Environment{}
	type plain Environment
	if err = unmarshal((*plain)(ee)); err != nil {
		return fmt.Errorf("environment.%s", err.Error())
	}

	validator.SetValidationFunc("volume_modes", validVolumeModes)

	*e = *ee

	if err = validator.Validate(e); err != nil {
		return fmt.Errorf("environment.%s", err.Error())
	}
	sort.Sort(e.Services)
	return nil
}
