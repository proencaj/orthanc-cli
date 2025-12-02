package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewSetCommand creates the config set command
func NewSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Long: `Set a configuration value in the config file.

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

			// Special handling for boolean values
			if key == "orthanc.insecure" || key == "output.json" {
				boolValue, err := strconv.ParseBool(value)
				if err != nil {
					return fmt.Errorf("invalid boolean value for %s: %s (use true or false)", key, value)
				}
				viper.Set(key, boolValue)
			} else {
				viper.Set(key, value)
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

			// Write config
			if err := viper.WriteConfigAs(configFile); err != nil {
				return fmt.Errorf("failed to write config file: %w", err)
			}

			// Mask password for security
			displayValue := value
			if key == "orthanc.password" {
				displayValue = "********"
			}

			fmt.Printf("✓ Set %s = %s\n", key, displayValue)
			fmt.Printf("  Config file: %s\n", configFile)

			return nil
		},
	}

	return cmd
}
