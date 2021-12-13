package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"nomad-gitops-operator/pkg/nomad"
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

		_, err := nomad.Apply()

		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}

		return nil
	},
}
