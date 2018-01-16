package git

import (
	"fmt"
	"os"
	"testing"
)

func TestPull(t *testing.T) {

	remotePath := createTestRepo(t)
	localPath := createSrcPath(t)
	//defer cleanupTestPath(localPath)
	//defer cleanupTestPath(remotePath)

	fmt.Printf("local path: %s\nremote path: %s", localPath, remotePath)

	g := initAndClone(t, localPath, remotePath)

	commitTestJunk(t, remotePath, "zzz.bitesize")

	if err := g.Pull(); err != nil {
		t.Errorf("Error on pull: %s", err.Error())
	}

	if _, err := os.Stat(localPath + "/zzz.bitesize"); os.IsNotExist(err) {
		t.Error("File zzz.bitesize is missing in cloned repo")
	}

}

func TestPullWithLocalChanges(t *testing.T) {
	remotePath := createTestRepo(t)
	localPath := createSrcPath(t)
	defer cleanupTestPath(remotePath)
	defer cleanupTestPath(localPath)

	g := initAndClone(t, localPath, remotePath)
	commitTestJunk(t, localPath, "zee.bitesize")
	commitTestJunk(t, remotePath, "zzz.bitesize")
	// commitTestJunk(t, srcPath, "zee.bitesize")

	if err := g.Pull(); err != nil {
		t.Errorf("Error on pull: %s", err.Error())
	}

	if _, err := os.Stat(localPath + "/zzz.bitesize"); os.IsNotExist(err) {
		t.Error("File zzz.bitesize is missing in cloned repo")
	}
}
