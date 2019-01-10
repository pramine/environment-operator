package git

import (
	gogit "gopkg.in/src-d/go-git.v4"
)

// UpdatesExist returns true if local HEAD is behind remote
func (g *Git) UpdatesExist() (bool, error) {
	err := g.Repository.Fetch(g.fetchOptions())

	if err == gogit.NoErrAlreadyUpToDate {
		return false, nil
	}

	return true, err
}
