package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jonasvinther/nomad-gitops-operator/pkg/nomad"
)

func init() {
	rootCmd.AddCommand(bootstrapCmd)
}

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap [git repo]",
	Short: "bootstrap a yaml file into a Vault instance",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := args[0]

		fmt.Println(repo)

		nomad.Apply()

		return nil
	},
}
