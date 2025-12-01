package studies

import (
	"encoding/json"
	"fmt"

	"github.com/proencaj/gorthanc/types"
	"github.com/spf13/cobra"
)

// ListSeriesFlags holds the flags for the list-series command
type ListSeriesFlags struct {
	expand     bool
	jsonOutput bool
}

// NewListSeriesCommand creates the studies list-series command
func NewListSeriesCommand() *cobra.Command {
	flags := &ListSeriesFlags{}

	command := &cobra.Command{
		Use:   "list-series <study-id>",
		Short: "List series in a study",
		Long:  `Retrieve and display a list of all series belonging to a specific study.`,
		Example: `  # List all series in a study (IDs only)
  orthanc studies list-series abc123

  # List series with detailed information
  orthanc studies list-series abc123 --expand

  # List series in JSON format
  orthanc studies list-series abc123 --json

  # List series with details in JSON format
  orthanc studies list-series abc123 --expand --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runListSeries(args[0], flags)
		},
	}

	// Add flags
	command.Flags().BoolVar(&flags.expand, "expand", false, "Show detailed information for each series")
	command.Flags().BoolVar(&flags.jsonOutput, "json", false, "Output in JSON format")

	return command
}

func runListSeries(studyID string, flags *ListSeriesFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Check if JSON output should be used (flag or config)
	jsonOutput := flags.jsonOutput || shouldUseJSON()

	if flags.expand {
		// Fetch expanded series information
		series, err := client.GetStudySeriesExpanded(studyID)
		if err != nil {
			return fmt.Errorf("failed to fetch series: %w", err)
		}
		return displaySeriesExpanded(series, jsonOutput)
	}

	// Fetch series IDs only
	seriesIDs, err := client.GetStudySeries(studyID)
	if err != nil {
		return fmt.Errorf("failed to fetch series: %w", err)
	}
	return displaySeriesIDs(seriesIDs, jsonOutput)
}

func displaySeriesIDs(seriesIDs []string, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(seriesIDs, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output - one ID per line
	for _, id := range seriesIDs {
		fmt.Println(id)
	}

	return nil
}

func displaySeriesExpanded(series []types.Series, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(series, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output - formatted information
	if len(series) == 0 {
		fmt.Println("No series found.")
		return nil
	}

	fmt.Printf("Found %d series:\n\n", len(series))

	for i, s := range series {
		fmt.Printf("Series %d:\n", i+1)
		fmt.Printf("  ID: %s\n", s.ID)
		if s.MainDicomTags.SeriesDescription != "" {
			fmt.Printf("  Description: %s\n", s.MainDicomTags.SeriesDescription)
		}
		if s.MainDicomTags.Modality != "" {
			fmt.Printf("  Modality: %s\n", s.MainDicomTags.Modality)
		}
		if s.MainDicomTags.SeriesInstanceUID != "" {
			fmt.Printf("  Series UID: %s\n", s.MainDicomTags.SeriesInstanceUID)
		}
		if s.MainDicomTags.SeriesNumber != "" {
			fmt.Printf("  Series Number: %s\n", s.MainDicomTags.SeriesNumber)
		}
		fmt.Printf("  Instances: %d\n", len(s.Instances))
		fmt.Printf("  Is Stable: %v\n", s.IsStable)
		if s.LastUpdate != "" {
			fmt.Printf("  Last Update: %s\n", s.LastUpdate)
		}
		fmt.Println()
	}

	return nil
}
