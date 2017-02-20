package git

import (
	"os"

	"testing"
)

func TestClone(t *testing.T) {
	remotePath := createTestRepo(t)
	localPath := createSrcPath(t)
	defer cleanupTestPath(remotePath)
	defer cleanupTestPath(localPath)

	initAndClone(t, localPath, remotePath)

	if _, err := os.Stat(localPath + "/environments.bitesize"); os.IsNotExist(err) {
		t.Error("File environments.bitesize is missing in cloned repo")
	}

}
