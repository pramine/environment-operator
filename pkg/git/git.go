package git

import (
	"os"

	git2go "gopkg.in/libgit2/git2go.v25"
)

// CloneOrPull checks if repo exists in local path. If it does, it
// pulls changes from remotePath, if it doesn't, performs a full git clone
func CloneOrPull(localPath string, remotePath string) error {
	if _, err := os.Stat("./conf/app.ini"); err == nil {
		return Pull(localPath, remotePath)
	} else {
		return Clone(localPath, remotePath)
	}
}

func Clone(localPath, remotePath string) error {
	return nil
}

func Pull(localPath, remotePath string) error {
	repo, err := git2go.OpenRepository(remotePath)
	if err != nil {
		return err
	}

	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		return err
	}

	if err := remote.Fetch([]string{}, nil, ""); err != nil {
		return err
	}
	return nil
}
