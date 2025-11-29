package studies

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

// NewStudiesCommand creates the studies command with all subcommands
func NewStudiesCommand() *cobra.Command {
	studiesCmd := &cobra.Command{
		Use:   "studies",
		Short: "Manage Orthanc studies",
		Long:  `Query, list, and manage DICOM studies in the Orthanc server.`,
	}

	// Add subcommands
	studiesCmd.AddCommand(NewListCommand())
	studiesCmd.AddCommand(NewGetCommand())
	studiesCmd.AddCommand(NewRemoveCommand())
	studiesCmd.AddCommand(NewAnonymizeCommand())
	studiesCmd.AddCommand(NewArchiveCommand())

	return studiesCmd
}
