package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewGetCommand creates the config get command
func NewGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Long: `Get a configuration value from the config file.

Available keys:
  orthanc.url       - Orthanc server URL
  orthanc.username  - Orthanc username
  orthanc.password  - Orthanc password
  orthanc.insecure  - Skip TLS verification
  output.json       - Default JSON output

Examples:
  orthanc config get orthanc.url
  orthanc config get orthanc.username
  orthanc config get output.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]

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

			value := viper.Get(key)
			if value == nil {
				fmt.Printf("%s is not set\n", key)
				return nil
			}

			// Mask password for security
			if key == "orthanc.password" {
				fmt.Printf("%s = ********\n", key)
			} else {
				fmt.Printf("%s = %v\n", key, value)
			}

			return nil
		},
	}

	return cmd
}
