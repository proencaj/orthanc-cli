package dicomweb

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/proencaj/gorthanc/types"
	"github.com/spf13/cobra"
)

// WadoRsFlags holds the flags for the wado-rs command
type WadoRsFlags struct {
	studyUID    string
	seriesUID   string
	instanceUID string
	frames      string
	metadata    bool
	rendered    bool
	accept      string
	quality     int
	viewport    string
	output      string
	outputDir   string
	jsonOutput  bool
}

// NewWadoRsCommand creates the dicomweb wado-rs command
func NewWadoRsCommand() *cobra.Command {
	flags := &WadoRsFlags{}

	command := &cobra.Command{
		Use:   "wado-rs",
		Short: "Retrieve DICOM objects using WADO-RS",
		Long: `Retrieve DICOM objects using the WADO-RS (RESTful Services) protocol.
WADO-RS provides web access to DICOM objects through RESTful web services.

For bulk data retrieval (studies, series), the response is multipart/related containing
multiple DICOM files. Use --output-dir to extract files to a directory, or --output
with a .zip extension to create a zip archive.`,
		Example: `  # Retrieve a complete study as a zip archive
  orthanc dicomweb wado-rs --study-uid 1.2.3 --output study.zip

  # Retrieve a study and extract files to a directory
  orthanc dicomweb wado-rs --study-uid 1.2.3 --output-dir ./study_files/

  # Retrieve a series as a zip archive
  orthanc dicomweb wado-rs --study-uid 1.2.3 --series-uid 1.2.3.4 --output series.zip

  # Retrieve a single instance (saves directly as .dcm)
  orthanc dicomweb wado-rs --study-uid 1.2.3 --series-uid 1.2.3.4 --instance-uid 1.2.3.4.5 --output instance.dcm

  # Retrieve study metadata as JSON
  orthanc dicomweb wado-rs --study-uid 1.2.3 --metadata

  # Retrieve specific frames
  orthanc dicomweb wado-rs --study-uid 1.2.3 --series-uid 1.2.3.4 --instance-uid 1.2.3.4.5 --frames 1,2,3 --output-dir ./frames/

  # Retrieve rendered instance as JPEG
  orthanc dicomweb wado-rs --study-uid 1.2.3 --series-uid 1.2.3.4 --instance-uid 1.2.3.4.5 --rendered --accept image/jpeg --output image.jpg`,
		Args: cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return runWadoRs(flags)
		},
	}

	// Required flags
	command.Flags().StringVar(&flags.studyUID, "study-uid", "", "Study Instance UID (required)")
	command.MarkFlagRequired("study-uid")

	// Optional hierarchy flags
	command.Flags().StringVar(&flags.seriesUID, "series-uid", "", "Series Instance UID")
	command.Flags().StringVar(&flags.instanceUID, "instance-uid", "", "SOP Instance UID")
	command.Flags().StringVar(&flags.frames, "frames", "", "Frame numbers to retrieve (comma-separated, e.g., 1,2,3)")

	// Operation mode flags
	command.Flags().BoolVar(&flags.metadata, "metadata", false, "Retrieve metadata instead of bulk data")
	command.Flags().BoolVar(&flags.rendered, "rendered", false, "Retrieve rendered image")

	// Rendering options
	command.Flags().StringVar(&flags.accept, "accept", "", "Accept header for content negotiation (e.g., image/jpeg, image/png)")
	command.Flags().IntVar(&flags.quality, "quality", 0, "Image quality (1-100, for rendered output)")
	command.Flags().StringVar(&flags.viewport, "viewport", "", "Viewport size for rendered output (format: width,height)")

	// Output flags
	command.Flags().StringVarP(&flags.output, "output", "o", "", "Output file path (.zip for archive, .dcm for single instance)")
	command.Flags().StringVar(&flags.outputDir, "output-dir", "", "Output directory for extracted DICOM files")
	command.Flags().BoolVar(&flags.jsonOutput, "json", false, "Output metadata in JSON format")

	return command
}

func runWadoRs(flags *WadoRsFlags) error {
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Handle metadata requests
	if flags.metadata {
		return runWadoRsMetadata(client, flags)
	}

	// Handle rendered requests
	if flags.rendered {
		return runWadoRsRendered(client, flags)
	}

	// Handle bulk data retrieval
	return runWadoRsBulkData(client, flags)
}

// wadoRsMetadataClient interface for metadata operations
type wadoRsMetadataClient interface {
	WadoRsRetrieveStudyMetadata(studyUID string) ([]map[string]interface{}, error)
	WadoRsRetrieveSeriesMetadata(studyUID, seriesUID string) ([]map[string]interface{}, error)
	WadoRsRetrieveInstanceMetadata(studyUID, seriesUID, instanceUID string) ([]map[string]interface{}, error)
}

func runWadoRsMetadata(client wadoRsMetadataClient, flags *WadoRsFlags) error {
	var metadata []map[string]interface{}
	var err error

	// Determine the level of retrieval
	switch {
	case flags.instanceUID != "":
		if flags.seriesUID == "" {
			return fmt.Errorf("series-uid is required when instance-uid is specified")
		}
		metadata, err = client.WadoRsRetrieveInstanceMetadata(flags.studyUID, flags.seriesUID, flags.instanceUID)
	case flags.seriesUID != "":
		metadata, err = client.WadoRsRetrieveSeriesMetadata(flags.studyUID, flags.seriesUID)
	default:
		metadata, err = client.WadoRsRetrieveStudyMetadata(flags.studyUID)
	}

	if err != nil {
		return fmt.Errorf("failed to retrieve metadata: %w", err)
	}

	// Output metadata
	jsonOutput := flags.jsonOutput || shouldUseJSON()
	if jsonOutput {
		data, err := json.MarshalIndent(metadata, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Default to JSON output for metadata since it's structured data
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

// wadoRsRenderedClient interface for rendered operations
type wadoRsRenderedClient interface {
	WadoRsRetrieveRenderedInstance(studyUID, seriesUID, instanceUID string, params *types.WadoRsRenderedParams) (*http.Response, error)
	WadoRsRetrieveRenderedFrames(studyUID, seriesUID, instanceUID, frameList string, params *types.WadoRsRenderedParams) (*http.Response, error)
}

func runWadoRsRendered(client wadoRsRenderedClient, flags *WadoRsFlags) error {
	if flags.instanceUID == "" || flags.seriesUID == "" {
		return fmt.Errorf("series-uid and instance-uid are required for rendered retrieval")
	}

	params := &types.WadoRsRenderedParams{}
	if flags.accept != "" {
		params.Accept = flags.accept
	}
	if flags.quality > 0 {
		params.Quality = flags.quality
	}
	if flags.viewport != "" {
		params.Viewport = flags.viewport
	}

	var resp *http.Response
	var err error

	if flags.frames != "" {
		resp, err = client.WadoRsRetrieveRenderedFrames(flags.studyUID, flags.seriesUID, flags.instanceUID, flags.frames, params)
	} else {
		resp, err = client.WadoRsRetrieveRenderedInstance(flags.studyUID, flags.seriesUID, flags.instanceUID, params)
	}

	if err != nil {
		return fmt.Errorf("failed to retrieve rendered content: %w", err)
	}
	defer resp.Body.Close()

	return writeSimpleOutput(resp.Body, flags.output)
}

// wadoRsBulkDataClient interface for bulk data operations
type wadoRsBulkDataClient interface {
	WadoRsRetrieveStudy(studyUID string) (*http.Response, error)
	WadoRsRetrieveSeries(studyUID, seriesUID string) (*http.Response, error)
	WadoRsRetrieveInstance(studyUID, seriesUID, instanceUID string) (*http.Response, error)
	WadoRsRetrieveFrames(studyUID, seriesUID, instanceUID, frameList string) (*http.Response, error)
}

func runWadoRsBulkData(client wadoRsBulkDataClient, flags *WadoRsFlags) error {
	var resp *http.Response
	var err error

	// Determine the level of retrieval
	switch {
	case flags.frames != "":
		if flags.instanceUID == "" || flags.seriesUID == "" {
			return fmt.Errorf("series-uid and instance-uid are required when frames is specified")
		}
		resp, err = client.WadoRsRetrieveFrames(flags.studyUID, flags.seriesUID, flags.instanceUID, flags.frames)
	case flags.instanceUID != "":
		if flags.seriesUID == "" {
			return fmt.Errorf("series-uid is required when instance-uid is specified")
		}
		resp, err = client.WadoRsRetrieveInstance(flags.studyUID, flags.seriesUID, flags.instanceUID)
	case flags.seriesUID != "":
		resp, err = client.WadoRsRetrieveSeries(flags.studyUID, flags.seriesUID)
	default:
		resp, err = client.WadoRsRetrieveStudy(flags.studyUID)
	}

	if err != nil {
		return fmt.Errorf("failed to retrieve DICOM data: %w", err)
	}
	defer resp.Body.Close()

	// Check if it's a multipart response
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/related") {
		return handleMultipartResponse(resp, contentType, flags)
	}

	// Single part response (e.g., single instance)
	return writeSimpleOutput(resp.Body, flags.output)
}

// handleMultipartResponse parses multipart/related response and saves DICOM files
func handleMultipartResponse(resp *http.Response, contentType string, flags *WadoRsFlags) error {
	// Parse the Content-Type header to get the boundary
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return fmt.Errorf("failed to parse content type: %w", err)
	}

	if !strings.HasPrefix(mediaType, "multipart/") {
		return fmt.Errorf("expected multipart response, got: %s", mediaType)
	}

	boundary := params["boundary"]
	if boundary == "" {
		return fmt.Errorf("no boundary found in multipart response")
	}

	// Determine output mode
	if flags.outputDir != "" {
		return extractToDirectory(resp.Body, boundary, flags.outputDir)
	}

	if flags.output != "" && strings.HasSuffix(strings.ToLower(flags.output), ".zip") {
		return extractToZip(resp.Body, boundary, flags.output)
	}

	if flags.output != "" {
		// Single output file specified but multipart response - extract to zip
		outputPath := flags.output
		if !strings.HasSuffix(strings.ToLower(outputPath), ".zip") {
			outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".zip"
			fmt.Fprintf(os.Stderr, "Multipart response detected, saving as: %s\n", outputPath)
		}
		return extractToZip(resp.Body, boundary, outputPath)
	}

	// No output specified - print info about parts
	return listMultipartParts(resp.Body, boundary)
}

// extractToDirectory extracts multipart DICOM files to a directory
func extractToDirectory(body io.Reader, boundary, outputDir string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	mr := multipart.NewReader(body, boundary)
	partNum := 0
	totalBytes := int64(0)

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read multipart: %w", err)
		}

		partNum++
		filename := generateFilename(part, partNum)
		outputPath := filepath.Join(outputDir, filename)

		file, err := os.Create(outputPath)
		if err != nil {
			part.Close()
			return fmt.Errorf("failed to create file %s: %w", outputPath, err)
		}

		n, err := io.Copy(file, part)
		file.Close()
		part.Close()

		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", outputPath, err)
		}

		totalBytes += n
		fmt.Fprintf(os.Stderr, "Extracted: %s (%d bytes)\n", filename, n)
	}

	fmt.Fprintf(os.Stderr, "\nExtracted %d files (%d bytes total) to %s\n", partNum, totalBytes, outputDir)
	return nil
}

// extractToZip extracts multipart DICOM files to a zip archive
func extractToZip(body io.Reader, boundary, outputPath string) error {
	zipFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	mr := multipart.NewReader(body, boundary)
	partNum := 0
	totalBytes := int64(0)

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read multipart: %w", err)
		}

		partNum++
		filename := generateFilename(part, partNum)

		writer, err := zipWriter.Create(filename)
		if err != nil {
			part.Close()
			return fmt.Errorf("failed to create zip entry %s: %w", filename, err)
		}

		n, err := io.Copy(writer, part)
		part.Close()

		if err != nil {
			return fmt.Errorf("failed to write zip entry %s: %w", filename, err)
		}

		totalBytes += n
	}

	fmt.Fprintf(os.Stderr, "Created %s with %d files (%d bytes total)\n", outputPath, partNum, totalBytes)
	return nil
}

// listMultipartParts lists the parts in a multipart response without saving
func listMultipartParts(body io.Reader, boundary string) error {
	mr := multipart.NewReader(body, boundary)
	partNum := 0

	fmt.Println("Multipart response contains:")
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read multipart: %w", err)
		}

		partNum++
		contentType := part.Header.Get("Content-Type")

		// Count bytes without storing
		n, _ := io.Copy(io.Discard, part)
		part.Close()

		fmt.Printf("  Part %d: %s (%d bytes)\n", partNum, contentType, n)
	}

	fmt.Printf("\nTotal: %d parts\n", partNum)
	fmt.Println("\nUse --output file.zip or --output-dir ./dir/ to save the files")
	return nil
}

// generateFilename generates a filename for a multipart part
func generateFilename(part *multipart.Part, partNum int) string {
	// Try to get filename from Content-Disposition
	if filename := part.FileName(); filename != "" {
		return filename
	}

	// Try to get Content-Location header (may contain SOP Instance UID)
	if loc := part.Header.Get("Content-Location"); loc != "" {
		// Extract the last part of the path as potential UID
		parts := strings.Split(strings.Trim(loc, "/"), "/")
		if len(parts) > 0 {
			uid := parts[len(parts)-1]
			if uid != "" {
				return uid + ".dcm"
			}
		}
	}

	// Default to numbered filename
	return fmt.Sprintf("instance_%04d.dcm", partNum)
}

// writeSimpleOutput writes a single response body to file or stdout
func writeSimpleOutput(body io.Reader, outputPath string) error {
	var writer io.Writer
	if outputPath != "" {
		file, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer file.Close()
		writer = file
	} else {
		writer = os.Stdout
	}

	n, err := io.Copy(writer, body)
	if err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	if outputPath != "" {
		fmt.Fprintf(os.Stderr, "Written %d bytes to %s\n", n, outputPath)
	}

	return nil
}
