package patients

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

// NewRemoveCommand creates the patients remove command
func NewRemoveCommand() *cobra.Command {
	flags := &RemoveFlags{}

	command := &cobra.Command{
		Use:   "remove <patient-id>",
		Short: "Remove a patient from the Orthanc server",
		Long:  `Delete a patient from the Orthanc server. This operation is irreversible.`,
		Example: `  # Remove a patient with confirmation prompt
  orthanc patients remove abc123

  # Remove a patient without confirmation
  orthanc patients remove abc123 --force`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runRemove(args[0], flags)
		},
	}

	// Add flags
	command.Flags().BoolVarP(&flags.force, "force", "f", false, "Skip confirmation prompt")

	return command
}

func runRemove(patientID string, flags *RemoveFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// If not using force flag, prompt for confirmation
	if !flags.force {
		confirmed, err := confirmRemoval(patientID)
		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}
		if !confirmed {
			fmt.Println("Operation cancelled")
			return nil
		}
	}

	// Delete the patient
	err = client.DeletePatient(patientID)
	if err != nil {
		return fmt.Errorf("failed to delete patient: %w", err)
	}

	fmt.Printf("Successfully deleted patient: %s\n", patientID)
	return nil
}

func confirmRemoval(patientID string) (bool, error) {
	fmt.Printf("\n⚠️  WARNING: You are about to delete patient '%s'\n", patientID)
	fmt.Println("This operation is NOT reversible and will permanently remove all associated data.")
	fmt.Print("\nDo you really want to delete this patient? (yes/no): ")

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
