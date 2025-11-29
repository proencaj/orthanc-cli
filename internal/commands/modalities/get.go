package modalities

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// GetFlags holds the flags for the get command
type GetFlags struct {
	jsonOutput bool
}

// NewGetCommand creates the modalities get command
func NewGetCommand() *cobra.Command {
	flags := &GetFlags{}

	command := &cobra.Command{
		Use:   "get <modality-name>",
		Short: "Get detailed information about a DICOM modality",
		Long:  `Retrieve and display detailed configuration information about a specific DICOM modality.`,
		Example: `  # Get modality details
  orthanc modalities get PACS_SERVER

  # Get modality details in JSON format
  orthanc modalities get PACS_SERVER --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runGet(args[0], flags)
		},
	}

	// Add flags
	command.Flags().BoolVar(&flags.jsonOutput, "json", false, "Output in JSON format")

	return command
}

func runGet(modalityName string, flags *GetFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Check if JSON output should be used (flag or config)
	jsonOutput := flags.jsonOutput || shouldUseJSON()

	// Fetch modality details
	modality, err := client.GetModalityDetails(modalityName)
	if err != nil {
		return fmt.Errorf("failed to fetch modality details: %w", err)
	}

	// Display the modality
	if jsonOutput {
		data, err := json.MarshalIndent(modality, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output
	fmt.Printf("Modality: %s\n", modalityName)
	fmt.Printf("AET: %s\n", modality.AET)
	fmt.Printf("Host: %s\n", modality.Host)
	fmt.Printf("Port: %d\n", modality.Port)

	if modality.Manufacturer != "" {
		fmt.Printf("Manufacturer: %s\n", modality.Manufacturer)
	}

	// Display permissions
	fmt.Println("\nPermissions:")
	fmt.Printf("  Allow Echo: %v\n", modality.AllowEcho)
	fmt.Printf("  Allow Find: %v\n", modality.AllowFind)
	fmt.Printf("  Allow Get: %v\n", modality.AllowGet)
	fmt.Printf("  Allow Move: %v\n", modality.AllowMove)
	fmt.Printf("  Allow Store: %v\n", modality.AllowStore)

	if modality.Timeout > 0 {
		fmt.Printf("\nTimeout: %d seconds\n", modality.Timeout)
	}

	return nil
}
