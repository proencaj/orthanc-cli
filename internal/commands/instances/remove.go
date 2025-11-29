package instances

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

// NewRemoveCommand creates the instances remove command
func NewRemoveCommand() *cobra.Command {
	flags := &RemoveFlags{}

	command := &cobra.Command{
		Use:   "remove <instance-id>",
		Short: "Remove an instance from the Orthanc server",
		Long:  `Delete an instance from the Orthanc server. This operation is irreversible.`,
		Example: `  # Remove an instance with confirmation prompt
  orthanc instances remove abc123

  # Remove an instance without confirmation
  orthanc instances remove abc123 --force`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runRemove(args[0], flags)
		},
	}

	// Add flags
	command.Flags().BoolVarP(&flags.force, "force", "f", false, "Skip confirmation prompt")

	return command
}

func runRemove(instanceID string, flags *RemoveFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// If not using force flag, prompt for confirmation
	if !flags.force {
		confirmed, err := confirmRemoval(instanceID)
		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}
		if !confirmed {
			fmt.Println("Operation cancelled")
			return nil
		}
	}

	// Delete the instance
	err = client.DeleteInstance(instanceID)
	if err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}

	fmt.Printf("Successfully deleted instance: %s\n", instanceID)
	return nil
}

func confirmRemoval(instanceID string) (bool, error) {
	fmt.Printf("\n⚠️  WARNING: You are about to delete instance '%s'\n", instanceID)
	fmt.Println("This operation is NOT reversible and will permanently remove all associated data.")
	fmt.Print("\nDo you really want to delete this instance? (yes/no): ")

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
