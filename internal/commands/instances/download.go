package instances

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// DownloadFlags holds the flags for the download command
type DownloadFlags struct {
	output string
}

// NewDownloadCommand creates the instances download command
func NewDownloadCommand() *cobra.Command {
	flags := &DownloadFlags{}

	command := &cobra.Command{
		Use:   "download <instance-id>",
		Short: "Download a DICOM instance file from the Orthanc server",
		Long:  `Download a DICOM instance file from the Orthanc server and save it to disk.`,
		Example: `  # Download an instance to current directory
  orthanc instances download abc123

  # Download an instance to a specific file
  orthanc instances download abc123 --output /path/to/instance.dcm

  # Download an instance to a specific directory
  orthanc instances download abc123 --output /path/to/directory/`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runDownload(args[0], flags)
		},
	}

	// Add flags
	command.Flags().StringVarP(&flags.output, "output", "o", "", "Output path (file or directory, defaults to current directory)")

	return command
}

func runDownload(instanceID string, flags *DownloadFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Download the DICOM file
	fmt.Printf("Downloading DICOM instance: %s\n", instanceID)
	resp, err := client.DownloadDicomFile(instanceID)
	if err != nil {
		return fmt.Errorf("failed to download DICOM file: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to download DICOM file: HTTP %d", resp.StatusCode)
	}

	// Determine the output path
	outputPath, err := determineOutputPath(flags.output, instanceID)
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
		return fmt.Errorf("failed to write DICOM file: %w", err)
	}

	fmt.Printf("Successfully downloaded DICOM file to: %s\n", outputPath)
	fmt.Printf("Size: %.2f MB\n", float64(written)/(1024*1024))

	return nil
}

// determineOutputPath determines the final output path for the DICOM file
func determineOutputPath(output string, instanceID string) (string, error) {
	// If no output specified, use current directory
	if output == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %w", err)
		}
		return filepath.Join(cwd, fmt.Sprintf("%s.dcm", instanceID)), nil
	}

	// Check if the output is a directory
	info, err := os.Stat(output)
	if err == nil && info.IsDir() {
		// It's an existing directory
		return filepath.Join(output, fmt.Sprintf("%s.dcm", instanceID)), nil
	}

	// Check if the output ends with a path separator (indicating it should be a directory)
	if output[len(output)-1] == os.PathSeparator || output[len(output)-1] == '/' {
		// Ensure the directory exists
		if err := os.MkdirAll(output, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
		return filepath.Join(output, fmt.Sprintf("%s.dcm", instanceID)), nil
	}

	// It's a file path - ensure the parent directory exists
	dir := filepath.Dir(output)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create parent directory: %w", err)
	}

	return output, nil
}
