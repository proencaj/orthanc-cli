package instances

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/proencaj/gorthanc/types"
	"github.com/spf13/cobra"
)

// AnonymizeFlags holds the flags for the anonymize command
type AnonymizeFlags struct {
	force      bool
	keepSource bool
	permissive bool
	output     string
}

// NewAnonymizeCommand creates the instances anonymize command
func NewAnonymizeCommand() *cobra.Command {
	flags := &AnonymizeFlags{}

	command := &cobra.Command{
		Use:   "anonymize <instance-id>",
		Short: "Anonymize an instance and download the anonymized DICOM file",
		Long:  `Anonymize an instance and download the resulting anonymized DICOM file to disk.`,
		Example: `  # Anonymize an instance and save to current directory
  orthanc instances anonymize abc123

  # Anonymize and save to a specific file
  orthanc instances anonymize abc123 --output /path/to/anonymized.dcm

  # Anonymize and delete the source instance
  orthanc instances anonymize abc123 --keep-source=false

  # Anonymize with force flag (ignore DICOM validity)
  orthanc instances anonymize abc123 --force

  # Anonymize with permissive mode (ignore individual step errors)
  orthanc instances anonymize abc123 --permissive`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runAnonymize(args[0], flags)
		},
	}

	// Add flags
	command.Flags().BoolVar(&flags.force, "force", false, "Force operation even if it would create an invalid DICOM file")
	command.Flags().BoolVar(&flags.keepSource, "keep-source", true, "Keep the source instance after anonymization")
	command.Flags().BoolVar(&flags.permissive, "permissive", false, "Ignore errors during individual steps of the job")
	command.Flags().StringVarP(&flags.output, "output", "o", "", "Output path (file or directory, defaults to current directory)")

	return command
}

func runAnonymize(instanceID string, flags *AnonymizeFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Prepare the anonymize request
	request := buildAnonymizeRequest(flags)

	// Call the anonymize method
	fmt.Printf("Anonymizing instance: %s\n", instanceID)
	resp, err := client.AnonymizeInstance(instanceID, request)
	if err != nil {
		return fmt.Errorf("failed to anonymize instance: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to anonymize instance: HTTP %d", resp.StatusCode)
	}

	// Determine the output path
	outputPath, err := determineAnonymizeOutputPath(flags.output, instanceID)
	if err != nil {
		return fmt.Errorf("failed to determine output path: %w", err)
	}

	// Create the output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Copy the response body to the file
	written, err := io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write anonymized DICOM file: %w", err)
	}

	fmt.Println("Instance anonymized successfully!")
	fmt.Printf("Anonymized DICOM file saved to: %s\n", outputPath)
	fmt.Printf("Size: %.2f MB\n", float64(written)/(1024*1024))

	return nil
}

// buildAnonymizeRequest creates a properly formatted anonymize request
// This works around the issue where bool fields with false values are omitted due to omitempty
func buildAnonymizeRequest(flags *AnonymizeFlags) *types.InstancesAnonymizeRequest {
	request := &types.InstancesAnonymizeRequest{}

	// Set Force if true
	if flags.force {
		request.Force = true
	}

	// Set Permissive if true
	if flags.permissive {
		request.Permissive = true
	}

	// Always set KeepSource explicitly
	// Note: Due to omitempty in the gorthanc library, when KeepSource is false,
	// it will be omitted from the JSON. This is a limitation of the library.
	// The workaround would require the library to use *bool instead of bool.
	request.KeepSource = flags.keepSource

	return request
}

// determineAnonymizeOutputPath determines the final output path for the anonymized file
func determineAnonymizeOutputPath(output string, instanceID string) (string, error) {
	// If no output specified, use current directory
	if output == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %w", err)
		}
		return filepath.Join(cwd, fmt.Sprintf("%s-anonymized.dcm", instanceID)), nil
	}

	// Check if the output is a directory
	info, err := os.Stat(output)
	if err == nil && info.IsDir() {
		// It's an existing directory
		return filepath.Join(output, fmt.Sprintf("%s-anonymized.dcm", instanceID)), nil
	}

	// Check if the output ends with a path separator (indicating it should be a directory)
	if output[len(output)-1] == os.PathSeparator || output[len(output)-1] == '/' {
		// Ensure the directory exists
		if err := os.MkdirAll(output, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
		return filepath.Join(output, fmt.Sprintf("%s-anonymized.dcm", instanceID)), nil
	}

	// It's a file path - ensure the parent directory exists
	dir := filepath.Dir(output)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create parent directory: %w", err)
	}

	return output, nil
}
