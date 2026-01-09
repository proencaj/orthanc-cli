package servers

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

// NewRemoveCommand creates the servers remove command
func NewRemoveCommand() *cobra.Command {
	flags := &RemoveFlags{}

	command := &cobra.Command{
		Use:   "remove <server-name>",
		Short: "Remove a DICOMweb server configuration",
		Long:  `Delete a DICOMweb server configuration from the Orthanc server. This operation is irreversible.`,
		Example: `  # Remove a server with confirmation prompt
  orthanc servers remove my-pacs

  # Remove a server without confirmation
  orthanc servers remove my-pacs --force`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runRemove(args[0], flags)
		},
	}

	// Add flags
	command.Flags().BoolVarP(&flags.force, "force", "f", false, "Skip confirmation prompt")

	return command
}

func runRemove(serverName string, flags *RemoveFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// If not using force flag, prompt for confirmation
	if !flags.force {
		confirmed, err := confirmRemoval(serverName)
		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}
		if !confirmed {
			fmt.Println("Operation cancelled")
			return nil
		}
	}

	// Delete the server
	err = client.DeleteDicomWebServer(serverName)
	if err != nil {
		return fmt.Errorf("failed to delete server: %w", err)
	}

	fmt.Printf("Successfully deleted DICOMweb server: %s\n", serverName)
	return nil
}

func confirmRemoval(serverName string) (bool, error) {
	fmt.Printf("\nWARNING: You are about to delete DICOMweb server '%s'\n", serverName)
	fmt.Println("This operation is NOT reversible and will permanently remove the server configuration.")
	fmt.Print("\nDo you really want to delete this server? (yes/no): ")

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
