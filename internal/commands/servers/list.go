package servers

import (
	"encoding/json"
	"fmt"

	"github.com/proencaj/gorthanc/types"
	"github.com/spf13/cobra"
)

// ListFlags holds the flags for the list command
type ListFlags struct {
	expand     bool
	jsonOutput bool
}

// NewListCommand creates the servers list command
func NewListCommand() *cobra.Command {
	flags := &ListFlags{}

	command := &cobra.Command{
		Use:   "list",
		Short: "List DICOMweb servers in the Orthanc server",
		Long:  `Retrieve and display a list of all configured DICOMweb servers in the Orthanc server.`,
		Example: `  # List all server names
  orthanc servers list

  # List servers with full details
  orthanc servers list --expand

  # Output in JSON format
  orthanc servers list --json
  orthanc servers list --expand --json`,
		RunE: func(c *cobra.Command, args []string) error {
			return runList(flags)
		},
	}

	// Add flags
	command.Flags().BoolVar(&flags.expand, "expand", false, "Show full server details")
	command.Flags().BoolVar(&flags.jsonOutput, "json", false, "Output in JSON format")

	return command
}

func runList(flags *ListFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Check if JSON output should be used (flag or config)
	jsonOutput := flags.jsonOutput || shouldUseJSON()

	// If expand is requested, fetch details for all servers
	if flags.expand {
		servers, err := client.GetDicomWebServersExpanded()
		if err != nil {
			return fmt.Errorf("failed to fetch servers: %w", err)
		}
		return displayServersExpanded(servers, jsonOutput)
	}

	// Fetch server names only
	serverNames, err := client.GetDicomWebServers()
	if err != nil {
		return fmt.Errorf("failed to fetch servers: %w", err)
	}

	return displayServerNames(serverNames, jsonOutput)
}

func displayServerNames(serverNames []string, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(serverNames, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output - one name per line
	if len(serverNames) == 0 {
		fmt.Println("No DICOMweb servers configured")
		return nil
	}

	for _, name := range serverNames {
		fmt.Println(name)
	}

	return nil
}

func displayServersExpanded(servers map[string]types.DicomWebServer, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(servers, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output - one server per block with details
	if len(servers) == 0 {
		fmt.Println("No DICOMweb servers configured")
		return nil
	}

	for name, server := range servers {
		fmt.Printf("Server: %s\n", name)
		fmt.Printf("  URL: %s\n", server.Url)
		if server.Username != "" {
			fmt.Printf("  Username: %s\n", server.Username)
		}
		if server.HasDelete != "" {
			fmt.Printf("  HasDelete: %s\n", server.HasDelete)
		}
		if server.ChunkedTransfers != "" {
			fmt.Printf("  ChunkedTransfers: %s\n", server.ChunkedTransfers)
		}
		if server.HasWadoRsUniversalTransferSyntax != "" {
			fmt.Printf("  HasWadoRsUniversalTransferSyntax: %s\n", server.HasWadoRsUniversalTransferSyntax)
		}
		fmt.Println()
	}

	return nil
}
