package config

import (
	"fmt"
	"os"
	"path/filepath"

	internalConfig "github.com/proencaj/orthanc-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewSetContextCommand creates the config set-context command
func NewSetContextCommand() *cobra.Command {
	var (
		url      string
		username string
		password string
		insecure bool
		current  bool
	)

	cmd := &cobra.Command{
		Use:   "set-context <name>",
		Short: "Create or update a context",
		Long: `Create a new context or update an existing context with the specified settings.

Examples:
  orthanc config set-context local --url http://localhost:8042 --username orthanc --password orthanc
  orthanc config set-context prod --url https://orthanc.prod.com --username admin --password secret
  orthanc config set-context dev --url http://dev:8042 --insecure --current`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]

			// Load the config
			cfgFile := viper.ConfigFileUsed()
			cfg, err := internalConfig.LoadConfig(cfgFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Initialize contexts map if needed
			if cfg.Contexts == nil {
				cfg.Contexts = make(map[string]*internalConfig.ContextConfig)
			}

			// Get or create context
			ctx, exists := cfg.Contexts[contextName]
			if !exists {
				ctx = &internalConfig.ContextConfig{}
				cfg.Contexts[contextName] = ctx
			}

			// Update fields if provided
			if cmd.Flags().Changed("url") {
				ctx.Orthanc.URL = url
			}
			if cmd.Flags().Changed("username") {
				ctx.Orthanc.Username = username
			}
			if cmd.Flags().Changed("password") {
				ctx.Orthanc.Password = password
			}
			if cmd.Flags().Changed("insecure") {
				ctx.Orthanc.Insecure = insecure
			}

			// Set as current context if requested or if it's the only context
			if current || len(cfg.Contexts) == 1 {
				cfg.CurrentContext = contextName
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

			if exists {
				fmt.Printf("Updated context %q\n", contextName)
			} else {
				fmt.Printf("Created context %q\n", contextName)
			}

			if current || len(cfg.Contexts) == 1 {
				fmt.Printf("Set as current context\n")
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "Orthanc server URL")
	cmd.Flags().StringVar(&username, "username", "", "Orthanc username")
	cmd.Flags().StringVar(&password, "password", "", "Orthanc password")
	cmd.Flags().BoolVar(&insecure, "insecure", false, "Skip TLS verification")
	cmd.Flags().BoolVar(&current, "current", false, "Set as current context")

	return cmd
}
