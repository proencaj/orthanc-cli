package patients

import (
	"encoding/json"
	"fmt"

	"github.com/proencaj/gorthanc/types"
	"github.com/proencaj/orthanc-cli/internal/helpers"
	"github.com/spf13/cobra"
)

// AnonymizeFlags holds the flags for the anonymize command
type AnonymizeFlags struct {
	force      bool
	keepSource bool
	permissive bool
	jsonOutput bool
}

// NewAnonymizeCommand creates the patients anonymize command
func NewAnonymizeCommand() *cobra.Command {
	flags := &AnonymizeFlags{}

	command := &cobra.Command{
		Use:   "anonymize <patient-id>",
		Short: "Anonymize a patient in the Orthanc server",
		Long:  `Anonymize a patient, creating a new anonymized copy in the Orthanc server.`,
		Example: `  # Anonymize a patient (keeps source by default)
  orthanc patients anonymize abc123

  # Anonymize and delete the source patient
  orthanc patients anonymize abc123 --keep-source=false

  # Anonymize with force flag (ignore DICOM validity)
  orthanc patients anonymize abc123 --force

  # Anonymize with permissive mode (ignore individual step errors)
  orthanc patients anonymize abc123 --permissive

  # Anonymize with JSON output
  orthanc patients anonymize abc123 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runAnonymize(args[0], flags)
		},
	}

	// Add flags
	command.Flags().BoolVar(&flags.force, "force", false, "Force operation even if it would create an invalid DICOM file")
	command.Flags().BoolVar(&flags.keepSource, "keep-source", true, "Keep the source patient after anonymization")
	command.Flags().BoolVar(&flags.permissive, "permissive", false, "Ignore errors during individual steps of the job")
	command.Flags().BoolVar(&flags.jsonOutput, "json", false, "Output in JSON format")

	return command
}

func runAnonymize(patientID string, flags *AnonymizeFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Check if JSON output should be used (flag or config)
	jsonOutput := flags.jsonOutput || shouldUseJSON()

	// Prepare the anonymize request
	request := buildAnonymizeRequest(flags)

	// Call the anonymize method
	response, err := client.AnonymizePatient(patientID, request)
	if err != nil {
		return fmt.Errorf("failed to anonymize patient: %w", err)
	}

	// Display the results
	return displayAnonymizeResponse(response, jsonOutput)
}

// buildAnonymizeRequest creates a properly formatted anonymize request
func buildAnonymizeRequest(flags *AnonymizeFlags) *types.PatientAnonymizeRequest {
	request := &types.PatientAnonymizeRequest{
		Force:      helpers.BoolPtr(flags.force),
		Permissive: helpers.BoolPtr(flags.permissive),
		KeepSource: helpers.BoolPtr(flags.keepSource),
	}

	return request
}

func displayAnonymizeResponse(response *types.PatientAnonymizeResponse, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output
	fmt.Println("Patient anonymized successfully!")
	fmt.Printf("New Patient ID: %s\n", response.PatientID)
	fmt.Printf("Instances Anonymized: %d\n", response.InstancesCount)
	if response.FailedInstancesCount > 0 {
		fmt.Printf("Failed Instances: %d\n", response.FailedInstancesCount)
	}
	fmt.Printf("Path: %s\n", response.Path)
	if len(response.ParentResources) > 0 {
		fmt.Printf("Parent Resources: %v\n", response.ParentResources)
	}

	return nil
}
