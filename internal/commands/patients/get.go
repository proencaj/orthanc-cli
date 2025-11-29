package patients

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// GetFlags holds the flags for the get command
type GetFlags struct {
	jsonOutput bool
}

// NewGetCommand creates the patients get command
func NewGetCommand() *cobra.Command {
	flags := &GetFlags{}

	command := &cobra.Command{
		Use:   "get <patient-id>",
		Short: "Get detailed information about a patient",
		Long:  `Retrieve and display detailed information about a specific patient from the Orthanc server.`,
		Example: `  # Get patient details
  orthanc patients get abc123

  # Get patient details in JSON format
  orthanc patients get abc123 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runGet(args[0], flags)
		},
	}

	// Add flags
	command.Flags().BoolVar(&flags.jsonOutput, "json", false, "Output in JSON format")

	return command
}

func runGet(patientID string, flags *GetFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Check if JSON output should be used (flag or config)
	jsonOutput := flags.jsonOutput || shouldUseJSON()

	// Fetch patient details
	patient, err := client.GetPatientDetails(patientID)
	if err != nil {
		return fmt.Errorf("failed to fetch patient details: %w", err)
	}

	// Display the patient
	if jsonOutput {
		data, err := json.MarshalIndent(patient, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output
	fmt.Printf("OrthancPatientID: %s\n", patient.ID)
	fmt.Printf("PatientName: %s\n", patient.MainDicomTags.PatientName)
	fmt.Printf("PatientID: %s\n", patient.MainDicomTags.PatientID)
	fmt.Printf("PatientBirthDate: %s\n", patient.MainDicomTags.PatientBirthDate)
	fmt.Printf("PatientSex: %s\n", patient.MainDicomTags.PatientSex)
	fmt.Printf("IsStable: %v\n", patient.IsStable)
	fmt.Printf("LastUpdate: %s\n", patient.LastUpdate)
	fmt.Printf("Studies: %d\n", len(patient.Studies))
	if len(patient.Studies) > 0 {
		fmt.Println("Study IDs:")
		for _, studyID := range patient.Studies {
			fmt.Printf("  - %s\n", studyID)
		}
	}

	return nil
}
