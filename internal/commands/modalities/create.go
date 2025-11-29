package modalities

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/proencaj/gorthanc/types"
	"github.com/spf13/cobra"
)

// CreateFlags holds the flags for the create command
type CreateFlags struct {
	file         string
	aet          string
	host         string
	port         int
	manufacturer string
	allowEcho    bool
	allowFind    bool
	allowGet     bool
	allowMove    bool
	allowStore   bool
	timeout      int
}

// NewCreateCommand creates the modalities create command
func NewCreateCommand() *cobra.Command {
	flags := &CreateFlags{}

	command := &cobra.Command{
		Use:   "create <modality-name>",
		Short: "Create or update a DICOM modality configuration",
		Long:  `Create or update a DICOM modality configuration in the Orthanc server. You can either provide a JSON file with --file or specify individual parameters.`,
		Example: `  # Create modality from JSON file
  orthanc modalities create PACS_SERVER --file modality.json

  # Create modality with individual parameters
  orthanc modalities create PACS_SERVER --aet REMOTE_AET --host 192.168.1.100 --port 11112

  # Create with all options
  orthanc modalities create MY_MODALITY \
    --aet MY_AET \
    --host localhost \
    --port 4242 \
    --manufacturer "GE Healthcare" \
    --timeout 30 \
    --allow-echo \
    --allow-find \
    --allow-store

  # Disable specific permissions
  orthanc modalities create RESTRICTED \
    --aet RESTRICTED_AET \
    --host 10.0.0.50 \
    --port 4242 \
    --allow-store=false \
    --allow-get=false`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runCreate(args[0], flags)
		},
	}

	// Add flags
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

func runCreate(modalityName string, flags *CreateFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	var request *types.ModalityCreateRequest

	// If file is provided, read from file
	if flags.file != "" {
		request, err = readModalityFromFile(flags.file)
		if err != nil {
			return err
		}
	} else {
		// Validate required fields when not using file
		if flags.aet == "" || flags.host == "" || flags.port == 0 {
			return fmt.Errorf("when not using --file, the following flags are required: --aet, --host, --port")
		}

		// Build request from flags
		request = buildModalityRequest(flags)
	}

	// Create or update the modality
	err = client.CreateOrUpdateModality(modalityName, request)
	if err != nil {
		return fmt.Errorf("failed to create/update modality: %w", err)
	}

	fmt.Printf("Successfully created/updated modality: %s\n", modalityName)
	fmt.Printf("AET: %s\n", request.AET)
	fmt.Printf("Host: %s\n", request.Host)
	fmt.Printf("Port: %d\n", request.Port)

	return nil
}

// readModalityFromFile reads modality configuration from a JSON file
func readModalityFromFile(filePath string) (*types.ModalityCreateRequest, error) {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON
	var request types.ModalityCreateRequest
	if err := json.Unmarshal(data, &request); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate required fields
	if request.AET == "" || request.Host == "" || request.Port == 0 {
		return nil, fmt.Errorf("JSON file must contain AET, Host, and Port fields")
	}

	return &request, nil
}

// buildModalityRequest builds a modality request from flags
func buildModalityRequest(flags *CreateFlags) *types.ModalityCreateRequest {
	request := &types.ModalityCreateRequest{
		AET:          flags.aet,
		Host:         flags.host,
		Port:         flags.port,
		Manufacturer: flags.manufacturer,
		Timeout:      flags.timeout,
	}

	// Set boolean pointers for permissions
	// Note: Due to omitempty in the gorthanc library, when these are false,
	// they will be omitted from the JSON. This is a limitation that will be
	// fixed in the library later.
	request.AllowEcho = &flags.allowEcho
	request.AllowFind = &flags.allowFind
	request.AllowGet = &flags.allowGet
	request.AllowMove = &flags.allowMove
	request.AllowStore = &flags.allowStore

	return request
}
