package modalities

import (
	"encoding/json"
	"fmt"

	"github.com/proencaj/gorthanc/types"
	"github.com/spf13/cobra"
)

// MoveFlags holds the flags for the move command
type MoveFlags struct {
	level        string
	targetAet    string
	resources    map[string]string
	timeout      int
	priority     int
	permissive   bool
	asynchronous bool
	limit        int
	jsonOutput   bool
}

// NewMoveCommand creates the modalities move command
func NewMoveCommand() *cobra.Command {
	flags := &MoveFlags{
		resources: make(map[string]string),
	}

	command := &cobra.Command{
		Use:   "move <modality-name>",
		Short: "Perform a C-MOVE operation to retrieve resources from a DICOM modality",
		Long: `Execute a DICOM C-MOVE operation to retrieve studies, series, or instances from a remote modality.
The resources will be moved to the specified target AET (usually your local Orthanc instance).`,
		Example: `  # Move a study by StudyInstanceUID to local Orthanc
  orthanc modalities move PACS_SERVER \
    --level Study \
    --target-aet ORTHANC \
    --resource StudyInstanceUID=1.2.3.4.5

  # Move a series
  orthanc modalities move PACS_SERVER \
    --level Series \
    --target-aet ORTHANC \
    --resource SeriesInstanceUID=1.2.3.4.5.6

  # Move with timeout and priority
  orthanc modalities move PACS_SERVER \
    --level Study \
    --target-aet ORTHANC \
    --resource StudyInstanceUID=1.2.3.4.5 \
    --timeout 60 \
    --priority 1

  # Asynchronous move (returns immediately)
  orthanc modalities move PACS_SERVER \
    --level Study \
    --target-aet ORTHANC \
    --resource StudyInstanceUID=1.2.3.4.5 \
    --asynchronous

  # Move with multiple resource identifiers
  orthanc modalities move PACS_SERVER \
    --level Study \
    --target-aet ORTHANC \
    --resource PatientID=12345 \
    --resource StudyDate=20240101`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runMove(args[0], flags)
		},
	}

	// Add flags
	command.Flags().StringVar(&flags.level, "level", "Study", "Query level (Patient, Study, Series, Instance)")
	command.Flags().StringVar(&flags.targetAet, "target-aet", "", "Target AET where resources should be moved (e.g., local Orthanc AET)")
	command.Flags().StringToStringVar(&flags.resources, "resource", nil, "DICOM tag and value to identify resources (can be specified multiple times)")
	command.Flags().IntVar(&flags.timeout, "timeout", 30, "Timeout in seconds")
	command.Flags().IntVar(&flags.priority, "priority", 0, "Priority level (0=medium, 1=high, 2=low)")
	command.Flags().BoolVar(&flags.permissive, "permissive", false, "Ignore errors during individual steps")
	command.Flags().BoolVar(&flags.asynchronous, "asynchronous", false, "Run the job asynchronously")
	command.Flags().IntVar(&flags.limit, "limit", 0, "Limit the number of resources (0 for no limit)")
	command.Flags().BoolVar(&flags.jsonOutput, "json", false, "Output in JSON format")

	// Mark required flags
	command.MarkFlagRequired("target-aet")
	command.MarkFlagRequired("resource")

	return command
}

func runMove(modalityName string, flags *MoveFlags) error {
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

	// Validate priority
	if flags.priority < 0 || flags.priority > 2 {
		return fmt.Errorf("invalid priority '%d', must be 0 (medium), 1 (high), or 2 (low)", flags.priority)
	}

	// Convert resources map to []map[string]interface{}
	// This is the tricky part similar to the Query in Find
	resourcesList := make([]map[string]interface{}, 0)
	if len(flags.resources) > 0 {
		resourceMap := make(map[string]interface{})
		for key, value := range flags.resources {
			resourceMap[key] = value
		}
		resourcesList = append(resourcesList, resourceMap)
	}

	// Build the move request
	request := &types.ModalityMoveRequest{
		Level:        flags.level,
		TargetAet:    flags.targetAet,
		Resources:    resourcesList,
		Timeout:      flags.timeout,
		Priority:     flags.priority,
		Permissive:   flags.permissive,
		Asynchronous: flags.asynchronous, // TODO: Cannot be async
		Limit:        flags.limit,
	}

	// Display operation details
	if !jsonOutput {
		fmt.Printf("Performing C-MOVE operation on modality: %s\n", modalityName)
		fmt.Printf("Level: %s\n", flags.level)
		fmt.Printf("Target AET: %s\n", flags.targetAet)
		fmt.Printf("Asynchronous: %v\n", flags.asynchronous)
		fmt.Println("Resource identifiers:")
		for key, value := range flags.resources {
			fmt.Printf("  %s: %s\n", key, value)
		}
		fmt.Println()
	}

	// Perform C-MOVE
	result, err := client.MoveFromModality(modalityName, request)
	if err != nil {
		return fmt.Errorf("C-MOVE failed: %w", err)
	}

	// Display results
	return displayMoveResult(result, jsonOutput)
}

func displayMoveResult(result *types.ModalityMoveResult, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output
	fmt.Println("C-MOVE operation completed successfully!")
	fmt.Println()
	fmt.Printf("Description: %s\n", result.Description)
	fmt.Printf("Local AET:   %s\n", result.LocalAet)
	fmt.Printf("Remote AET:  %s\n", result.RemoteAet)
	fmt.Printf("Target AET:  %s\n", result.TargetAet)

	if len(result.Query) > 0 {
		fmt.Println()
		fmt.Println("Query parameters:")
		for _, q := range result.Query {
			for key, value := range q {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}
	}

	fmt.Println()
	return nil
}
