package tools

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ShutdownFlags holds the flags for the shutdown command
type ShutdownFlags struct {
	force bool
}

// NewShutdownCommand creates the tools shutdown command
func NewShutdownCommand() *cobra.Command {
	flags := &ShutdownFlags{}

	command := &cobra.Command{
		Use:   "shutdown",
		Short: "Shut down the Orthanc server",
		Long: `Shut down the Orthanc server completely. This will stop the Orthanc process.

WARNING: This is a destructive operation that will stop the server entirely.
You will need to manually restart the Orthanc process after shutdown.`,
		Example: `  # Shutdown the Orthanc server (with confirmation)
  orthanc tools shutdown

  # Force shutdown without confirmation
  orthanc tools shutdown --force`,
		Args: cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return runShutdown(flags)
		},
	}

	// Add flags
	command.Flags().BoolVarP(&flags.force, "force", "f", false, "Skip confirmation prompt")

	return command
}

func runShutdown(flags *ShutdownFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Confirmation prompt unless force is used
	if !flags.force {
		fmt.Println("WARNING: This will shut down the Orthanc server completely.")
		fmt.Println("You will need to manually restart the Orthanc process.")
		fmt.Print("Are you sure you want to continue? (yes/no): ")

		var response string
		fmt.Scanln(&response)

		if response != "yes" && response != "y" {
			fmt.Println("Shutdown cancelled.")
			return nil
		}
	}

	fmt.Println("Shutting down Orthanc server...")

	// Perform the shutdown
	err = client.Shutdown()
	if err != nil {
		return fmt.Errorf("failed to shutdown Orthanc server: %w", err)
	}

	fmt.Println("Orthanc server shutdown initiated successfully.")
	fmt.Println("The server process will terminate shortly.")

	return nil
}
