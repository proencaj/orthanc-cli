package instances

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/proencaj/gorthanc/types"
	"github.com/spf13/cobra"
)

// UploadFlags holds the flags for the upload command
type UploadFlags struct {
	jsonOutput bool
}

// NewUploadCommand creates the instances upload command
func NewUploadCommand() *cobra.Command {
	flags := &UploadFlags{}

	command := &cobra.Command{
		Use:   "upload <file-path>",
		Short: "Upload a DICOM file to the Orthanc server",
		Long:  `Upload a DICOM file from disk to the Orthanc server.`,
		Example: `  # Upload a DICOM file
  orthanc instances upload /path/to/file.dcm

  # Upload a DICOM file with JSON output
  orthanc instances upload /path/to/file.dcm --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runUpload(args[0], flags)
		},
	}

	// Add flags
	command.Flags().BoolVar(&flags.jsonOutput, "json", false, "Output in JSON format")

	return command
}

func runUpload(filePath string, flags *UploadFlags) error {
	// Get the Orthanc client
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Check if JSON output should be used (flag or config)
	jsonOutput := flags.jsonOutput || shouldUseJSON()

	// Check if file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", filePath)
		}
		return fmt.Errorf("failed to stat file: %w", err)
	}

	// Check if it's a file (not a directory)
	if fileInfo.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", filePath)
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file name for display
	fileName := filepath.Base(filePath)
	fileSize := fileInfo.Size()

	fmt.Printf("Uploading DICOM file: %s (%.2f MB)\n", fileName, float64(fileSize)/(1024*1024))

	// Upload the file
	response, err := client.UploadDicomFile(file)
	if err != nil {
		return fmt.Errorf("failed to upload DICOM file: %w", err)
	}

	// Display the results
	return displayUploadResponse(response, jsonOutput)
}

func displayUploadResponse(response *types.UploadDicomFileResponse, jsonOutput bool) error {
	if jsonOutput {
		data, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Raw text output
	fmt.Println("DICOM file uploaded successfully!")
	fmt.Printf("Instance ID: %s\n", response.ID)
	fmt.Printf("Status: %s\n", response.Status)
	fmt.Printf("Path: %s\n", response.Path)

	if response.ParentPatient != "" {
		fmt.Printf("Parent Patient: %s\n", response.ParentPatient)
	}
	if response.ParentStudy != "" {
		fmt.Printf("Parent Study: %s\n", response.ParentStudy)
	}
	if response.ParentSeries != "" {
		fmt.Printf("Parent Series: %s\n", response.ParentSeries)
	}

	return nil
}
