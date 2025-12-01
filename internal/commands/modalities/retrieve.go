package modalities

import (
	"fmt"

	"github.com/proencaj/gorthanc/types"
	"github.com/proencaj/orthanc-cli/internal/helpers"
	"github.com/spf13/cobra"
)

// RetrieveFlags holds the flags for the retrieve command
type RetrieveFlags struct {
	level        string
	resources    map[string]string
	timeout      int
	permissive   bool
	asynchronous bool
	jsonOutput   bool
}

// NewRetrieveCommand creates the modalities retrieve command (C-GET)
func NewRetrieveCommand() *cobra.Command {
	flags := &RetrieveFlags{
		resources: make(map[string]string),
	}

	command := &cobra.Command{
		Use:   "retrieve <modality-name>",
		Short: "Perform a C-GET operation to retrieve resources from a DICOM modality",
		Long: `Execute a DICOM C-GET operation to retrieve studies, series, or instances from a remote modality.
Unlike C-MOVE, C-GET retrieves resources directly without requiring a target AET.
The resources will be retrieved and stored in your local Orthanc instance.`,
		Example: `  # Retrieve a study by StudyInstanceUID
  orthanc modalities retrieve PACS_SERVER \
    --level Study \
    --resource StudyInstanceUID=1.2.3.4.5

  # Retrieve a series
  orthanc modalities retrieve PACS_SERVER \
    --level Series \
    --resource SeriesInstanceUID=1.2.3.4.5.6

  # Retrieve with timeout
  orthanc modalities retrieve PACS_SERVER \
    --level Study \
    --resource StudyInstanceUID=1.2.3.4.5 \
    --timeout 60

  # Asynchronous retrieve (returns immediately)
  orthanc modalities retrieve PACS_SERVER \
    --level Study \
    --resource StudyInstanceUID=1.2.3.4.5 \
    --asynchronous

  # Retrieve with multiple resource identifiers
  orthanc modalities retrieve PACS_SERVER \
    --level Study \
    --resource PatientID=12345 \
    --resource StudyDate=20240101

  # Retrieve with permissive mode (ignore errors)
  orthanc modalities retrieve PACS_SERVER \
    --level Study \
    --resource StudyInstanceUID=1.2.3.4.5 \
    --permissive`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runRetrieve(args[0], flags)
		},
	}

	// Add flags
	command.Flags().StringVar(&flags.level, "level", "Study", "Query level (Patient, Study, Series, Instance)")
	command.Flags().StringToStringVar(&flags.resources, "resource", nil, "DICOM tag and value to identify resources (can be specified multiple times)")
	command.Flags().IntVar(&flags.timeout, "timeout", 30, "Timeout in seconds")
	command.Flags().BoolVar(&flags.permissive, "permissive", false, "Ignore errors during individual steps")
	command.Flags().BoolVar(&flags.asynchronous, "asynchronous", false, "Run the job asynchronously")
	command.Flags().BoolVar(&flags.jsonOutput, "json", false, "Output in JSON format")

	// Mark required flags
	command.MarkFlagRequired("resource")

	return command
}

func runRetrieve(modalityName string, flags *RetrieveFlags) error {
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

	// Convert resources map to []map[string]interface{}
	resourcesList := make([]map[string]interface{}, 0)
	if len(flags.resources) > 0 {
		resourceMap := make(map[string]interface{})
		for key, value := range flags.resources {
			resourceMap[key] = value
		}
		resourcesList = append(resourcesList, resourceMap)
	}

	// Build the retrieve request
	request := &types.ModalityGetRequest{
		Level:        flags.level,
		Resources:    resourcesList,
		Timeout:      flags.timeout,
		Permissive:   helpers.BoolPtr(flags.permissive),
		Asynchronous: helpers.BoolPtr(flags.asynchronous),
	}

	// Display operation details
	if !jsonOutput {
		fmt.Printf("Performing C-GET operation on modality: %s\n", modalityName)
		fmt.Printf("Level: %s\n", flags.level)
		fmt.Printf("Asynchronous: %v\n", flags.asynchronous)
		fmt.Println("Resource identifiers:")
		for key, value := range flags.resources {
			fmt.Printf("  %s: %s\n", key, value)
		}
		fmt.Println()
	}

	// Perform C-GET
	err = client.GetFromModality(modalityName, request)
	if err != nil {
		return fmt.Errorf("C-GET failed: %w", err)
	}

	// Display results
	if !jsonOutput {
		fmt.Println("C-GET operation completed successfully!")
		fmt.Println()
		fmt.Println("Resources have been retrieved and stored in your local Orthanc instance.")

		if flags.asynchronous {
			fmt.Println()
			fmt.Println("Note: The operation is running asynchronously. Check your Orthanc instance for completion status.")
		}
		fmt.Println()
	} else {
		fmt.Println(`{"status": "success", "message": "C-GET operation completed"}`)
	}

	return nil
}
