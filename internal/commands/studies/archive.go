package studies

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// ArchiveFlags holds the flags for the archive command
type ArchiveFlags struct {
	output string
}

// NewArchiveCommand creates the studies archive command
func NewArchiveCommand() *cobra.Command {
	flags := &ArchiveFlags{}

	command := &cobra.Command{
		Use:   "archive <study-id>",
		Short: "Download and archive a study from the Orthanc server",
		Long:  `Download a study as a ZIP archive from the Orthanc server and save it to disk.`,
		Example: `  # Archive a study to current directory
  orthanc studies archive abc123

  # Archive a study to a specific file
  orthanc studies archive abc123 --output /path/to/study.zip

  # Archive a study to a specific directory
  orthanc studies archive abc123 --output /path/to/directory/`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runArchive(args[0], flags)
		},
	}

	// Add flags
	command.Flags().StringVarP(&flags.output, "output", "o", "", "Output path (file or directory, defaults to current directory)")

	return command
}

func runArchive(studyID string, flags *ArchiveFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Download the study archive
	fmt.Printf("Downloading study archive: %s\n", studyID)
	resp, err := client.DownloadStudyArchive(studyID)
	if err != nil {
		return fmt.Errorf("failed to download study archive: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to download study archive: HTTP %d", resp.StatusCode)
	}

	// Determine the output path
	outputPath, err := determineOutputPath(flags.output, studyID)
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
		return fmt.Errorf("failed to write archive to file: %w", err)
	}

	fmt.Printf("Successfully downloaded study archive to: %s\n", outputPath)
	fmt.Printf("Size: %.2f MB\n", float64(written)/(1024*1024))

	return nil
}

// determineOutputPath determines the final output path for the archive
func determineOutputPath(output string, studyID string) (string, error) {
	// If no output specified, use current directory
	if output == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %w", err)
		}
		return filepath.Join(cwd, fmt.Sprintf("%s.zip", studyID)), nil
	}

	// Check if the output is a directory
	info, err := os.Stat(output)
	if err == nil && info.IsDir() {
		// It's an existing directory
		return filepath.Join(output, fmt.Sprintf("%s.zip", studyID)), nil
	}

	// Check if the output ends with a path separator (indicating it should be a directory)
	if output[len(output)-1] == os.PathSeparator || output[len(output)-1] == '/' {
		// Ensure the directory exists
		if err := os.MkdirAll(output, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
		return filepath.Join(output, fmt.Sprintf("%s.zip", studyID)), nil
	}

	// It's a file path - ensure the parent directory exists
	dir := filepath.Dir(output)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create parent directory: %w", err)
	}

	return output, nil
}
