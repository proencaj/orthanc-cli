package system

import (
	"encoding/json"
	"fmt"

	"github.com/proencaj/orthanc-cli/internal/client"
	"github.com/proencaj/orthanc-cli/internal/config"
	"github.com/spf13/cobra"
)

// clientGetter is a function type that returns an Orthanc client
var clientGetter func() (*client.Client, error)

// SetClientGetter sets the function to get the Orthanc client
func SetClientGetter(getter func() (*client.Client, error)) {
	clientGetter = getter
}

// getClient returns the Orthanc client using the configured getter
func getClient() (*client.Client, error) {
	if clientGetter != nil {
		return clientGetter()
	}
	// Fallback: try to load config from default location
	cfg, err := config.LoadConfig("")
	if err != nil {
		return nil, err
	}
	return client.NewClient(cfg)
}

// shouldUseJSON checks if JSON output is enabled in config
func shouldUseJSON() bool {
	cfg, err := config.LoadConfig("")
	if err != nil {
		return false
	}
	return cfg.Output.JSON
}

// SystemFlags holds the flags for the system command
type SystemFlags struct {
	jsonOutput bool
}

// NewSystemCommand creates the system command
func NewSystemCommand() *cobra.Command {
	flags := &SystemFlags{}

	command := &cobra.Command{
		Use:   "system",
		Short: "Get Orthanc server system information",
		Long:  `Retrieve system information about the Orthanc server including version, ports, plugins, and database details.`,
		Example: `  # Get system information
  orthanc system

  # Get system information in JSON format
  orthanc system --json`,
		Args: cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return runSystem(flags)
		},
	}

	// Add flags
	command.Flags().BoolVar(&flags.jsonOutput, "json", false, "Output in JSON format")

	return command
}

func runSystem(flags *SystemFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Check if JSON output should be used (flag or config)
	jsonOutput := flags.jsonOutput || shouldUseJSON()

	// Get system information
	systemInfo, err := client.GetSystem()
	if err != nil {
		return fmt.Errorf("failed to get system information: %w", err)
	}

	// Display results
	if jsonOutput {
		data, err := json.MarshalIndent(systemInfo, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Human-readable output
	fmt.Println("Orthanc System Information")
	fmt.Println("==========================")
	fmt.Println()

	fmt.Printf("Server Name:         %s\n", systemInfo.Name)
	fmt.Printf("Version:             %s\n", systemInfo.Version)
	fmt.Printf("API Version:         %d\n", systemInfo.ApiVersion)
	fmt.Println()

	fmt.Println("Network Configuration:")
	fmt.Printf("  HTTP Port:         %d\n", systemInfo.HttpPort)
	fmt.Printf("  DICOM AET:         %s\n", systemInfo.DicomAet)
	fmt.Printf("  DICOM Port:        %d\n", systemInfo.DicomPort)
	fmt.Println()

	fmt.Println("Database:")
	fmt.Printf("  Database Version:  %d\n", systemInfo.DatabaseVersion)
	if systemInfo.DatabaseBackendPlugin != "" {
		fmt.Printf("  Backend Plugin:    %s\n", systemInfo.DatabaseBackendPlugin)
	} else {
		fmt.Printf("  Backend:           SQLite (built-in)\n")
	}
	if systemInfo.InMemoryDatabaseIdentifier != "" {
		fmt.Printf("  In-Memory ID:      %s\n", systemInfo.InMemoryDatabaseIdentifier)
	}
	fmt.Println()

	fmt.Println("Storage:")
	if systemInfo.StorageAreaPlugin != "" {
		fmt.Printf("  Storage Plugin:    %s\n", systemInfo.StorageAreaPlugin)
	} else {
		fmt.Printf("  Storage:           File system (built-in)\n")
	}
	if systemInfo.MaximumStorageSize > 0 {
		fmt.Printf("  Max Storage Size:  %d bytes (%.2f GB)\n",
			systemInfo.MaximumStorageSize,
			float64(systemInfo.MaximumStorageSize)/(1024*1024*1024))
	} else {
		fmt.Printf("  Max Storage Size:  Unlimited\n")
	}
	fmt.Println()

	fmt.Println("Features:")
	fmt.Printf("  Plugins Enabled:   %v\n", systemInfo.PluginsEnabled)
	fmt.Printf("  Check Revisions:   %v\n", systemInfo.CheckRevisions)

	return nil
}
