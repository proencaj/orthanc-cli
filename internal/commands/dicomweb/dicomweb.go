package dicomweb

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

// NewDicomwebCommand creates the dicomweb command with all subcommands
func NewDicomwebCommand() *cobra.Command {
	dicomwebCmd := &cobra.Command{
		Use:   "dicomweb",
		Short: "DICOMweb network operations",
		Long:  `Perform DICOMweb operations including WADO-URI, WADO-RS, and QIDO-RS queries.`,
	}

	// Add subcommands
	dicomwebCmd.AddCommand(NewWadoCommand())
	dicomwebCmd.AddCommand(NewWadoRsCommand())
	dicomwebCmd.AddCommand(NewQidoCommand())

	return dicomwebCmd
}
