package git

import "testing"

func TestNoUpdatesExist(t *testing.T) {
	local := createSrcPath(t)
	remote := createTestRepo(t)

	defer cleanupTestPath(local)
	defer cleanupTestPath(remote)

	g := initAndClone(t, local, remote)

	if exist, _ := g.UpdatesExist(); exist {
		t.Error("Invalid state: no update should be available")
	}
}

func TestUpdatesExist(t *testing.T) {
	local := createSrcPath(t)
	remote := createTestRepo(t)

	defer cleanupTestPath(local)
	defer cleanupTestPath(remote)

	g := initAndClone(t, local, remote)

	commitTestJunk(t, remote, "")

	if exist, e := g.UpdatesExist(); !exist {
		if e != nil {
			t.Error("Exception on update exist: %s", e.Error())
		} else {
			t.Error("Invalid state: remote update not found.")
		}
	}

}
