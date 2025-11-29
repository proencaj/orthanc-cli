package series

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// RemoveFlags holds the flags for the remove command
type RemoveFlags struct {
	force bool
}

// NewRemoveCommand creates the series remove command
func NewRemoveCommand() *cobra.Command {
	flags := &RemoveFlags{}

	command := &cobra.Command{
		Use:   "remove <series-id>",
		Short: "Remove a series from the Orthanc server",
		Long:  `Delete a series from the Orthanc server. This operation is irreversible.`,
		Example: `  # Remove a series with confirmation prompt
  orthanc series remove abc123

  # Remove a series without confirmation
  orthanc series remove abc123 --force`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runRemove(args[0], flags)
		},
	}

	// Add flags
	command.Flags().BoolVarP(&flags.force, "force", "f", false, "Skip confirmation prompt")

	return command
}

func runRemove(seriesID string, flags *RemoveFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// If not using force flag, prompt for confirmation
	if !flags.force {
		confirmed, err := confirmRemoval(seriesID)
		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}
		if !confirmed {
			fmt.Println("Operation cancelled")
			return nil
		}
	}

	// Delete the series
	err = client.DeleteSeries(seriesID)
	if err != nil {
		return fmt.Errorf("failed to delete series: %w", err)
	}

	fmt.Printf("Successfully deleted series: %s\n", seriesID)
	return nil
}

func confirmRemoval(seriesID string) (bool, error) {
	fmt.Printf("\n⚠️  WARNING: You are about to delete series '%s'\n", seriesID)
	fmt.Println("This operation is NOT reversible and will permanently remove all associated data.")
	fmt.Print("\nDo you really want to delete this series? (yes/no): ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	// Clean up the response
	response = strings.TrimSpace(strings.ToLower(response))

	// Accept "yes" or "y" as confirmation
	return response == "yes" || response == "y", nil
}
