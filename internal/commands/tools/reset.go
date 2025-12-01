package tools

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ResetFlags holds the flags for the reset command
type ResetFlags struct {
	force bool
}

// NewResetCommand creates the tools reset command
func NewResetCommand() *cobra.Command {
	flags := &ResetFlags{}

	command := &cobra.Command{
		Use:   "reset",
		Short: "Perform a hot restart of the Orthanc server",
		Long: `Perform a hot restart of the Orthanc server. The configuration file will be read again,
and the server will reload without stopping the process. This is useful for applying
configuration changes without shutting down the server.

WARNING: This will temporarily interrupt all ongoing operations.`,
		Example: `  # Reset the Orthanc server (with confirmation)
  orthanc tools reset

  # Force reset without confirmation
  orthanc tools reset --force`,
		Args: cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return runReset(flags)
		},
	}

	// Add flags
	command.Flags().BoolVarP(&flags.force, "force", "f", false, "Skip confirmation prompt")

	return command
}

func runReset(flags *ResetFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Confirmation prompt unless force is used
	if !flags.force {
		fmt.Println("WARNING: This will perform a hot restart of the Orthanc server.")
		fmt.Println("All ongoing operations will be temporarily interrupted.")
		fmt.Print("Are you sure you want to continue? (yes/no): ")

		var response string
		fmt.Scanln(&response)

		if response != "yes" && response != "y" {
			fmt.Println("Reset cancelled.")
			return nil
		}
	}

	fmt.Println("Resetting Orthanc server...")

	// Perform the reset
	err = client.Reset()
	if err != nil {
		return fmt.Errorf("failed to reset Orthanc server: %w", err)
	}

	fmt.Println("Orthanc server has been reset successfully.")
	fmt.Println("The configuration file has been reloaded.")

	return nil
}
