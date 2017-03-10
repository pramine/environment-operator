package git

import (
	"fmt"
	"strings"

	git2go "gopkg.in/libgit2/git2go.v24"
)

// UpdatesExist returns true if local HEAD is behind remote
func (g *Git) UpdatesExist() (bool, error) {
	repo, err := git2go.OpenRepository(g.LocalPath)
	if err != nil {
		return true, err
	}

	head, err := repo.Head()
	if err != nil {
		return true, err
	}

	srcTag := head.Target()

	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		return true, err
	}

	if err = remote.Fetch([]string{}, g.fetchOptions(), ""); err != nil {
		return true, err
	}

	branch, err := repo.References.Lookup("refs/remotes/origin/" + g.BranchName)
	if err != nil {
		return true, err
	}

	remoteTag := branch.Target()
	sTag := fmt.Sprintf("%s", srcTag)
	rTag := fmt.Sprintf("%s", remoteTag)

	// log.Infof("git remote tag: %s, local tag: %s", rTag, sTag)

	return (strings.Compare(sTag, rTag) != 0), nil

}
