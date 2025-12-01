package modalities

import (
	"encoding/json"
	"fmt"

	"github.com/proencaj/gorthanc/types"
	"github.com/proencaj/orthanc-cli/internal/helpers"
	"github.com/spf13/cobra"
)

// FindFlags holds the flags for the find command
type FindFlags struct {
	level      string
	tags       map[string]string
	normalize  bool
	timeout    int
	jsonOutput bool
}

// NewFindCommand creates the modalities find command
func NewFindCommand() *cobra.Command {
	flags := &FindFlags{
		tags: make(map[string]string),
	}

	command := &cobra.Command{
		Use:   "find <modality-name>",
		Short: "Perform a C-FIND query on a DICOM modality",
		Long:  `Execute a DICOM C-FIND query to search for studies, series, or instances on a remote modality.`,
		Example: `  # Find all studies for a patient
  orthanc modalities find PACS_SERVER --level Study --tag PatientID=12345

  # Find studies by patient name and date
  orthanc modalities find PACS_SERVER \
    --level Study \
    --tag PatientName="DOE^JOHN" \
    --tag StudyDate=20240101

  # Find series within a study
  orthanc modalities find PACS_SERVER \
    --level Series \
    --tag StudyInstanceUID=1.2.3.4.5

  # Find with multiple tags
  orthanc modalities find PACS_SERVER \
    --level Study \
    --tag PatientID=12345 \
    --tag Modality=CT \
    --tag StudyDate=20240101-20240131

  # Output in JSON format
  orthanc modalities find PACS_SERVER --level Study --tag PatientID=12345 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runFind(args[0], flags)
		},
	}

	// Add flags
	command.Flags().StringVar(&flags.level, "level", "Study", "Query level (Patient, Study, Series, Instance)")
	command.Flags().StringToStringVar(&flags.tags, "tag", nil, "DICOM tag and value for query (can be specified multiple times)")
	command.Flags().BoolVar(&flags.normalize, "normalize", false, "Normalize the query")
	command.Flags().IntVar(&flags.timeout, "timeout", 0, "Timeout in seconds (0 for default)")
	command.Flags().BoolVar(&flags.jsonOutput, "json", false, "Output in JSON format")

	// Mark required flags
	command.MarkFlagRequired("tag")

	return command
}

func runFind(modalityName string, flags *FindFlags) error {
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
	request := &types.ModalityFindRequest{
		Level:     flags.level,
		Query:     flags.tags,
		Normalize: helpers.BoolPtr(flags.normalize),
		Timeout:   flags.timeout,
	}

	// Perform C-FIND
	fmt.Printf("Performing C-FIND query on modality: %s\n", modalityName)
	fmt.Printf("Level: %s\n", flags.level)
	fmt.Println("Query tags:")
	for key, value := range flags.tags {
		fmt.Printf("  %s: %s\n", key, value)
	}
	fmt.Println()

	results, err := client.FindInModality(modalityName, request)
	if err != nil {
		return fmt.Errorf("C-FIND failed: %w", err)
	}

	// Display results
	return displayFindResults(results, jsonOutput)
}

func displayFindResults(results []map[string]interface{}, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output
	if len(results) == 0 {
		fmt.Println("No results found.")
		return nil
	}

	fmt.Printf("Found %d result(s):\n\n", len(results))

	for i, result := range results {
		// Display the Path if available
		if path, ok := result["Path"].(string); ok {
			fmt.Printf("Result %d: %s\n", i+1, path)
		} else {
			fmt.Printf("Result %d:\n", i+1)
			// Fallback: display all fields if Path is not available
			for key, value := range result {
				fmt.Printf("  %s: %v\n", key, value)
			}
		}
	}
	fmt.Println()

	return nil
}
