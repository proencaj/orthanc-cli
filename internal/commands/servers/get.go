package servers

import (
	"encoding/json"
	"fmt"

	"github.com/proencaj/gorthanc/types"
	"github.com/spf13/cobra"
)

// GetFlags holds the flags for the get command
type GetFlags struct {
	jsonOutput bool
}

// NewGetCommand creates the servers get command
func NewGetCommand() *cobra.Command {
	flags := &GetFlags{}

	command := &cobra.Command{
		Use:   "get <server-name>",
		Short: "Get detailed information about a DICOMweb server",
		Long:  `Retrieve and display detailed configuration information about a specific DICOMweb server.`,
		Example: `  # Get server details
  orthanc servers get my-pacs

  # Get server details in JSON format
  orthanc servers get my-pacs --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runGet(args[0], flags)
		},
	}

	// Add flags
	command.Flags().BoolVar(&flags.jsonOutput, "json", false, "Output in JSON format")

	return command
}

func runGet(serverName string, flags *GetFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Check if JSON output should be used (flag or config)
	jsonOutput := flags.jsonOutput || shouldUseJSON()

	// Fetch all servers expanded and find the one we want
	servers, err := client.GetDicomWebServersExpanded()
	if err != nil {
		return fmt.Errorf("failed to fetch servers: %w", err)
	}

	server, ok := servers[serverName]
	if !ok {
		return fmt.Errorf("server '%s' not found", serverName)
	}

	return displayServer(serverName, &server, jsonOutput)
}

func displayServer(serverName string, server *types.DicomWebServer, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(server, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output
	fmt.Printf("Server: %s\n", serverName)
	fmt.Printf("URL: %s\n", server.Url)

	if server.Username != "" {
		fmt.Printf("Username: %s\n", server.Username)
	}

	fmt.Println("\nOptions:")
	if server.HasDelete != "" {
		fmt.Printf("  HasDelete: %s\n", server.HasDelete)
	}
	if server.ChunkedTransfers != "" {
		fmt.Printf("  ChunkedTransfers: %s\n", server.ChunkedTransfers)
	}
	if server.HasWadoRsUniversalTransferSyntax != "" {
		fmt.Printf("  HasWadoRsUniversalTransferSyntax: %s\n", server.HasWadoRsUniversalTransferSyntax)
	}

	return nil
}
