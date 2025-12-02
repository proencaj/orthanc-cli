package modalities

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

// NewModalitiesCommand creates the modalities command with all subcommands
func NewModalitiesCommand() *cobra.Command {
	modalitiesCmd := &cobra.Command{
		Use:   "modalities",
		Short: "Manage Orthanc DICOM modalities",
		Long:  `List and manage DICOM modality configurations in the Orthanc server.`,
	}

	// Add subcommands
	modalitiesCmd.AddCommand(NewListCommand())
	modalitiesCmd.AddCommand(NewGetCommand())
	modalitiesCmd.AddCommand(NewCreateCommand())
	modalitiesCmd.AddCommand(NewUpdateCommand())
	modalitiesCmd.AddCommand(NewRemoveCommand())
	modalitiesCmd.AddCommand(NewEchoCommand())
	modalitiesCmd.AddCommand(NewFindCommand())
	modalitiesCmd.AddCommand(NewMoveCommand())
	modalitiesCmd.AddCommand(NewStoreCommand())
	modalitiesCmd.AddCommand(NewRetrieveCommand())

	return modalitiesCmd
}
