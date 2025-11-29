package series

import (
	"encoding/json"
	"fmt"

	"github.com/proencaj/gorthanc/types"
	"github.com/spf13/cobra"
)

// ListFlags holds the flags for the list command
type ListFlags struct {
	limit      int
	since      int
	expand     bool
	jsonOutput bool
}

// NewListCommand creates the series list command
func NewListCommand() *cobra.Command {
	flags := &ListFlags{}

	command := &cobra.Command{
		Use:   "list",
		Short: "List series in the Orthanc server",
		Long:  `Retrieve and display a list of all series stored in the Orthanc server.`,
		Example: `  # List all series (IDs only)
  orthanc series list

  # List first 10 series
  orthanc series list --limit 10

  # List series with full details
  orthanc series list --expand

  # Output in JSON format
  orthanc series list --json
  orthanc series list --expand --json`,
		RunE: func(c *cobra.Command, args []string) error {
			return runList(flags)
		},
	}

	// Add flags
	command.Flags().IntVar(&flags.limit, "limit", 100, "Maximum number of series to return")
	command.Flags().IntVar(&flags.since, "since", 0, "Start from this index")
	command.Flags().BoolVar(&flags.expand, "expand", false, "Show full series details")
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

	// Prepare query parameters
	params := &types.SeriesQueryParams{
		Expand: flags.expand,
		Since:  flags.since,
		Limit:  flags.limit,
	}

	// Fetch series
	if flags.expand {
		series, err := client.GetSeriesExpanded(params)
		if err != nil {
			return fmt.Errorf("failed to fetch series: %w", err)
		}

		return displaySeriesExpanded(series, jsonOutput)
	}

	// Fetch series IDs only
	seriesIDs, err := client.GetSeries(params)
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

	// Raw text output - one series per line with key info
	for _, s := range series {
		fmt.Printf("OrthancSeriesID: %s\n", s.ID)
		fmt.Printf("SeriesDescription: %s\n", s.MainDicomTags.SeriesDescription)
		fmt.Printf("SeriesInstanceUID: %s\n", s.MainDicomTags.SeriesInstanceUID)
		fmt.Printf("Modality: %s\n", s.MainDicomTags.Modality)
		fmt.Printf("ParentStudy: %s\n", s.ParentStudy)
		fmt.Printf("Instances: %d\n", len(s.Instances))
		fmt.Printf("\n")
	}

	return nil
}
