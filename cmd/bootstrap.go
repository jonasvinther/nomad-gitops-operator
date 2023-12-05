package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(bootstrapCmd)
}

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "bootstrap into a Nomad instance",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
}
