package repository

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"golang.org/x/crypto/ssh"
)

// CLone will clone a git repo
func CLone(ctx context.Context, repositoryURL *url.URL, branch string, username string, sshkey string, password string, sshinsecure bool) (*git.Worktree, error) {
	storer := memory.NewStorage()
	fs := memfs.New()
	ranchRef := plumbing.NewBranchReferenceName(branch)

	cloneOptions := &git.CloneOptions{
		URL:           repositoryURL.String(),
		ReferenceName: ranchRef,
		SingleBranch:  true,
		Depth:         1,
		Tags:          git.NoTags,
	}

	if sshkey != "" {
		publicKey, err := gitssh.NewPublicKeys(username, []byte(sshkey), password)
		if err != nil {
			return nil, fmt.Errorf("unable to use ssh key for git auth: %w", err)
		}
		if sshinsecure {
			publicKey.HostKeyCallback = ssh.InsecureIgnoreHostKey()
		}
		cloneOptions.Auth = publicKey
	}

	repo, err := git.CloneContext(ctx, storer, fs, cloneOptions)
	if err != nil {
		return nil, fmt.Errorf("unable to clone repository %s for templating: %w", repositoryURL.String(), err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("unable to get worktree: %w", err)
	}

	return w, nil
}
