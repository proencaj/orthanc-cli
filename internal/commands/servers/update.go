package servers

import (
	"github.com/spf13/cobra"
)

// NewUpdateCommand creates the servers update command
func NewUpdateCommand() *cobra.Command {
	flags := &CreateFlags{}

	command := &cobra.Command{
		Use:   "update <server-name>",
		Short: "Update a DICOMweb server configuration",
		Long:  `Update an existing DICOMweb server configuration in the Orthanc server. You can either provide a JSON file with --file or specify individual parameters.`,
		Example: `  # Update server from JSON file
  orthanc servers update my-pacs --file server.json

  # Update specific fields
  orthanc servers update my-pacs --url https://new-pacs.example.com/dicom-web

  # Update with all options
  orthanc servers update my-pacs \
    --url https://pacs.example.com/dicom-web \
    --username newadmin \
    --password newsecret`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			// Reuse the same logic as create since the API endpoint is the same
			return runCreate(args[0], flags)
		},
	}

	// Add flags (same as create)
	command.Flags().StringVar(&flags.file, "file", "", "JSON file containing server configuration")
	command.Flags().StringVar(&flags.url, "url", "", "URL of the remote DICOMweb server")
	command.Flags().StringVar(&flags.username, "username", "", "Username for authentication (optional)")
	command.Flags().StringVar(&flags.password, "password", "", "Password for authentication (optional)")
	command.Flags().BoolVar(&flags.hasDelete, "has-delete", false, "Whether the server supports DELETE operations")
	command.Flags().BoolVar(&flags.chunkedTransfers, "chunked-transfers", true, "Whether to use chunked transfers (set to false for Orthanc <= 1.5.6)")
	command.Flags().BoolVar(&flags.hasWadoRsUniversalTransferSyntax, "has-wado-rs-universal-transfer-syntax", true, "Whether the server supports WADO-RS universal transfer syntax (set to false for Orthanc DICOMweb plugin <= 1.0)")

	return command
}
