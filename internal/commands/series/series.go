package series

import (
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

// NewSeriesCommand creates the series command with all subcommands
func NewSeriesCommand() *cobra.Command {
	seriesCmd := &cobra.Command{
		Use:   "series",
		Short: "Manage Orthanc series",
		Long:  `Query, list, and manage DICOM series in the Orthanc server.`,
	}

	// Add subcommands
	seriesCmd.AddCommand(NewListCommand())
	seriesCmd.AddCommand(NewGetCommand())
	seriesCmd.AddCommand(NewRemoveCommand())
	seriesCmd.AddCommand(NewAnonymizeCommand())
	seriesCmd.AddCommand(NewArchiveCommand())
	seriesCmd.AddCommand(NewListInstancesCommand())

	return seriesCmd
}
