package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	internalConfig "github.com/proencaj/orthanc-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewSetCommand creates the config set command
func NewSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value in the current context",
		Long: `Set a configuration value in the current context.

Available keys:
  orthanc.url       - Orthanc server URL (e.g., http://localhost:8042)
  orthanc.username  - Orthanc username
  orthanc.password  - Orthanc password
  orthanc.insecure  - Skip TLS verification (true/false)
  output.json       - Output in JSON format by default (true/false)

Examples:
  orthanc config set orthanc.url http://localhost:8042
  orthanc config set orthanc.username myuser
  orthanc config set orthanc.password mypassword
  orthanc config set orthanc.insecure false
  orthanc config set output.json true`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			// Validate key
			validKeys := map[string]bool{
				"orthanc.url":      true,
				"orthanc.username": true,
				"orthanc.password": true,
				"orthanc.insecure": true,
				"output.json":      true,
			}

			if !validKeys[key] {
				return fmt.Errorf("invalid configuration key: %s\nValid keys: orthanc.url, orthanc.username, orthanc.password, orthanc.insecure, output.json", key)
			}

			// Load the config
			cfgFile := viper.ConfigFileUsed()
			cfg, err := internalConfig.LoadConfig(cfgFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// For context-specific keys, update current context
			if key != "output.json" {
				if cfg.CurrentContext == "" {
					return fmt.Errorf("no current context set (use 'orthanc config set-context <name>' to create one)")
				}

				ctx, exists := cfg.Contexts[cfg.CurrentContext]
				if !exists {
					return fmt.Errorf("current context %q not found", cfg.CurrentContext)
				}

				switch key {
				case "orthanc.url":
					ctx.Orthanc.URL = value
				case "orthanc.username":
					ctx.Orthanc.Username = value
				case "orthanc.password":
					ctx.Orthanc.Password = value
				case "orthanc.insecure":
					boolValue, err := strconv.ParseBool(value)
					if err != nil {
						return fmt.Errorf("invalid boolean value for %s: %s (use true or false)", key, value)
					}
					ctx.Orthanc.Insecure = boolValue
				}
			} else {
				// output.json is global
				boolValue, err := strconv.ParseBool(value)
				if err != nil {
					return fmt.Errorf("invalid boolean value for %s: %s (use true or false)", key, value)
				}
				cfg.Output.JSON = boolValue
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

			// Mask password for security
			displayValue := value
			if key == "orthanc.password" {
				displayValue = "********"
			}

			fmt.Printf("✓ Set %s = %s\n", key, displayValue)
			if key != "output.json" {
				fmt.Printf("  Context: %s\n", cfg.CurrentContext)
			}
			fmt.Printf("  Config file: %s\n", configFile)

			return nil
		},
	}

	return cmd
}
