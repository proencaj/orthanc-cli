package modalities

import (
	"github.com/spf13/cobra"
)

// NewUpdateCommand creates the modalities update command
func NewUpdateCommand() *cobra.Command {
	flags := &CreateFlags{}

	command := &cobra.Command{
		Use:   "update <modality-name>",
		Short: "Update a DICOM modality configuration",
		Long:  `Update an existing DICOM modality configuration in the Orthanc server. You can either provide a JSON file with --file or specify individual parameters.`,
		Example: `  # Update modality from JSON file
  orthanc modalities update PACS_SERVER --file modality.json

  # Update specific fields
  orthanc modalities update PACS_SERVER --host 192.168.1.200 --port 11113

  # Update with all options
  orthanc modalities update MY_MODALITY \
    --aet NEW_AET \
    --host localhost \
    --port 4242 \
    --manufacturer "GE Healthcare" \
    --timeout 30`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			// Reuse the same logic as create since the API endpoint is the same
			return runCreate(args[0], flags)
		},
	}

	// Add flags (same as create)
	command.Flags().StringVar(&flags.file, "file", "", "JSON file containing modality configuration")
	command.Flags().StringVar(&flags.aet, "aet", "", "Application Entity Title (AET) of the remote modality")
	command.Flags().StringVar(&flags.host, "host", "", "Host/IP address of the remote modality")
	command.Flags().IntVar(&flags.port, "port", 0, "Port number of the remote modality")
	command.Flags().StringVar(&flags.manufacturer, "manufacturer", "", "Manufacturer name (optional)")
	command.Flags().BoolVar(&flags.allowEcho, "allow-echo", true, "Allow DICOM C-ECHO requests")
	command.Flags().BoolVar(&flags.allowFind, "allow-find", true, "Allow DICOM C-FIND requests")
	command.Flags().BoolVar(&flags.allowGet, "allow-get", true, "Allow DICOM C-GET requests")
	command.Flags().BoolVar(&flags.allowMove, "allow-move", true, "Allow DICOM C-MOVE requests")
	command.Flags().BoolVar(&flags.allowStore, "allow-store", true, "Allow DICOM C-STORE requests")
	command.Flags().IntVar(&flags.timeout, "timeout", 10, "Timeout for DICOM operations in seconds")

	return command
}
