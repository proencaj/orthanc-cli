package modalities

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

// NewRemoveCommand creates the modalities remove command
func NewRemoveCommand() *cobra.Command {
	flags := &RemoveFlags{}

	command := &cobra.Command{
		Use:   "remove <modality-name>",
		Short: "Remove a DICOM modality configuration",
		Long:  `Delete a DICOM modality configuration from the Orthanc server. This operation is irreversible.`,
		Example: `  # Remove a modality with confirmation prompt
  orthanc modalities remove PACS_SERVER

  # Remove a modality without confirmation
  orthanc modalities remove PACS_SERVER --force`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runRemove(args[0], flags)
		},
	}

	// Add flags
	command.Flags().BoolVarP(&flags.force, "force", "f", false, "Skip confirmation prompt")

	return command
}

func runRemove(modalityName string, flags *RemoveFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// If not using force flag, prompt for confirmation
	if !flags.force {
		confirmed, err := confirmRemoval(modalityName)
		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}
		if !confirmed {
			fmt.Println("Operation cancelled")
			return nil
		}
	}

	// Delete the modality
	err = client.DeleteModality(modalityName)
	if err != nil {
		return fmt.Errorf("failed to delete modality: %w", err)
	}

	fmt.Printf("Successfully deleted modality: %s\n", modalityName)
	return nil
}

func confirmRemoval(modalityName string) (bool, error) {
	fmt.Printf("\n⚠️  WARNING: You are about to delete modality '%s'\n", modalityName)
	fmt.Println("This operation is NOT reversible and will permanently remove the modality configuration.")
	fmt.Print("\nDo you really want to delete this modality? (yes/no): ")

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
