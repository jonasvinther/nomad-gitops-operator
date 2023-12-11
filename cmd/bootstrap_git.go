package cmd

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"

	"nomad-gitops-operator/pkg/reconcile"
	"nomad-gitops-operator/pkg/repository"
)

type gitFlags struct {
	url         string
	branch      string
	path        string
	var_path    string
	username    string
	password    string
	sshkey      string
	sshinsecure bool
	watch       bool
	delete      bool
}

var gitArgs gitFlags

func init() {
	bootstrapCmd.AddCommand(bootstrapGitCmd)
	bootstrapGitCmd.Flags().StringVar(&gitArgs.url, "url", "", "git repository URL")
	bootstrapGitCmd.Flags().StringVar(&gitArgs.branch, "branch", "main", "git branch")
	bootstrapGitCmd.Flags().StringVar(&gitArgs.path, "path", "**/*.nomad", "glob pattern relative to the repository root")
	bootstrapGitCmd.Flags().StringVar(&gitArgs.var_path, "var-path", "**/*.vars.yml", "var glob pattern relative to the repository root")
	bootstrapGitCmd.Flags().StringVar(&gitArgs.username, "username", "git", "SSH username")
	bootstrapGitCmd.Flags().StringVar(&gitArgs.username, "password", "", "SSH private key password")
	bootstrapGitCmd.Flags().StringVar(&gitArgs.sshkey, "ssh-key", "", "SSH private key")
	bootstrapGitCmd.Flags().BoolVar(&gitArgs.sshinsecure, "ssh-insecure-ignore-host-key", false, "Ignore insecure SSH host key")
	bootstrapGitCmd.Flags().BoolVar(&gitArgs.watch, "watch", true, "Enable watch mode")
	bootstrapGitCmd.Flags().BoolVar(&gitArgs.watch, "delete", true, "Enable delete missing jobs")
}

var bootstrapGitCmd = &cobra.Command{
	Use:   "git [git repo]",
	Short: "Bootstrap Nomad using a git repository",
	Long:  ``,
	// Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return reconcile.Run(reconcile.ReconcileOptions{
			Path:    gitArgs.path,
			VarPath: gitArgs.var_path,
			Watch:   gitArgs.watch,
			Delete:  gitArgs.delete,
			Fs: func() (billy.Filesystem, error) {
				repositoryURL, err := url.Parse(gitArgs.url)
				if err != nil {
					return nil, err
				}

				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Minute*5))
				defer cancel()

				worktree, err := repository.CLone(ctx, repositoryURL, gitArgs.branch, gitArgs.username, gitArgs.sshkey, gitArgs.password, gitArgs.sshinsecure)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
				}

				worktree.Pull(&git.PullOptions{RemoteName: "origin"})

				fs := worktree.Filesystem
				return fs, nil
			}})
	},
}
