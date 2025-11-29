package modalities

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewEchoCommand creates the modalities echo command
func NewEchoCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "echo <modality-name>",
		Short: "Test connectivity to a DICOM modality using C-ECHO",
		Long:  `Perform a DICOM C-ECHO operation to test connectivity and verify that the modality is responding.`,
		Example: `  # Test connection to a modality
  orthanc modalities echo PACS_SERVER

  # Echo a specific modality to verify it's online
  orthanc modalities echo MY_MODALITY`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runEcho(args[0])
		},
	}

	return command
}

func runEcho(modalityName string) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Perform C-ECHO
	fmt.Printf("Performing C-ECHO to modality: %s\n", modalityName)
	err = client.EchoModality(modalityName)
	if err != nil {
		fmt.Printf("✗ C-ECHO failed: \n")
		fmt.Printf("The modality '%s' is not responding or not reachable.\n", modalityName)
		return nil
	}

	fmt.Printf("✓ C-ECHO successful! Modality '%s' is responding.\n", modalityName)
	return nil
}
