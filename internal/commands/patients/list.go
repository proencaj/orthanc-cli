package patients

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

// NewListCommand creates the patients list command
func NewListCommand() *cobra.Command {
	flags := &ListFlags{}

	command := &cobra.Command{
		Use:   "list",
		Short: "List patients in the Orthanc server",
		Long:  `Retrieve and display a list of all patients stored in the Orthanc server.`,
		Example: `  # List all patients (IDs only)
  orthanc patients list

  # List first 10 patients
  orthanc patients list --limit 10

  # List patients with full details
  orthanc patients list --expand

  # Output in JSON format
  orthanc patients list --json
  orthanc patients list --expand --json`,
		RunE: func(c *cobra.Command, args []string) error {
			return runList(flags)
		},
	}

	// Add flags
	command.Flags().IntVar(&flags.limit, "limit", 100, "Maximum number of patients to return")
	command.Flags().IntVar(&flags.since, "since", 0, "Start from this index")
	command.Flags().BoolVar(&flags.expand, "expand", false, "Show full patient details")
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
	params := &types.PatientQueryParams{
		Expand: false, // API doesn't support expand, we'll do it manually if needed
		Since:  flags.since,
		Limit:  flags.limit,
	}

	// Fetch patient IDs
	patientIDs, err := client.GetPatients(params)
	if err != nil {
		return fmt.Errorf("failed to fetch patients: %w", err)
	}

	// If expand is requested, fetch details for each patient
	if flags.expand {
		patients := make([]types.Patient, 0, len(patientIDs))
		for _, id := range patientIDs {
			patient, err := client.GetPatientDetails(id)
			if err != nil {
				return fmt.Errorf("failed to fetch patient details for %s: %w", id, err)
			}
			patients = append(patients, *patient)
		}
		return displayPatientsExpanded(patients, jsonOutput)
	}

	return displayPatientIDs(patientIDs, jsonOutput)
}

func displayPatientIDs(patientIDs []string, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(patientIDs, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output - one ID per line
	for _, id := range patientIDs {
		fmt.Println(id)
	}

	return nil
}

func displayPatientsExpanded(patients []types.Patient, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(patients, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output - one patient per line with key info
	for _, p := range patients {
		fmt.Printf("OrthancPatientID: %s\n", p.ID)
		fmt.Printf("PatientName: %s\n", p.MainDicomTags.PatientName)
		fmt.Printf("PatientID: %s\n", p.MainDicomTags.PatientID)
		fmt.Printf("PatientBirthDate: %s\n", p.MainDicomTags.PatientBirthDate)
		fmt.Printf("PatientSex: %s\n", p.MainDicomTags.PatientSex)
		fmt.Printf("Studies: %d\n", len(p.Studies))
		fmt.Printf("\n")
	}

	return nil
}
