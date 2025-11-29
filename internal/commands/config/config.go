package config

import (
	"github.com/spf13/cobra"
)

// NewConfigCommand creates the config command with all subcommands
func NewConfigCommand() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
		Long:  `Manage the Orthanc CLI configuration including server URL, credentials, and TLS settings.`,
	}

	// Add subcommands
	configCmd.AddCommand(NewInitCommand())
	configCmd.AddCommand(NewSetCommand())
	configCmd.AddCommand(NewGetCommand())
	configCmd.AddCommand(NewListCommand())

	return configCmd
}
