package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
	"github.com/pearsontechnology/environment-operator/pkg/git"
	"github.com/pearsontechnology/environment-operator/pkg/kubernetes"
	"github.com/pearsontechnology/environment-operator/pkg/reaper"
	"github.com/pearsontechnology/environment-operator/version"
)

func main() {
	var cfg Config
	err := envconfig.Process("operator", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	g := &git.Git{
		LocalPath:  cfg.GitLocalPath,
		RemotePath: cfg.GitRepo,
		BranchName: cfg.GitBranch,
		SSHKey:     cfg.GitKey,
	}

	client, err := kubernetes.NewWrapper()
	if err != nil {
		log.Fatalf("Error creating kubernetes client: %s", err.Error())
	}

	log.Infof("Starting up environment-operator version %s", version.Version)

	for {
		g.Refresh()
		gitConfiguration, _ := cfg.LoadEnvironment()
		client.ApplyIfChanged(gitConfiguration)
		go reaper.Cleanup(gitConfiguration, client)
		time.Sleep(30000 * time.Millisecond)
	}

}
