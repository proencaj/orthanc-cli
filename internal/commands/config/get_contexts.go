package config

import (
	"fmt"
	"sort"

	internalConfig "github.com/proencaj/orthanc-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewGetContextsCommand creates the config get-contexts command
func NewGetContextsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-contexts",
		Short: "List all available contexts",
		Long:  `Display all available contexts with the current context marked.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load the config to get contexts
			cfgFile := viper.ConfigFileUsed()
			cfg, err := internalConfig.LoadConfig(cfgFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			if len(cfg.Contexts) == 0 {
				fmt.Println("No contexts found")
				fmt.Println("\nCreate a context with: orthanc config set-context <name> --url <url>")
				return nil
			}

			// Sort context names for consistent output
			names := make([]string, 0, len(cfg.Contexts))
			for name := range cfg.Contexts {
				names = append(names, name)
			}
			sort.Strings(names)

			fmt.Println("CURRENT   NAME")
			for _, name := range names {
				current := "  "
				if name == cfg.CurrentContext {
					current = "* "
				}
				fmt.Printf("%s        %s\n", current, name)
			}

			return nil
		},
	}

	return cmd
}
