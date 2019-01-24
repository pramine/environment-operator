package git

import log "github.com/Sirupsen/logrus"

// Refresh checks if local git repository copy is outdated. If it is,
// changes are pulled in.
func (g *Git) Refresh() error {
	ok, err := g.UpdatesExist()

	//TODO update to return the repo status and stop from comparing if there are no new changes.
	if err != nil {
		log.Errorf("Error while checking for updates: %s", err.Error())
		return err
	}

	if ok {
		log.Infof("Updates in repository: %s", g.RemotePath)
		if err := g.Pull(); err != nil {
			log.Errorf("Error while pulling the changes from repository: %s", err)
			return err
		}
	}

	return nil
}
