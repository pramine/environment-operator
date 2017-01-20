package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	validator "gopkg.in/validator.v2"
	yaml "gopkg.in/yaml.v2"
)

// EnvironmentsBitesize is a 1:1 mapping to environments.bitesize file
type EnvironmentsBitesize struct {
	Project      string                 `yaml:"project"`
	Environments []BitesizeEnvironment  `yaml:"environments"`
	XXX          map[string]interface{} `yaml:",inline"`
}

// BitesizeEnvironment represents full managed environment,
// including services, HTTP endpoints and deployments. It can
// be either built from environments.bitesize configuration file
// or Kubernetes cluster
type BitesizeEnvironment struct {
	Name       string                 `yaml:"name" validate:"nonzero"`
	Namespace  string                 `yaml:"namespace,omitempty" validate:"regexp=^[a-zA-Z\\-]*$"` // This field should be optional now
	Deployment *DeploymentSettings    `yaml:"deployment,omitempty"`
	Services   []BitesizeService      `yaml:"services"`
	Tests      []BitesizeTest         `yaml:"tests,omitempty"`
	XXX        map[string]interface{} `yaml:",inline"`
}

// DeploymentSettings represent "deployment" block in environments.bitesize
type DeploymentSettings struct {
	Method string                 `yaml:"method,omitempty" validate:"regexp=^(bluegreen|rolling-upgrade)*$"`
	Mode   string                 `yaml:"mode,omitempty" validate:"regexp=^(manual|auto)*$"`
	Active string                 `yaml:"active,omitempty" validate:"regexp=^(blue|green)*$"`
	XXX    map[string]interface{} `yaml:",inline"`
}

// BitesizeService represents a single service and it's configuration,
// running in environment
type BitesizeService struct {
	Name        string                 `yaml:"name" validate:"nonzero"`
	ExternalURL string                 `yaml:"external_url,omitempty" validate:"regexp=^([a-zA-Z\\.\\-]+$)*"`
	Port        int                    `yaml:"port,omitempty" validate:"max=65535"`
	Ssl         string                 `yaml:"ssl,omitempty" validate:"regexp=^(true|false)*$"`
	Replicas    int                    `yaml:"replicas,omitempty"`
	Deployment  *DeploymentSettings    `yaml:"deployment,omitempty"`
	HealthCheck *BitesizeLiveness      `yaml:"health_check,omitempty"`
	XXX         map[string]interface{} `yaml:",inline"`
}

// BitesizeTest is obsolete and not used by environment-operator,
// but it's here for configuration compatability
type BitesizeTest struct {
	Name       string                 `yaml:"name"`
	Repository string                 `yaml:"repository"`
	Branch     string                 `yaml:"branch"`
	Commands   map[string]string      `yaml:"commands"`
	XXX        map[string]interface{} `yaml:",inline"`
}

// BitesizeLiveness maps to LivenessProbe in Kubernetes
type BitesizeLiveness struct {
	Command      []string               `yaml:"command"`
	InitialDelay int                    `yaml:"initial_delay,omitempty"`
	Timeout      int                    `yaml:"timeout,omitempty"`
	XXX          map[string]interface{} `yaml:",inline"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for BitesizeEnvironment.
func (e *BitesizeEnvironment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error
	ee := &BitesizeEnvironment{}
	type plain BitesizeEnvironment
	if err = unmarshal((*plain)(ee)); err != nil {
		return fmt.Errorf("environment.%s", err.Error())
	}

	*e = *ee

	if err = checkOverflow(e.XXX, "environment"); err != nil {
		return err
	}

	if err = validator.Validate(e); err != nil {
		// return err
		return fmt.Errorf("environment.%s", err.Error())
	}
	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for BitesizeService.
func (e *BitesizeService) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error
	ee := &BitesizeService{}
	type plain BitesizeService
	if err = unmarshal((*plain)(ee)); err != nil {
		return fmt.Errorf("service.%s", err.Error())
	}

	*e = *ee
	if err = validator.Validate(e); err != nil {
		return fmt.Errorf("service.%s", err.Error())
	}
	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for BitesizeLiveness.
func (e *BitesizeLiveness) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error
	ee := &BitesizeLiveness{}
	type plain BitesizeLiveness
	if err = unmarshal((*plain)(ee)); err != nil {
		return fmt.Errorf("health_check.%s", err.Error())
	}

	*e = *ee

	if err = checkOverflow(e.XXX, "health_check"); err != nil {
		return err
	}

	if err = validator.Validate(e); err != nil {
		return fmt.Errorf("health_check.%s", err.Error())
	}
	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for DeploymentSettings.
func (e *DeploymentSettings) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error
	ee := &DeploymentSettings{}
	type plain DeploymentSettings
	if err = unmarshal((*plain)(ee)); err != nil {
		return fmt.Errorf("deployment.%s", err.Error())
	}

	*e = *ee
	if err = validator.Validate(e); err != nil {
		return fmt.Errorf("deployment.%s", err.Error())
	}
	return nil
}

// LoadFromString returns BitesizeEnvironment object from yaml string
func LoadFromString(cfg string) (*EnvironmentsBitesize, error) {
	t := &EnvironmentsBitesize{}

	err := yaml.Unmarshal([]byte(cfg), t)
	if err != nil {
		return t, err
	}

	return t, err
}

// LoadFromFile returns BitesizeEnvironment object loaded from file, passed
// as a path argument.
func LoadFromFile(path string) (*EnvironmentsBitesize, error) {
	var err error
	var contents []byte

	contents, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadFromString(string(contents))
}

func checkOverflow(m map[string]interface{}, ctx string) error {
	if len(m) > 0 {
		var keys []string
		for k := range m {
			keys = append(keys, k)
		}
		return fmt.Errorf("%s: unknown fields (%s)", ctx, strings.Join(keys, ", "))
	}
	return nil
}
