package git

import git2go "gopkg.in/libgit2/git2go.v24"
import log "github.com/Sirupsen/logrus"

// Clone clonse remote git repo remotePath to localPath
func (g *Git) Clone() error {
	log.Debugf("Cloning remote repository %s", g.RemotePath)

	_, err := git2go.Clone(g.RemotePath, g.LocalPath, g.cloneOptions())
	return err
}
