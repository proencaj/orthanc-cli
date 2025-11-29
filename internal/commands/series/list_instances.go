package series

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// ListInstancesFlags holds the flags for the list-instances command
type ListInstancesFlags struct {
	jsonOutput bool
}

// NewListInstancesCommand creates the series list-instances command
func NewListInstancesCommand() *cobra.Command {
	flags := &ListInstancesFlags{}

	command := &cobra.Command{
		Use:   "list-instances <series-id>",
		Short: "List instances in a series",
		Long:  `Retrieve and display a list of all instances belonging to a specific series.`,
		Example: `  # List all instances in a series
  orthanc series list-instances abc123

  # List instances in JSON format
  orthanc series list-instances abc123 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runListInstances(args[0], flags)
		},
	}

	// Add flags
	command.Flags().BoolVar(&flags.jsonOutput, "json", false, "Output in JSON format")

	return command
}

func runListInstances(seriesID string, flags *ListInstancesFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Check if JSON output should be used (flag or config)
	jsonOutput := flags.jsonOutput || shouldUseJSON()

	// Fetch instances
	instanceIDs, err := client.GetSeriesInstances(seriesID)
	if err != nil {
		return fmt.Errorf("failed to fetch instances: %w", err)
	}

	// Display the instances
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
