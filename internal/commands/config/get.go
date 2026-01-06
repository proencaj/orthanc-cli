package config

import (
	"fmt"

	internalConfig "github.com/proencaj/orthanc-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewGetCommand creates the config get command
func NewGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value from the current context",
		Long: `Get a configuration value from the current context.

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

			// For context-specific keys, get from current context
			if key != "output.json" {
				cfgFile := viper.ConfigFileUsed()
				cfg, err := internalConfig.LoadConfig(cfgFile)
				if err != nil {
					return fmt.Errorf("failed to load config: %w", err)
				}

				orthancCfg, err := cfg.GetCurrentContext()
				if err != nil {
					return fmt.Errorf("failed to get current context: %w", err)
				}

				var value interface{}
				switch key {
				case "orthanc.url":
					value = orthancCfg.URL
				case "orthanc.username":
					value = orthancCfg.Username
				case "orthanc.password":
					if orthancCfg.Password == "" {
						fmt.Printf("%s is not set\n", key)
						return nil
					}
					fmt.Printf("%s = ********\n", key)
					return nil
				case "orthanc.insecure":
					value = orthancCfg.Insecure
				}

				if value == "" || value == nil {
					fmt.Printf("%s is not set\n", key)
					return nil
				}

				fmt.Printf("%s = %v\n", key, value)
			} else {
				// output.json is global
				value := viper.Get(key)
				if value == nil {
					fmt.Printf("%s is not set\n", key)
					return nil
				}
				fmt.Printf("%s = %v\n", key, value)
			}

			return nil
		},
	}

	return cmd
}
