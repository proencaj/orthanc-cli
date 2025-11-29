package config

import (
	"fmt"

	internalConfig "github.com/proencaj/orthanc-cli/internal/config"
	"github.com/spf13/cobra"
)

// NewInitCommand creates the config init command
func NewInitCommand() *cobra.Command {
	var outputPath string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a default configuration file",
		Long:  `Create a default configuration file at ~/.orthanc-cli.yaml with example settings.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := internalConfig.SaveConfig(outputPath); err != nil {
				return fmt.Errorf("failed to create config file: %w", err)
			}

			if outputPath == "" {
				fmt.Println("✓ Configuration file created at ~/.orthanc-cli.yaml")
			} else {
				fmt.Printf("✓ Configuration file created at %s\n", outputPath)
			}

			fmt.Println("\nEdit the file to configure your Orthanc server settings:")
			fmt.Println("  - url: Your Orthanc server URL")
			fmt.Println("  - username: Your Orthanc username")
			fmt.Println("  - password: Your Orthanc password")
			fmt.Println("  - insecure: Set to true to skip TLS verification (not recommended)")

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output path for config file (default: ~/.orthanc-cli.yaml)")

	return cmd
}
