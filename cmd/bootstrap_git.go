package cmd

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"

	"nomad-gitops-operator/pkg/nomad"
	"nomad-gitops-operator/pkg/repository"
)

type gitFlags struct {
	url      string
	branch   string
	path     string
	username string
	password string
}

var gitArgs gitFlags

func init() {
	bootstrapCmd.AddCommand(bootstrapGitCmd)
	bootstrapGitCmd.Flags().StringVar(&gitArgs.url, "url", "", "git repository URL")
	bootstrapGitCmd.Flags().StringVar(&gitArgs.branch, "branch", "main", "git branch")
	bootstrapGitCmd.Flags().StringVar(&gitArgs.path, "path", "/", "path relative to the repository root")
}

var bootstrapGitCmd = &cobra.Command{
	Use:   "git [git repo]",
	Short: "Bootstrap Nomad using a git repository",
	Long:  ``,
	// Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repositoryURL, err := url.Parse(gitArgs.url)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Minute*5))
		defer cancel()

		worktree, err := repository.CLone(ctx, repositoryURL, gitArgs.branch)

		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}

		// Reconcile
		for true {
			worktree.Pull(&git.PullOptions{RemoteName: "origin"})

			fs := worktree.Filesystem
			path := gitArgs.path
			files, err := fs.ReadDir(path)
			if err != nil {
				return err
			}

			for _, file := range files {
				filePath := fs.Join(path, file.Name())
				f, err := fs.Open(filePath)
				if err != nil {
					return err
				}

				b, err := io.ReadAll(f)
				if err != nil {
					return err
				}

				status, err := nomad.ApplyJob(string(b))
				if err != nil {
					return err
				}
				fmt.Println(status)
			}
			time.Sleep(30 * time.Second)
		}

		return nil
	},
}
