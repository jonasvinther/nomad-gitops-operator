package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "nomoperator",
	Short: "Nomoperator is a cli tool for keeping a HashiCorp Nomad cluster in sync with a Git repository containing Nomad job specifications.",
	Long: `Nomoperator is a cli tool for keeping a HashiCorp Nomad cluster in sync with a Git repository containing Nomad job specifications.
Created by Jonas Vinther.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Apply the viper config value to the flag when the flag is not set and viper has a value
		address, _ := cmd.Flags().GetString("address")
		if address != "" {
			os.Setenv("NOMAD_ADDR", address)
		}

		if viper.IsSet("NOMAD_ADDR") && address == "" {
			value := viper.Get("NOMAD_ADDR").(string)
			err := cmd.Flags().Set("address", value)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("address", "a", "", "Address of the Nomad server")

	// AutomaticEnv makes Viper load environment variables
	viper.AutomaticEnv()

	// Explicitly defines the path, name and type of the config file.
	// Viper will use this and not check any of the config paths.
	// It will search for the "config" file in the ~/.nomoperator
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.nomoperator")
	viper.SetConfigName("config")

	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil {
		// log.Fatalf("Error while reading config file %s", err)
	}

}
