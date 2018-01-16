package git

// Pull performs git pull for remote path
func (g *Git) Pull() error {
	tree, err := g.Repository.Worktree()
	if err != nil {
		return err
	}

	return tree.Pull(g.pullOptions())
}
