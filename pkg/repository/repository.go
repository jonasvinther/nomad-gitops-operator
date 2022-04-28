package repository

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

// CLone will clone a git repo
func CLone(ctx context.Context, repositoryURL *url.URL, branch string) (*git.Worktree, error) {
	storer := memory.NewStorage()
	fs := memfs.New()
	ranchRef := plumbing.NewBranchReferenceName(branch)
	repo, err := git.CloneContext(ctx, storer, fs, &git.CloneOptions{
		// Auth: auth,
		URL:           repositoryURL.String(),
		ReferenceName: ranchRef,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to clone repository %s for templating: %w", repositoryURL.String(), err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("unable to get worktree: %w", err)
	}

	return w, nil
}
