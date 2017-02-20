package git

import (
	"os"

	git2go "gopkg.in/libgit2/git2go.v25"
)

// Git represents repository object and wraps git2go calls
type Git struct {
	SSHKey     string
	LocalPath  string
	RemotePath string
	BranchName string
}

// CloneOrPull checks if repo exists in local path. If it does, it
// pulls changes from remotePath, if it doesn't, performs a full git clone
func (g *Git) CloneOrPull() error {
	if _, err := os.Stat(g.LocalPath); err == nil {
		return g.Pull()
	}
	return g.Clone()
}

func credentialsCallback(url string, username string, allowedTypes git2go.CredType) (git2go.ErrorCode, *git2go.Cred) {
	sshKey := os.Getenv("PRIVATE_SSH_KEY")

	ret, cred := git2go.NewCredSshKeyFromMemory("git", "", sshKey, "")
	return git2go.ErrorCode(ret), &cred
}

// Made this one just return 0 during troubleshooting...
func certificateCheckCallback(cert *git2go.Certificate, valid bool, hostname string) git2go.ErrorCode {
	return 0
}

func cloneOptions() *git2go.CloneOptions {
	opts := &git2go.CloneOptions{}
	opts.FetchOptions = fetchOptions()
	return opts
}

func fetchOptions() *git2go.FetchOptions {
	return &git2go.FetchOptions{
		RemoteCallbacks: git2go.RemoteCallbacks{
			CredentialsCallback:      credentialsCallback,
			CertificateCheckCallback: certificateCheckCallback,
		},
	}
}
