package config

import (
	"fmt"
	"os"
	"path/filepath"

	internalConfig "github.com/proencaj/orthanc-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewUseContextCommand creates the config use-context command
func NewUseContextCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use-context <name>",
		Short: "Switch to a different context",
		Long:  `Set the specified context as the current active context.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]

			// Load the config
			cfgFile := viper.ConfigFileUsed()
			cfg, err := internalConfig.LoadConfig(cfgFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Check if context exists
			if _, exists := cfg.Contexts[contextName]; !exists {
				return fmt.Errorf("context %q not found", contextName)
			}

			// Update current context
			cfg.CurrentContext = contextName

			// Determine config file path
			configFile := viper.ConfigFileUsed()
			if configFile == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get home directory: %w", err)
				}
				configFile = filepath.Join(home, ".orthanc-cli.yaml")
			}

			// Save the updated config
			if err := internalConfig.SaveConfigToFile(cfg, configFile); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("Switched to context %q\n", contextName)
			return nil
		},
	}

	return cmd
}
