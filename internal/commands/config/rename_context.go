package config

import (
	"fmt"
	"os"
	"path/filepath"

	internalConfig "github.com/proencaj/orthanc-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewRenameContextCommand creates the config rename-context command
func NewRenameContextCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename-context <old-name> <new-name>",
		Short: "Rename a context",
		Long:  `Rename an existing context to a new name.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			oldName := args[0]
			newName := args[1]

			// Load the config
			cfgFile := viper.ConfigFileUsed()
			cfg, err := internalConfig.LoadConfig(cfgFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Check if old context exists
			ctx, exists := cfg.Contexts[oldName]
			if !exists {
				return fmt.Errorf("context %q not found", oldName)
			}

			// Check if new name already exists
			if _, exists := cfg.Contexts[newName]; exists {
				return fmt.Errorf("context %q already exists", newName)
			}

			// Rename the context
			cfg.Contexts[newName] = ctx
			delete(cfg.Contexts, oldName)

			// Update current context if needed
			if cfg.CurrentContext == oldName {
				cfg.CurrentContext = newName
			}

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

			fmt.Printf("Renamed context %q to %q\n", oldName, newName)
			return nil
		},
	}

	return cmd
}
