package studies

import (
	"encoding/json"
	"fmt"

	"github.com/proencaj/gorthanc/types"
	"github.com/spf13/cobra"
)

// GetFlags holds the flags for the get command
type GetFlags struct {
	jsonOutput bool
}

// NewGetCommand creates the studies get command
func NewGetCommand() *cobra.Command {
	flags := &GetFlags{}

	command := &cobra.Command{
		Use:   "get <study-id>",
		Short: "Get detailed information about a specific study",
		Long:  `Retrieve and display detailed information about a study using its Orthanc Study ID.`,
		Example: `  # Get study information
  orthanc studies get abc123def456ghi789

  # Get study information in JSON format
  orthanc studies get abc123def456ghi789 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			studyID := args[0]
			return runGet(studyID, flags)
		},
	}

	// Add flags
	command.Flags().BoolVar(&flags.jsonOutput, "json", false, "Output in JSON format")

	return command
}

func runGet(studyID string, flags *GetFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Check if JSON output should be used (flag or config)
	jsonOutput := flags.jsonOutput || shouldUseJSON()

	// Fetch study details
	study, err := client.GetStudy(studyID)
	if err != nil {
		return fmt.Errorf("failed to fetch study: %w", err)
	}

	return displayStudy(study, jsonOutput)
}

func displayStudy(study *types.Study, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(study, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	fmt.Printf("OrthancStudyID: %s\n", study.ID)
	fmt.Printf("AccessionNumber: %s\n", study.MainDicomTags.AccessionNumber)
	fmt.Printf("StudyInstanceUID: %s\n", study.MainDicomTags.StudyInstanceUID)
	fmt.Printf("StudyDate: %s\n", study.MainDicomTags.StudyDate)
	fmt.Printf("StudyTime: %s\n", study.MainDicomTags.StudyTime)
	fmt.Printf("StudyDescription: %s\n", study.MainDicomTags.StudyDescription)
	fmt.Printf("Series: \n")
	for _, seriesId := range study.Series {
		fmt.Printf("  %s\n", seriesId)
	}
	fmt.Printf("PatientId: %s\n", study.PatientMainDicomTags.PatientID)
	fmt.Printf("PatientName: %s\n", study.PatientMainDicomTags.PatientName)
	fmt.Printf("PatientBirthDate: %s\n", study.PatientMainDicomTags.PatientBirthDate)
	fmt.Printf("PatientSex: %s\n", study.PatientMainDicomTags.PatientSex)
	fmt.Printf("IsStable: %v\n", study.IsStable)
	fmt.Printf("LastUpdate: %v\n", study.LastUpdate)
	fmt.Printf("\n")

	return nil
}
