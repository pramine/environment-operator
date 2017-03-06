package main

import (
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
	"github.com/kylelemons/godebug/pretty"
	"github.com/pearsontechnology/environment-operator/pkg/config"
	"github.com/pearsontechnology/environment-operator/pkg/git"
	"github.com/pearsontechnology/environment-operator/version"
)

// Config contains environment variables used to configure the app
type Config struct {
	GitRepo   string `envconfig:"GIT_REMOTE_REPOSITORY"`
	GitBranch string `envconfig:"GIT_BRANCH" default:"master"`
	GitKey    string `envconfig:"GIT_PRIVATE_KEY"`
	EnvName   string `envconfig:"ENVIRONMENT_NAME"`
	EnvFile   string `envconfig:"BITESIZE_FILE"`
	Namespace string `envconfig:"NAMESPACE"`
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

	client, err := config.NewKubernetesWrapper()

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

func compareConfig(cfg Config, g *git.Git, client *config.KubernetesWrapper) {
	fp := filepath.Join(g.LocalPath, cfg.EnvFile)
	gitEnv, _ := config.Environment(fp, cfg.EnvName)
	kubeEnv, _ := config.LoadFromClient(client, cfg.Namespace)

	compareConfig := &pretty.Config{
		Diffable:       true,
		SkipZeroFields: true,
	}
	diff := compareConfig.Compare(kubeEnv, gitEnv)
	if diff != "" {
		log.Infof(diff)
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
