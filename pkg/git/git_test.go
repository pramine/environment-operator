package git

import (
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	git2go "gopkg.in/libgit2/git2go.v25"
)

func initAndClone(t *testing.T, local, remote string) *Git {
	g := &Git{
		LocalPath:  local,
		RemotePath: remote,
		BranchName: "master",
	}

	if err := g.Clone(); err != nil {
		t.Errorf("Error on clone: %s", err.Error())
	}
	return g
}

func createTestRepo(t *testing.T) string {
	path, err := ioutil.TempDir("", "env-operator-remote")
	checkFatal(t, err)

	repo, err := git2go.InitRepository(path, false)
	checkFatal(t, err)

	tmpfile := "environments.bitesize"
	err = ioutil.WriteFile(path+"/"+tmpfile, []byte("foo\n"), 0644)
	checkFatal(t, err)

	sig := &git2go.Signature{
		Name:  "Author",
		Email: "author@pearson.com",
		When:  time.Now(),
	}

	idx, err := repo.Index()
	checkFatal(t, err)
	err = idx.AddByPath("environments.bitesize")
	checkFatal(t, err)

	err = idx.Write()
	checkFatal(t, err)
	treeID, err := idx.WriteTree()
	checkFatal(t, err)

	message := "Commit message\n"
	tree, err := repo.LookupTree(treeID)
	checkFatal(t, err)
	_, err = repo.CreateCommit("HEAD", sig, sig, message, tree)
	checkFatal(t, err)
	return path
}

func commitTestJunk(t *testing.T, dest string, file string) {

	repo, err := git2go.OpenRepository(dest)
	checkFatal(t, err)

	if file == "" {
		f, e := ioutil.TempFile(dest, "junk")
		file = f.Name()
		checkFatal(t, e)
	}
	path := dest + "/" + file
	contents := randomString(64)

	err = ioutil.WriteFile(path, contents, 0644)
	checkFatal(t, err)

	loc, err := time.LoadLocation("UTC")
	checkFatal(t, err)

	sig := &git2go.Signature{
		Name:  "Author",
		Email: "author@pearson.com",
		When:  time.Date(2013, 03, 06, 14, 30, 0, 0, loc),
	}

	idx, err := repo.Index()
	checkFatal(t, err)

	err = idx.AddByPath(file)
	checkFatal(t, err)

	err = idx.Write()
	checkFatal(t, err)
	treeID, err := idx.WriteTree()
	checkFatal(t, err)

	currentBranch, err := repo.Head()
	checkFatal(t, err)
	currentTip, err := repo.LookupCommit(currentBranch.Target())
	checkFatal(t, err)

	message := "Commit message\n"
	tree, err := repo.LookupTree(treeID)
	checkFatal(t, err)
	_, err = repo.CreateCommit("HEAD", sig, sig, message, tree, currentTip)
	checkFatal(t, err)

}

func checkFatal(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func createSrcPath(t *testing.T) string {
	path, err := ioutil.TempDir("", "env-operator-local")
	checkFatal(t, err)
	return path
}

func cleanupTestPath(path string) {
	os.RemoveAll(path + "/")
}

func randomString(strlen int) []byte {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return result
}
