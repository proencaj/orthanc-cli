package modalities

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// ListFlags holds the flags for the list command
type ListFlags struct {
	expand     bool
	jsonOutput bool
}

// NewListCommand creates the modalities list command
func NewListCommand() *cobra.Command {
	flags := &ListFlags{}

	command := &cobra.Command{
		Use:   "list",
		Short: "List DICOM modalities in the Orthanc server",
		Long:  `Retrieve and display a list of all configured DICOM modalities in the Orthanc server.`,
		Example: `  # List all modality names
  orthanc modalities list

  # List modalities with full details
  orthanc modalities list --expand

  # Output in JSON format
  orthanc modalities list --json
  orthanc modalities list --expand --json`,
		RunE: func(c *cobra.Command, args []string) error {
			return runList(flags)
		},
	}

	// Add flags
	command.Flags().BoolVar(&flags.expand, "expand", false, "Show full modality details")
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

	// Fetch modality names
	modalityNames, err := client.GetModalities()
	if err != nil {
		return fmt.Errorf("failed to fetch modalities: %w", err)
	}

	// If expand is requested, fetch details for each modality
	if flags.expand {
		modalities := make(map[string]interface{})
		for _, name := range modalityNames {
			modality, err := client.GetModalityDetails(name)
			if err != nil {
				return fmt.Errorf("failed to fetch modality details for %s: %w", name, err)
			}
			modalities[name] = modality
		}
		return displayModalitiesExpanded(modalities, jsonOutput)
	}

	return displayModalityNames(modalityNames, jsonOutput)
}

func displayModalityNames(modalityNames []string, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(modalityNames, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output - one name per line
	if len(modalityNames) == 0 {
		fmt.Println("No modalities configured")
		return nil
	}

	for _, name := range modalityNames {
		fmt.Println(name)
	}

	return nil
}

func displayModalitiesExpanded(modalities map[string]interface{}, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(modalities, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output - one modality per block with details
	if len(modalities) == 0 {
		fmt.Println("No modalities configured")
		return nil
	}

	// We need to convert the interface{} back to types.Modality for display
	// For text output, we'll use a simple format
	for name, details := range modalities {
		fmt.Printf("Modality: %s\n", name)
		// Marshal and unmarshal to display nicely
		data, _ := json.MarshalIndent(details, "  ", "  ")
		fmt.Printf("  %s\n\n", string(data))
	}

	return nil
}
