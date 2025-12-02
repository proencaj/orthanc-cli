package cmd

import (
	"fmt"
	"os"

	"github.com/proencaj/orthanc-cli/internal/client"
	configCmd "github.com/proencaj/orthanc-cli/internal/commands/config"
	"github.com/proencaj/orthanc-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "orthanc",
	Short: "A CLI tool to interact with Orthanc DICOM servers",
	SilenceUsage: true,
	Long: `orthanc is a command-line interface for managing and querying
Orthanc DICOM servers. It provides commands to interact with instances,
studies, series, patients, and other Orthanc resources.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize configuration
		var err error
		cfg, err = config.LoadConfig(cfgFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.orthanc-cli.yaml)")

	// Register subcommands
	rootCmd.AddCommand(configCmd.NewConfigCommand())
}

// AddCommand adds a command to the root command
func AddCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}

// GetConfig returns the current configuration
func GetConfig() *config.Config {
	return cfg
}

// GetClient creates a new Orthanc client using the current configuration
func GetClient() (*client.Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}
	return client.NewClient(cfg)
}
