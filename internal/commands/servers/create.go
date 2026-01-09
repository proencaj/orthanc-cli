package servers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/proencaj/gorthanc"
	"github.com/proencaj/gorthanc/types"
	"github.com/spf13/cobra"
)

// CreateFlags holds the flags for the create command
type CreateFlags struct {
	file                           string
	url                            string
	username                       string
	password                       string
	hasDelete                      bool
	chunkedTransfers               bool
	hasWadoRsUniversalTransferSyntax bool
}

// NewCreateCommand creates the servers create command
func NewCreateCommand() *cobra.Command {
	flags := &CreateFlags{}

	command := &cobra.Command{
		Use:   "create <server-name>",
		Short: "Create or update a DICOMweb server configuration",
		Long:  `Create or update a DICOMweb server configuration in the Orthanc server. You can either provide a JSON file with --file or specify individual parameters.`,
		Example: `  # Create server from JSON file
  orthanc servers create my-pacs --file server.json

  # Create server with individual parameters
  orthanc servers create my-pacs --url https://pacs.example.com/dicom-web

  # Create with authentication
  orthanc servers create my-pacs \
    --url https://pacs.example.com/dicom-web \
    --username admin \
    --password secret

  # Create with all options
  orthanc servers create my-pacs \
    --url https://pacs.example.com/dicom-web \
    --username admin \
    --password secret \
    --has-delete \
    --chunked-transfers \
    --has-wado-rs-universal-transfer-syntax`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runCreate(args[0], flags)
		},
	}

	// Add flags
	command.Flags().StringVar(&flags.file, "file", "", "JSON file containing server configuration")
	command.Flags().StringVar(&flags.url, "url", "", "URL of the remote DICOMweb server")
	command.Flags().StringVar(&flags.username, "username", "", "Username for authentication (optional)")
	command.Flags().StringVar(&flags.password, "password", "", "Password for authentication (optional)")
	command.Flags().BoolVar(&flags.hasDelete, "has-delete", false, "Whether the server supports DELETE operations")
	command.Flags().BoolVar(&flags.chunkedTransfers, "chunked-transfers", true, "Whether to use chunked transfers (set to false for Orthanc <= 1.5.6)")
	command.Flags().BoolVar(&flags.hasWadoRsUniversalTransferSyntax, "has-wado-rs-universal-transfer-syntax", true, "Whether the server supports WADO-RS universal transfer syntax (set to false for Orthanc DICOMweb plugin <= 1.0)")

	return command
}

func runCreate(serverName string, flags *CreateFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	var request *types.DicomWebServerCreateRequest

	// If file is provided, read from file
	if flags.file != "" {
		request, err = readServerFromFile(flags.file)
		if err != nil {
			return err
		}
	} else {
		// Validate required fields when not using file
		if flags.url == "" {
			return fmt.Errorf("when not using --file, the --url flag is required")
		}

		// Build request from flags
		request = buildServerRequest(flags)
	}

	// Create or update the server
	err = client.CreateOrUpdateDicomWebServer(serverName, request)
	if err != nil {
		return fmt.Errorf("failed to create/update server: %w", err)
	}

	fmt.Printf("Successfully created/updated DICOMweb server: %s\n", serverName)
	fmt.Printf("URL: %s\n", request.Url)
	if request.Username != "" {
		fmt.Printf("Username: %s\n", request.Username)
	}

	return nil
}

// readServerFromFile reads server configuration from a JSON file
func readServerFromFile(filePath string) (*types.DicomWebServerCreateRequest, error) {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON
	var request types.DicomWebServerCreateRequest
	if err := json.Unmarshal(data, &request); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate required fields
	if request.Url == "" {
		return nil, fmt.Errorf("JSON file must contain Url field")
	}

	return &request, nil
}

// buildServerRequest builds a server request from flags
func buildServerRequest(flags *CreateFlags) *types.DicomWebServerCreateRequest {
	request := &types.DicomWebServerCreateRequest{
		Url:      flags.url,
		Username: flags.username,
		Password: flags.password,
	}

	// Set boolean pointers for optional fields
	request.HasDelete = gorthanc.BoolPtr(flags.hasDelete)
	request.ChunkedTransfers = gorthanc.BoolPtr(flags.chunkedTransfers)
	request.HasWadoRsUniversalTransferSyntax = gorthanc.BoolPtr(flags.hasWadoRsUniversalTransferSyntax)

	return request
}
