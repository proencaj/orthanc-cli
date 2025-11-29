package instances

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

// NewListCommand creates the instances list command
func NewListCommand() *cobra.Command {
	flags := &ListFlags{}

	command := &cobra.Command{
		Use:   "list",
		Short: "List instances in the Orthanc server",
		Long:  `Retrieve and display a list of all instances stored in the Orthanc server.`,
		Example: `  # List all instances (IDs only)
  orthanc instances list

  # List first 10 instances
  orthanc instances list --limit 10

  # List instances with full details
  orthanc instances list --expand

  # Output in JSON format
  orthanc instances list --json
  orthanc instances list --expand --json`,
		RunE: func(c *cobra.Command, args []string) error {
			return runList(flags)
		},
	}

	// Add flags
	command.Flags().IntVar(&flags.limit, "limit", 100, "Maximum number of instances to return")
	command.Flags().IntVar(&flags.since, "since", 0, "Start from this index")
	command.Flags().BoolVar(&flags.expand, "expand", false, "Show full instance details")
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
	params := &types.InstancesQueryParams{
		Since: flags.since,
		Limit: flags.limit,
	}

	// Fetch instance IDs
	instanceIDs, err := client.GetAllInstances(params)
	if err != nil {
		return fmt.Errorf("failed to fetch instances: %w", err)
	}

	// If expand is requested, fetch details for each instance
	if flags.expand {
		instances := make([]types.Instance, 0, len(instanceIDs))
		for _, id := range instanceIDs {
			instance, err := client.GetInstanceDetails(id)
			if err != nil {
				return fmt.Errorf("failed to fetch instance details for %s: %w", id, err)
			}
			instances = append(instances, *instance)
		}
		return displayInstancesExpanded(instances, jsonOutput)
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

	// Raw text output - one instance per line with key info
	for _, i := range instances {
		fmt.Printf("OrthancInstanceID: %s\n", i.ID)
		fmt.Printf("SOPInstanceUID: %s\n", i.MainDicomTags.SOPInstanceUID)
		fmt.Printf("InstanceNumber: %s\n", i.MainDicomTags.InstanceNumber)
		fmt.Printf("ParentSeries: %s\n", i.ParentSeries)
		fmt.Printf("FileSize: %d bytes\n", i.FileSize)
		if i.IndexInSeries > 0 {
			fmt.Printf("IndexInSeries: %d\n", i.IndexInSeries)
		}
		fmt.Printf("\n")
	}

	return nil
}
