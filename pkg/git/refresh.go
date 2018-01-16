package git

import log "github.com/Sirupsen/logrus"

// Refresh checks if local git repository copy is outdated. If it is,
// changes are pulled in.
func (g *Git) Refresh() error {
	ok, err := g.UpdatesExist()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if ok {
		log.Infof("Updates in repository: %s", g.RemotePath)
		g.Pull()
	}
	return nil
}
