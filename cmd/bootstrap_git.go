package cmd

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"

	"nomad-gitops-operator/pkg/nomad"
	"nomad-gitops-operator/pkg/repository"
)

type gitFlags struct {
	url         string
	branch      string
	path        string
	username    string
	password    string
	sshkey      string
	sshinsecure bool
}

var gitArgs gitFlags

func init() {
	bootstrapCmd.AddCommand(bootstrapGitCmd)
	bootstrapGitCmd.Flags().StringVar(&gitArgs.url, "url", "", "git repository URL")
	bootstrapGitCmd.Flags().StringVar(&gitArgs.branch, "branch", "main", "git branch")
	bootstrapGitCmd.Flags().StringVar(&gitArgs.path, "path", "**/*.nomad", "glob pattern relative to the repository root")
	bootstrapGitCmd.Flags().StringVar(&gitArgs.username, "username", "git", "SSH username")
	bootstrapGitCmd.Flags().StringVar(&gitArgs.username, "password", "", "SSH private key password")
	bootstrapGitCmd.Flags().StringVar(&gitArgs.sshkey, "ssh-key", "", "SSH private key")
	bootstrapGitCmd.Flags().BoolVar(&gitArgs.sshinsecure, "ssh-insecure-ignore-host-key", false, "Ignore insecure SSH host key")
}

var bootstrapGitCmd = &cobra.Command{
	Use:   "git [git repo]",
	Short: "Bootstrap Nomad using a git repository",
	Long:  ``,
	// Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create Nomad client
		client, err := nomad.NewClient()
		if err != nil {
			fmt.Printf("Error %s\n", err)
		}

		// Reconcile
		for true {
			repositoryURL, err := url.Parse(gitArgs.url)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Minute*5))
			defer cancel()

			worktree, err := repository.CLone(ctx, repositoryURL, gitArgs.branch, gitArgs.username, gitArgs.sshkey, gitArgs.password, gitArgs.sshinsecure)

			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
			worktree.Pull(&git.PullOptions{RemoteName: "origin"})

			fs := worktree.Filesystem
			files, err := util.Glob(fs, gitArgs.path)
			if err != nil {
				return err
			}

			desiredStateJobs := make(map[string]interface{})

			// Parse and apply all jobs from within the git repo
			for _, filePath := range files {
				f, err := fs.Open(filePath)
				if err != nil {
					return err
				}

				b, err := io.ReadAll(f)
				if err != nil {
					return err
				}

				// Parse job
				job, err := client.ParseJob(string(b))
				if err != nil {
					// If a parse error occurs we skip the job an continue with the next job
					fmt.Printf("Failed to parse file [%s]: %s\n", filePath, err)
					continue
				}
				desiredStateJobs[*job.Name] = job

				// Apply job
				fmt.Printf("Applying job [%s][%s]\n", *job.Name, filePath)
				_, err = client.ApplyJob(job)
				if err != nil {
					return err
				}
			}

			// List all jobs managed by Monoporator
			currentStateJobs, err := client.ListJobs()
			if err != nil {
				fmt.Printf("Error %s\n", err)
			}

			// Check if job has the required metadata
			// Check if job is one of the parsed jobs
			for _, job := range currentStateJobs {
				meta := job.Meta

				if _, isManaged := meta["nomoporater"]; isManaged {
					// If the job is managed by Nomoporator and is part of the desired state
					if _, inDesiredState := desiredStateJobs[*job.Name]; inDesiredState {

					} else {
						fmt.Printf("Deleting job [%s]\n", *job.Name)
						err = client.DeleteJob(job)
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}

			time.Sleep(30 * time.Second)
		}

		return nil
	},
}
