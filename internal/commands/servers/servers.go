package servers

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

// NewServersCommand creates the servers command with all subcommands
func NewServersCommand() *cobra.Command {
	serversCmd := &cobra.Command{
		Use:   "servers",
		Short: "Manage DICOMweb servers",
		Long:  `List and manage DICOMweb server configurations in the Orthanc server.`,
	}

	// Add subcommands
	serversCmd.AddCommand(NewListCommand())
	serversCmd.AddCommand(NewGetCommand())
	serversCmd.AddCommand(NewCreateCommand())
	serversCmd.AddCommand(NewUpdateCommand())
	serversCmd.AddCommand(NewRemoveCommand())

	return serversCmd
}
