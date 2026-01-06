package config

import (
	"fmt"

	internalConfig "github.com/proencaj/orthanc-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewCurrentContextCommand creates the config current-context command
func NewCurrentContextCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current-context",
		Short: "Display the current context",
		Long:  `Show the currently active context name.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load the config to get current context
			cfgFile := viper.ConfigFileUsed()
			cfg, err := internalConfig.LoadConfig(cfgFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			if cfg.CurrentContext == "" {
				fmt.Println("No current context set")
				return nil
			}

			fmt.Println(cfg.CurrentContext)
			return nil
		},
	}

	return cmd
}
