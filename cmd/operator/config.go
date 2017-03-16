package main

import (
	"path/filepath"

	"github.com/pearsontechnology/environment-operator/pkg/bitesize"
)

// Config contains environment variables used to configure the app
type Config struct {
	GitRepo        string `envconfig:"GIT_REMOTE_REPOSITORY"`
	GitBranch      string `envconfig:"GIT_BRANCH" default:"master"`
	GitKey         string `envconfig:"GIT_PRIVATE_KEY"`
	GitLocalPath   string `envconfig:"GIT_LOCAL_PATH" default:"/tmp/repository"`
	EnvName        string `envconfig:"ENVIRONMENT_NAME"`
	EnvFile        string `envconfig:"BITESIZE_FILE"`
	Namespace      string `envconfig:"NAMESPACE"`
	DockerRegistry string `envconfig:"DOCKER_REGISTRY" default:"bitesize-registry.default.svc.cluster.local"`
}

// LoadEnvironment returns bitesize.Environment object
// constructed from environment variables
func (c Config) LoadEnvironment() (*bitesize.Environment, error) {
	fp := filepath.Join(c.GitLocalPath, c.EnvFile)
	return bitesize.LoadEnvironment(fp, c.EnvName)
}
