package instances

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// GetFlags holds the flags for the get command
type GetFlags struct {
	jsonOutput bool
}

// NewGetCommand creates the instances get command
func NewGetCommand() *cobra.Command {
	flags := &GetFlags{}

	command := &cobra.Command{
		Use:   "get <instance-id>",
		Short: "Get detailed information about an instance",
		Long:  `Retrieve and display detailed information about a specific instance from the Orthanc server.`,
		Example: `  # Get instance details
  orthanc instances get abc123

  # Get instance details in JSON format
  orthanc instances get abc123 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runGet(args[0], flags)
		},
	}

	// Add flags
	command.Flags().BoolVar(&flags.jsonOutput, "json", false, "Output in JSON format")

	return command
}

func runGet(instanceID string, flags *GetFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Check if JSON output should be used (flag or config)
	jsonOutput := flags.jsonOutput || shouldUseJSON()

	// Fetch instance details
	instance, err := client.GetInstanceDetails(instanceID)
	if err != nil {
		return fmt.Errorf("failed to fetch instance details: %w", err)
	}

	// Display the instance
	if jsonOutput {
		data, err := json.MarshalIndent(instance, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output
	fmt.Printf("OrthancInstanceID: %s\n", instance.ID)
	fmt.Printf("SOPInstanceUID: %s\n", instance.MainDicomTags.SOPInstanceUID)
	fmt.Printf("InstanceNumber: %s\n", instance.MainDicomTags.InstanceNumber)
	fmt.Printf("ImageIndex: %s\n", instance.MainDicomTags.ImageIndex)
	fmt.Printf("ParentSeries: %s\n", instance.ParentSeries)
	fmt.Printf("FileSize: %d bytes (%.2f MB)\n", instance.FileSize, float64(instance.FileSize)/(1024*1024))
	fmt.Printf("FileUUID: %s\n", instance.FileUuid)
	if instance.IndexInSeries > 0 {
		fmt.Printf("IndexInSeries: %d\n", instance.IndexInSeries)
	}
	if instance.ModifiedFrom != "" {
		fmt.Printf("ModifiedFrom: %s\n", instance.ModifiedFrom)
	}

	return nil
}
