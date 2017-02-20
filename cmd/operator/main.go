package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/pearsontechnology/environment-operator/pkg/git"
	"github.com/pearsontechnology/environment-operator/version"
)

// Config contains environment variables used to configure the app
type Config struct {
	GitRepo   string `envconfig:"GIT_REMOTE_REPOSITORY"`
	GitBranch string `envconfig:"GIT_BRANCH" default:"master"`
	GitKey    string `envconfig:"GIT_PRIVATE_KEY"`
}

func main() {
	var cfg Config
	err := envconfig.Process("operator", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	g := &git.Git{
		Src:        "/tmp/asd",
		Dest:       cfg.GitRepo,
		BranchName: cfg.GitBranch,
		SSHKey:     cfg.GitKey,
	}

	if ok, err := g.UpdatesExist(); ok {
		log.Error(err.Error())
		log.Info("must update")
	}

	fmt.Println(version.Version)
}
