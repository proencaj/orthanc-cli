package modalities

import (
	"encoding/json"
	"fmt"

	"github.com/proencaj/gorthanc/types"
	"github.com/spf13/cobra"
)

// StoreFlags holds the flags for the store command
type StoreFlags struct {
	resources         []string
	synchronous       bool
	localAet          string
	remoteAet         string
	timeout           int
	moveOriginatorAet string
	moveOriginatorID  int
	permissive        int
	storageCommitment int
	jsonOutput        bool
}

// NewStoreCommand creates the modalities store command
func NewStoreCommand() *cobra.Command {
	flags := &StoreFlags{}

	command := &cobra.Command{
		Use:   "store <modality-name> <resource-id> [resource-id...]",
		Short: "Perform a C-STORE operation to send resources to a DICOM modality",
		Long: `Execute a DICOM C-STORE operation to send studies, series, or instances to a remote modality.
The resources (identified by their Orthanc IDs) will be transmitted to the specified DICOM modality.`,
		Example: `  # Store a single study to a PACS
  orthanc modalities store PACS_SERVER a1b2c3d4-e5f6-7890-abcd-ef1234567890

  # Store multiple resources
  orthanc modalities store PACS_SERVER \
    study-id-1 \
    study-id-2 \
    series-id-1

  # Store with synchronous mode (wait for completion)
  orthanc modalities store PACS_SERVER study-id \
    --synchronous

  # Store with custom timeout and local AET
  orthanc modalities store PACS_SERVER study-id \
    --timeout 60 \
    --local-aet MY_LOCAL_AET

  # Store with all options
  orthanc modalities store PACS_SERVER study-id \
    --synchronous \
    --timeout 120 \
    --local-aet ORTHANC \
    --remote-aet PACS \
    --json`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			modalityName := args[0]
			flags.resources = args[1:]
			return runStore(modalityName, flags)
		},
	}

	// Add flags
	command.Flags().BoolVar(&flags.synchronous, "synchronous", false, "Wait synchronously for the transfer to complete")
	command.Flags().StringVar(&flags.localAet, "local-aet", "", "Local AET to use for the transfer")
	command.Flags().StringVar(&flags.remoteAet, "remote-aet", "", "Remote AET (if different from modality's configured AET)")
	command.Flags().IntVar(&flags.timeout, "timeout", 30, "Timeout in seconds")
	command.Flags().StringVar(&flags.moveOriginatorAet, "move-originator-aet", "", "Move Originator AET (for C-MOVE operations)")
	command.Flags().IntVar(&flags.moveOriginatorID, "move-originator-id", 0, "Move Originator ID (for C-MOVE operations)")
	command.Flags().IntVar(&flags.permissive, "permissive", 0, "Permissive mode (0=strict, 1=permissive)")
	command.Flags().IntVar(&flags.storageCommitment, "storage-commitment", 0, "Storage commitment (0=disabled, 1=enabled)")
	command.Flags().BoolVar(&flags.jsonOutput, "json", false, "Output in JSON format")

	return command
}

func runStore(modalityName string, flags *StoreFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Check if JSON output should be used (flag or config)
	jsonOutput := flags.jsonOutput || shouldUseJSON()

	// Build the store request
	request := &types.ModalityStoreRequest{
		Resources:         flags.resources,
		Synchronous:       flags.synchronous,
		LocalAet:          flags.localAet,
		RemoteAet:         flags.remoteAet,
		Timeout:           flags.timeout,
		MoveOriginatorAet: flags.moveOriginatorAet,
		MoveOriginatorID:  flags.moveOriginatorID,
		Permissive:        flags.permissive,
		StorageCommitment: flags.storageCommitment,
	}

	// Display operation details
	if !jsonOutput {
		fmt.Printf("Performing C-STORE operation to modality: %s\n", modalityName)
		fmt.Printf("Resources to send: %d\n", len(flags.resources))
		fmt.Printf("Synchronous: %v\n", flags.synchronous)
		if flags.localAet != "" {
			fmt.Printf("Local AET: %s\n", flags.localAet)
		}
		if flags.remoteAet != "" {
			fmt.Printf("Remote AET: %s\n", flags.remoteAet)
		}
		fmt.Printf("Timeout: %d seconds\n", flags.timeout)
		fmt.Println()
		fmt.Println("Resource IDs:")
		for i, resourceID := range flags.resources {
			fmt.Printf("  %d. %s\n", i+1, resourceID)
		}
		fmt.Println()
	}

	// Perform C-STORE
	result, err := client.StoreToModalityWithOptions(modalityName, request)
	if err != nil {
		return fmt.Errorf("C-STORE failed: %w", err)
	}

	// Display results
	return displayStoreResult(result, jsonOutput)
}

func displayStoreResult(result *types.ModalityStoreResult, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output
	fmt.Println("C-STORE operation completed successfully!")
	fmt.Println()
	fmt.Printf("Description: %s\n", result.Description)
	fmt.Printf("Local AET:   %s\n", result.LocalAet)
	fmt.Printf("Remote AET:  %s\n", result.RemoteAet)

	if len(result.ParentResources) > 0 {
		fmt.Println()
		fmt.Println("Parent resources sent:")
		for i, resource := range result.ParentResources {
			fmt.Printf("  %d. %s\n", i+1, resource)
		}
	}

	fmt.Println()
	return nil
}
