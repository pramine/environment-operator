package git

import (
	"errors"
	"fmt"

	git2go "gopkg.in/libgit2/git2go.v24"
)

// Pull performs git pull for remote path
func (g *Git) Pull() error {

	repo, err := git2go.OpenRepository(g.LocalPath)
	if err != nil {
		return err
	}

	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		return err
	}

	remoteBranch, err := repo.References.Lookup("refs/remotes/origin/" + g.BranchName)
	if err != nil {
		return err
	}
	remoteBranchID := remoteBranch.Target()

	if err = remote.Fetch([]string{}, g.fetchOptions(), ""); err != nil {
		return err
	}

	// Get annotated commit
	annotatedCommit, err := repo.AnnotatedCommitFromRef(remoteBranch)
	if err != nil {
		return err
	}

	// Do the merge analysis
	mergeHeads := make([]*git2go.AnnotatedCommit, 1)
	mergeHeads[0] = annotatedCommit
	analysis, _, err := repo.MergeAnalysis(mergeHeads)
	if err != nil {
		return err
	}

	// Get repo head
	head, err := repo.Head()
	if err != nil {
		return err
	}

	if analysis&git2go.MergeAnalysisUpToDate != 0 {
		return nil
	} else if analysis&git2go.MergeAnalysisNormal != 0 {
		// Just merge changes
		if err := repo.Merge([]*git2go.AnnotatedCommit{annotatedCommit}, nil, nil); err != nil {
			return err
		}
		// Check for conflicts
		index, err := repo.Index()
		if err != nil {
			return err
		}

		if index.HasConflicts() {
			return errors.New("Conflicts encountered. Please resolve them.")
		}

		// Make the merge commit
		sig, err := repo.DefaultSignature()
		if err != nil {
			return err
		}

		// Get Write Tree
		treeID, err := index.WriteTree()
		if err != nil {
			return err
		}

		tree, err := repo.LookupTree(treeID)
		if err != nil {
			return err
		}

		localCommit, err := repo.LookupCommit(head.Target())
		if err != nil {
			return err
		}

		remoteCommit, err := repo.LookupCommit(remoteBranchID)
		if err != nil {
			return err
		}

		repo.CreateCommit("HEAD", sig, sig, "", tree, localCommit, remoteCommit)

		// Clean up
		repo.StateCleanup()
	} else if analysis&git2go.MergeAnalysisFastForward != 0 {
		// Fast-forward changes
		// Get remote tree
		remoteTree, err := repo.LookupTree(remoteBranchID)
		if err != nil {
			return err
		}

		// Checkout
		if err = repo.CheckoutTree(remoteTree, nil); err != nil {
			return err
		}

		branchRef, err := repo.References.Lookup("refs/heads/" + g.BranchName)
		if err != nil {
			return err
		}

		// Point branch to the object
		branchRef.SetTarget(remoteBranchID, "")
		if _, err := head.SetTarget(remoteBranchID, ""); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("Unexpected merge analysis result %d", analysis)
	}

	return nil
}
