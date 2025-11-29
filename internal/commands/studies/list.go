package studies

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

// NewListCommand creates the studies list command
func NewListCommand() *cobra.Command {
	flags := &ListFlags{}

	command := &cobra.Command{
		Use:   "list",
		Short: "List studies in the Orthanc server",
		Long:  `Retrieve and display a list of all studies stored in the Orthanc server.`,
		Example: `  # List all studies (IDs only)
  orthanc studies list

  # List first 10 studies
  orthanc studies list --limit 10

  # List studies with full details
  orthanc studies list --expand

  # Output in JSON format
  orthanc studies list --json
  orthanc studies list --expand --json`,
		RunE: func(c *cobra.Command, args []string) error {
			return runList(flags)
		},
	}

	// Add flags
	command.Flags().IntVar(&flags.limit, "limit", 100, "Maximum number of studies to return")
	command.Flags().IntVar(&flags.since, "since", 0, "Start from this index")
	command.Flags().BoolVar(&flags.expand, "expand", false, "Show full study details")
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
	params := &types.StudiesQueryParams{
		Expand: flags.expand,
		Since:  flags.since,
		Limit:  flags.limit,
	}

	// Fetch studies
	if flags.expand {
		studies, err := client.GetStudiesExpanded(params)
		if err != nil {
			return fmt.Errorf("failed to fetch studies: %w", err)
		}

		return displayStudiesExpanded(studies, jsonOutput)
	}

	// Fetch study IDs only
	studyIDs, err := client.GetStudies(params)
	if err != nil {
		return fmt.Errorf("failed to fetch studies: %w", err)
	}

	return displayStudyIDs(studyIDs, jsonOutput)
}

func displayStudyIDs(studyIDs []string, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(studyIDs, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output - one ID per line
	for _, id := range studyIDs {
		fmt.Println(id)
	}

	return nil
}

func displayStudiesExpanded(studies []types.Study, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(studies, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output - one study per line with key info
	for _, study := range studies {
		fmt.Printf("OrthancStudyID: %s\n", study.ID)
		fmt.Printf("AccessionNumber: %s\n", study.MainDicomTags.AccessionNumber)
		fmt.Printf("StudyInstanceUID: %s\n", study.MainDicomTags.StudyInstanceUID)
		fmt.Printf("StudyDate: %s\n", study.MainDicomTags.StudyDate)
		fmt.Printf("StudyDescription: %s\n", study.MainDicomTags.StudyDescription)
		fmt.Printf("PatientName: %s\n", study.MainDicomTags.PatientName)
		fmt.Printf("\n")
	}

	return nil
}
