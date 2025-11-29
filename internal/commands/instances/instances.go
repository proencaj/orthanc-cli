package instances

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

// NewInstancesCommand creates the instances command with all subcommands
func NewInstancesCommand() *cobra.Command {
	instancesCmd := &cobra.Command{
		Use:   "instances",
		Short: "Manage Orthanc instances",
		Long:  `Query, list, and manage DICOM instances in the Orthanc server.`,
	}

	// Add subcommands
	instancesCmd.AddCommand(NewListCommand())
	instancesCmd.AddCommand(NewGetCommand())
	instancesCmd.AddCommand(NewRemoveCommand())
	instancesCmd.AddCommand(NewDownloadCommand())
	instancesCmd.AddCommand(NewUploadCommand())
	instancesCmd.AddCommand(NewAnonymizeCommand())

	return instancesCmd
}
