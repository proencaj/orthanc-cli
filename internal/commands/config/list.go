package config

import (
	"fmt"

	internalConfig "github.com/proencaj/orthanc-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewListCommand creates the config list command
func NewListCommand() *cobra.Command {
	var showPassword bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all configuration values from the current context",
		Long:  `Display all current configuration values from the current context.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load the config
			cfgFile := viper.ConfigFileUsed()
			cfg, err := internalConfig.LoadConfig(cfgFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			configFile := viper.ConfigFileUsed()
			if configFile != "" {
				fmt.Printf("Configuration file: %s\n", configFile)
			} else {
				fmt.Println("No configuration file loaded")
			}

			if cfg.CurrentContext != "" {
				fmt.Printf("Current context: %s\n\n", cfg.CurrentContext)
			} else {
				fmt.Println("Current context: (not set)\n")
			}

			// Get current context config
			if cfg.CurrentContext == "" {
				fmt.Println("No current context set")
				fmt.Println("\nCreate a context with: orthanc config set-context <name> --url <url>")
				return nil
			}

			orthancCfg, err := cfg.GetCurrentContext()
			if err != nil {
				return fmt.Errorf("failed to get current context: %w", err)
			}

			fmt.Println("Orthanc Configuration:")
			fmt.Println("----------------------")

			if orthancCfg.URL != "" {
				fmt.Printf("  URL:      %s\n", orthancCfg.URL)
			} else {
				fmt.Println("  URL:      (not set)")
			}

			if orthancCfg.Username != "" {
				fmt.Printf("  Username: %s\n", orthancCfg.Username)
			} else {
				fmt.Println("  Username: (not set)")
			}

			if orthancCfg.Password != "" {
				if showPassword {
					fmt.Printf("  Password: %s\n", orthancCfg.Password)
				} else {
					fmt.Println("  Password: ********")
				}
			} else {
				fmt.Println("  Password: (not set)")
			}

			fmt.Printf("  Insecure: %v\n", orthancCfg.Insecure)

			fmt.Println()
			fmt.Println("Output Configuration:")
			fmt.Println("---------------------")
			fmt.Printf("  JSON:     %v\n", cfg.Output.JSON)

			return nil
		},
	}

	cmd.Flags().BoolVar(&showPassword, "show-password", false, "Show password in plain text")

	return cmd
}
