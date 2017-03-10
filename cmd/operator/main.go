package main

import (
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
	"github.com/kylelemons/godebug/pretty"
	"github.com/pearsontechnology/environment-operator/pkg/bitesize"
	"github.com/pearsontechnology/environment-operator/pkg/git"
	"github.com/pearsontechnology/environment-operator/pkg/kubernetes"
	"github.com/pearsontechnology/environment-operator/version"
)

// Config contains environment variables used to configure the app
type Config struct {
	GitRepo        string `envconfig:"GIT_REMOTE_REPOSITORY"`
	GitBranch      string `envconfig:"GIT_BRANCH" default:"master"`
	GitKey         string `envconfig:"GIT_PRIVATE_KEY"`
	EnvName        string `envconfig:"ENVIRONMENT_NAME"`
	EnvFile        string `envconfig:"BITESIZE_FILE"`
	Namespace      string `envconfig:"NAMESPACE"`
	DockerRegistry string `envconfig:"DOCKER_REGISTRY" default:"bitesize-registry.default.svc.cluster.local"`
}

func main() {
	var cfg Config
	err := envconfig.Process("operator", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	g := &git.Git{
		LocalPath:  "/tmp/repository",
		RemotePath: cfg.GitRepo,
		BranchName: cfg.GitBranch,
		SSHKey:     cfg.GitKey,
	}

	client, err := kubernetes.NewWrapper()

	if err != nil {
		log.Errorf("Error creating kubernetes client: %s", err.Error())
	}

	log.Infof("Staring up environment-operator version %s", version.Version)

	if err := g.Clone(); err != nil {
		log.Errorf("Error on initial clone: %s", err.Error())
	}

	for {
		updateGitRepo(g)
		compareConfig(cfg, g, client)

		time.Sleep(30000 * time.Millisecond)
	}

}

func compareConfig(cfg Config, g *git.Git, client *kubernetes.Wrapper) {
	fp := filepath.Join(g.LocalPath, cfg.EnvFile)
	gitEnv, _ := bitesize.LoadEnvironment(fp, cfg.EnvName)
	kubeEnv, _ := kubernetes.LoadFromClient(client, cfg.Namespace)

	// Tests are obsolete.
	gitEnv.Tests = []bitesize.Test{}
	kubeEnv.Tests = []bitesize.Test{}
	gitEnv.Deployment = nil
	kubeEnv.Deployment = nil

	// XXX: remove tprs for now
	var newServices bitesize.Services
	for _, s := range gitEnv.Services {
		if s.Type == "" {
			d := kubeEnv.Services.FindByName(s.Name)
			if d != nil {
				s.Version = d.Version
			}
			newServices = append(newServices, s)
		}
	}
	gitEnv.Services = newServices
	// XXX: the end

	compareConfig := &pretty.Config{
		Diffable:       true,
		SkipZeroFields: true,
	}
	diff := compareConfig.Compare(kubeEnv, gitEnv)
	if diff != "" {
		log.Infof(diff)
		client.ApplyEnvironment(gitEnv)
		// Need to apply gitEnv here
	}
}

func updateGitRepo(g *git.Git) {
	if ok, err := g.UpdatesExist(); ok {
		if err != nil {
			log.Error(err.Error())
		}
		log.Infof("Updates in repository: %s", g.RemotePath)
		g.CloneOrPull()
	}
}
