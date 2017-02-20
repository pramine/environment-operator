package git

import git2go "gopkg.in/libgit2/git2go.v25"

// Clone clonse remote git repo remotePath to localPath
func (g *Git) Clone() error {

	_, err := git2go.Clone(g.RemotePath, g.LocalPath, cloneOptions())
	return err
}
