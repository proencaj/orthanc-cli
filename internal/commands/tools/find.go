package tools

import (
	"encoding/json"
	"fmt"

	"github.com/proencaj/gorthanc/types"
	"github.com/proencaj/orthanc-cli/internal/helpers"
	"github.com/spf13/cobra"
)

// FindFlags holds the flags for the tools find command
type FindFlags struct {
	level            string
	tags             map[string]string
	expand           bool
	limit            int
	since            int
	requestedTags    []string
	labels           []string
	labelsConstraint string
	jsonOutput       bool
}

// NewFindCommand creates the tools find command
func NewFindCommand() *cobra.Command {
	flags := &FindFlags{
		tags: make(map[string]string),
	}

	command := &cobra.Command{
		Use:   "find",
		Short: "Search for DICOM resources in the local Orthanc database",
		Long:  `Execute a search query to find patients, studies, series, or instances stored in the local Orthanc database.`,
		Example: `  # Find all studies (no filter)
  orthanc tools find --level Study

  # Find all studies for a patient
  orthanc tools find --level Study --tag PatientID=12345

  # Find studies by patient name with expanded details
  orthanc tools find --level Study --tag PatientName="DOE^JOHN" --expand

  # Find series by modality (Modality is at Series level in DICOM)
  orthanc tools find --level Series --tag Modality=CT

  # Find studies with CT modality (use ModalitiesInStudy for Study level)
  orthanc tools find --level Study --tag ModalitiesInStudy=CT

  # Find series within a study
  orthanc tools find --level Series --tag StudyInstanceUID=1.2.3.4.5 --limit 10

  # Find with multiple tags and labels (Orthanc 1.12.0+)
  orthanc tools find \
    --level Study \
    --tag PatientID=12345 \
    --label urgent

  # Find with specific requested tags (Orthanc 1.11.0+)
  orthanc tools find \
    --level Study \
    --tag PatientID=12345 \
    --requested-tag PatientName \
    --requested-tag StudyDescription

  # Output in JSON format
  orthanc tools find --level Study --tag PatientID=12345 --json`,
		Args: cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return runFind(flags)
		},
	}

	// Add flags
	command.Flags().StringVar(&flags.level, "level", "Study", "Query level (Patient, Study, Series, Instance)")
	command.Flags().StringToStringVar(&flags.tags, "tag", nil, "DICOM tag and value for query (can be specified multiple times)")
	command.Flags().BoolVar(&flags.expand, "expand", false, "Return expanded information about the resources")
	command.Flags().IntVar(&flags.limit, "limit", 0, "Limit the number of results (0 for no limit)")
	command.Flags().IntVar(&flags.since, "since", 0, "Return results starting from this index")
	command.Flags().StringSliceVar(&flags.requestedTags, "requested-tag", nil, "Specific DICOM tags to include in response (Orthanc 1.11.0+)")
	command.Flags().StringSliceVar(&flags.labels, "label", nil, "Filter resources by labels (Orthanc 1.12.0+)")
	command.Flags().StringVar(&flags.labelsConstraint, "labels-constraint", "", "How to apply label filters: All, Any, None (Orthanc 1.12.0+)")
	command.Flags().BoolVar(&flags.jsonOutput, "json", false, "Output in JSON format")

	return command
}

func runFind(flags *FindFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Check if JSON output should be used (flag or config)
	jsonOutput := flags.jsonOutput || shouldUseJSON()

	// Validate level
	validLevels := map[string]bool{
		"Patient":  true,
		"Study":    true,
		"Series":   true,
		"Instance": true,
	}
	if !validLevels[flags.level] {
		return fmt.Errorf("invalid level '%s', must be one of: Patient, Study, Series, Instance", flags.level)
	}

	// Build the find request
	request := &types.ToolsFindRequest{
		Level: types.ResourceLevel(flags.level),
		Query: flags.tags,
	}

	// Add optional parameters
	if flags.expand {
		request.Expand = helpers.BoolPtr(true)
	}

	if flags.limit > 0 {
		request.Limit = &flags.limit
	}

	if flags.since > 0 {
		request.Since = &flags.since
	}

	if len(flags.requestedTags) > 0 {
		request.RequestedTags = flags.requestedTags
	}

	if len(flags.labels) > 0 {
		request.Labels = flags.labels
	}

	if flags.labelsConstraint != "" {
		request.LabelsConstraint = flags.labelsConstraint
	}

	// Display query information
	if !jsonOutput {
		fmt.Printf("Searching local Orthanc database\n")
		fmt.Printf("Level: %s\n", flags.level)
		fmt.Println("Query tags:")
		for key, value := range flags.tags {
			fmt.Printf("  %s: %s\n", key, value)
		}
		if flags.expand {
			fmt.Println("Mode: Expanded")
		}
		if flags.limit > 0 {
			fmt.Printf("Limit: %d\n", flags.limit)
		}
		if flags.since > 0 {
			fmt.Printf("Starting from index: %d\n", flags.since)
		}
		if len(flags.labels) > 0 {
			fmt.Printf("Labels: %v\n", flags.labels)
			if flags.labelsConstraint != "" {
				fmt.Printf("Labels constraint: %s\n", flags.labelsConstraint)
			}
		}
		fmt.Println()
	}

	// Perform the search
	if flags.expand {
		results, err := client.FindExpanded(request)
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}
		return displayExpandedResults(results, jsonOutput)
	} else {
		results, err := client.Find(request)
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}
		return displaySimpleResults(results, jsonOutput)
	}
}

func displaySimpleResults(results []string, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Text output
	if len(results) == 0 {
		fmt.Println("No results found.")
		return nil
	}

	fmt.Printf("Found %d result(s):\n\n", len(results))
	for i, id := range results {
		fmt.Printf("%d. %s\n", i+1, id)
	}

	return nil
}

func displayExpandedResults(results []types.ToolsFindExpandedResource, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Text output
	if len(results) == 0 {
		fmt.Println("No results found.")
		return nil
	}

	fmt.Printf("Found %d result(s):\n\n", len(results))

	for i, result := range results {
		fmt.Printf("Result %d:\n", i+1)
		fmt.Printf("  ID: %s\n", result.ID)
		fmt.Printf("  Type: %s\n", result.Type)

		if result.IsStable {
			fmt.Printf("  Stable: Yes\n")
		}

		if result.LastUpdate != "" {
			fmt.Printf("  Last Update: %s\n", result.LastUpdate)
		}

		// Display main DICOM tags
		if len(result.MainDicomTags) > 0 {
			fmt.Println("  Main DICOM Tags:")
			for key, value := range result.MainDicomTags {
				fmt.Printf("    %s: %v\n", key, value)
			}
		}

		// Display patient DICOM tags if available
		if len(result.PatientMainDicomTags) > 0 {
			fmt.Println("  Patient DICOM Tags:")
			for key, value := range result.PatientMainDicomTags {
				fmt.Printf("    %s: %v\n", key, value)
			}
		}

		// Display labels if available
		if len(result.Labels) > 0 {
			fmt.Printf("  Labels: %v\n", result.Labels)
		}

		fmt.Println()
	}

	return nil
}
