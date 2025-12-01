package studies

import (
	"encoding/json"
	"fmt"

	"github.com/proencaj/gorthanc/types"
	"github.com/spf13/cobra"
)

// ListInstancesFlags holds the flags for the list-instances command
type ListInstancesFlags struct {
	expand     bool
	jsonOutput bool
}

// NewListInstancesCommand creates the studies list-instances command
func NewListInstancesCommand() *cobra.Command {
	flags := &ListInstancesFlags{}

	command := &cobra.Command{
		Use:   "list-instances <study-id>",
		Short: "List instances in a study",
		Long:  `Retrieve and display a list of all instances belonging to a specific study.`,
		Example: `  # List all instances in a study (IDs only)
  orthanc studies list-instances abc123

  # List instances with detailed information
  orthanc studies list-instances abc123 --expand

  # List instances in JSON format
  orthanc studies list-instances abc123 --json

  # List instances with details in JSON format
  orthanc studies list-instances abc123 --expand --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runListInstances(args[0], flags)
		},
	}

	// Add flags
	command.Flags().BoolVar(&flags.expand, "expand", false, "Show detailed information for each instance")
	command.Flags().BoolVar(&flags.jsonOutput, "json", false, "Output in JSON format")

	return command
}

func runListInstances(studyID string, flags *ListInstancesFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Check if JSON output should be used (flag or config)
	jsonOutput := flags.jsonOutput || shouldUseJSON()

	if flags.expand {
		// Fetch expanded instances information
		instances, err := client.GetStudyInstancesExpanded(studyID)
		if err != nil {
			return fmt.Errorf("failed to fetch instances: %w", err)
		}
		return displayInstancesExpanded(instances, jsonOutput)
	}

	// Fetch instance IDs only
	instanceIDs, err := client.GetStudyInstances(studyID)
	if err != nil {
		return fmt.Errorf("failed to fetch instances: %w", err)
	}
	return displayInstanceIDs(instanceIDs, jsonOutput)
}

func displayInstanceIDs(instanceIDs []string, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(instanceIDs, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output - one ID per line
	for _, id := range instanceIDs {
		fmt.Println(id)
	}

	return nil
}

func displayInstancesExpanded(instances []types.Instance, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(instances, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output - formatted information
	if len(instances) == 0 {
		fmt.Println("No instances found.")
		return nil
	}

	fmt.Printf("Found %d instances:\n\n", len(instances))

	for i, inst := range instances {
		fmt.Printf("Instance %d:\n", i+1)
		fmt.Printf("  ID: %s\n", inst.ID)
		if inst.MainDicomTags.SOPInstanceUID != "" {
			fmt.Printf("  SOP Instance UID: %s\n", inst.MainDicomTags.SOPInstanceUID)
		}
		if inst.MainDicomTags.InstanceNumber != "" {
			fmt.Printf("  Instance Number: %s\n", inst.MainDicomTags.InstanceNumber)
		}
		if inst.IndexInSeries > 0 {
			fmt.Printf("  Index in Series: %d\n", inst.IndexInSeries)
		}
		if inst.MainDicomTags.ImagePositionPatient != "" {
			fmt.Printf("  Image Position: %s\n", inst.MainDicomTags.ImagePositionPatient)
		}
		if inst.FileSize > 0 {
			fmt.Printf("  File Size: %.2f KB\n", float64(inst.FileSize)/1024)
		}
		if inst.FileUuid != "" {
			fmt.Printf("  File UUID: %s\n", inst.FileUuid)
		}
		if inst.ParentSeries != "" {
			fmt.Printf("  Parent Series: %s\n", inst.ParentSeries)
		}
		fmt.Println()
	}

	return nil
}
