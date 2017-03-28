package main

import (
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
	"github.com/pearsontechnology/environment-operator/pkg/bitesize"
	"github.com/pearsontechnology/environment-operator/pkg/cluster"
	"github.com/pearsontechnology/environment-operator/pkg/git"
	"github.com/pearsontechnology/environment-operator/pkg/reaper"
	"github.com/pearsontechnology/environment-operator/version"
)

// Config contains environment variables used to configure the app
type Config struct {
	LogLevel       string `envconfig:"LOG_LEVEL"`
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

func main() {
	var cfg Config
	err := envconfig.Process("operator", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	if cfg.LogLevel == "debug" {
		log.SetLevel(log.DebugLevel)
	}

	g := &git.Git{
		LocalPath:  cfg.GitLocalPath,
		RemotePath: cfg.GitRepo,
		BranchName: cfg.GitBranch,
		SSHKey:     cfg.GitKey,
	}
	g.Clone()

	client, err := cluster.NewClusterClient()
	if err != nil {
		log.Fatalf("Error creating kubernetes client: %s", err.Error())
	}

	r := reaper.Reaper{
		Namespace: cfg.Namespace,
		Wrapper:   client,
	}

	log.Infof("Starting up environment-operator version %s", version.Version)

	for {
		g.Refresh()
		gitConfiguration, _ := cfg.LoadEnvironment()
		client.ApplyIfChanged(gitConfiguration)

		go r.Cleanup(gitConfiguration)
		time.Sleep(30000 * time.Millisecond)
	}

}
