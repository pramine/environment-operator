package git

import (
	"fmt"
	"os"

	"gopkg.in/src-d/go-git.v4/plumbing"

	log "github.com/Sirupsen/logrus"
	"github.com/pearsontechnology/environment-operator/pkg/config"
	gogit "gopkg.in/src-d/go-git.v4"
	gitconfig "gopkg.in/src-d/go-git.v4/config"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

// Git represents repository object and wraps git2go calls
type Git struct {
	SSHKey     string
	LocalPath  string
	RemotePath string
	BranchName string
	Repository *gogit.Repository
}

func Client() *Git {
	var repository *gogit.Repository
	var err error

	if _, err := os.Stat(config.Env.GitLocalPath); os.IsNotExist(err) {
		repository, err = gogit.PlainInit(config.Env.GitLocalPath, false)
		if err != nil {
			log.Errorf("could not init local repository %s: %s", config.Env.GitLocalPath, err.Error())
		}
	} else {
		repository, err = gogit.PlainOpen(config.Env.GitLocalPath)
	}

	_, err = repository.CreateRemote(&gitconfig.RemoteConfig{
		Name: "origin",
		URLs: []string{config.Env.GitRepo},
	})
	if err != nil {
		log.Errorf("could not attach to origin %s: %s", config.Env.GitRepo, err.Error())
	}

	return &Git{
		LocalPath:  config.Env.GitLocalPath,
		RemotePath: config.Env.GitRepo,
		BranchName: config.Env.GitBranch,
		SSHKey:     config.Env.GitKey,
		Repository: repository,
	}
}

func (g *Git) pullOptions() *gogit.PullOptions {
	branch := fmt.Sprintf("refs/heads/%s", g.BranchName)
	return &gogit.PullOptions{
		ReferenceName: plumbing.ReferenceName(branch),
		Auth:          g.sshKeys(),
	}
}

func (g *Git) fetchOptions() *gogit.FetchOptions {
	return &gogit.FetchOptions{
		Auth: g.sshKeys(),
	}
}

func (g *Git) sshKeys() *gitssh.PublicKeys {
	if g.SSHKey == "" {
		return nil
	}
	auth, err := gitssh.NewPublicKeys("", []byte(g.SSHKey), "")
	if err != nil {
		log.Warningf("error on parsing private key: %s", err.Error())
		return nil
	}
	return auth
}
