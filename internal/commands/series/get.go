package series

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// GetFlags holds the flags for the get command
type GetFlags struct {
	jsonOutput bool
}

// NewGetCommand creates the series get command
func NewGetCommand() *cobra.Command {
	flags := &GetFlags{}

	command := &cobra.Command{
		Use:   "get <series-id>",
		Short: "Get detailed information about a series",
		Long:  `Retrieve and display detailed information about a specific series from the Orthanc server.`,
		Example: `  # Get series details
  orthanc series get abc123

  # Get series details in JSON format
  orthanc series get abc123 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runGet(args[0], flags)
		},
	}

	// Add flags
	command.Flags().BoolVar(&flags.jsonOutput, "json", false, "Output in JSON format")

	return command
}

func runGet(seriesID string, flags *GetFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Check if JSON output should be used (flag or config)
	jsonOutput := flags.jsonOutput || shouldUseJSON()

	// Fetch series details
	series, err := client.GetSeriesDetail(seriesID)
	if err != nil {
		return fmt.Errorf("failed to fetch series details: %w", err)
	}

	// Display the series
	if jsonOutput {
		data, err := json.MarshalIndent(series, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output
	fmt.Printf("OrthancSeriesID: %s\n", series.ID)
	fmt.Printf("SeriesDescription: %s\n", series.MainDicomTags.SeriesDescription)
	fmt.Printf("SeriesInstanceUID: %s\n", series.MainDicomTags.SeriesInstanceUID)
	fmt.Printf("Modality: %s\n", series.MainDicomTags.Modality)
	fmt.Printf("SeriesNumber: %s\n", series.MainDicomTags.SeriesNumber)
	fmt.Printf("ParentStudy: %s\n", series.ParentStudy)
	fmt.Printf("IsStable: %v\n", series.IsStable)
	fmt.Printf("LastUpdate: %s\n", series.LastUpdate)
	fmt.Printf("Instances: %d\n", len(series.Instances))
	if series.ExpectedNumberOfInstances > 0 {
		fmt.Printf("ExpectedInstances: %d\n", series.ExpectedNumberOfInstances)
	}
	if series.Status != "" {
		fmt.Printf("Status: %s\n", series.Status)
	}

	return nil
}
