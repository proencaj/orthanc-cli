package dicomweb

import (
	"fmt"
	"io"
	"os"

	"github.com/proencaj/gorthanc/types"
	"github.com/spf13/cobra"
)

// WadoFlags holds the flags for the wado command
type WadoFlags struct {
	studyUID       string
	seriesUID      string
	objectUID      string
	contentType    string
	transferSyntax string
	anonymize      bool
	frameNumber    int
	imageQuality   int
	windowCenter   string
	windowWidth    string
	rows           int
	columns        int
	region         string
	output         string
}

// NewWadoCommand creates the dicomweb wado command
func NewWadoCommand() *cobra.Command {
	flags := &WadoFlags{}

	command := &cobra.Command{
		Use:   "wado",
		Short: "Retrieve DICOM objects using WADO-URI",
		Long: `Retrieve DICOM objects using the legacy WADO-URI protocol.
WADO-URI provides web access to DICOM persistent objects through URI-based requests.`,
		Example: `  # Retrieve a DICOM instance
  orthanc dicomweb wado --study-uid 1.2.3 --series-uid 1.2.3.4 --object-uid 1.2.3.4.5 --output instance.dcm

  # Retrieve as JPEG with specific window settings
  orthanc dicomweb wado --study-uid 1.2.3 --series-uid 1.2.3.4 --object-uid 1.2.3.4.5 --content-type image/jpeg --window-center 40 --window-width 400

  # Retrieve a specific frame
  orthanc dicomweb wado --study-uid 1.2.3 --series-uid 1.2.3.4 --object-uid 1.2.3.4.5 --frame 1 --output frame.dcm`,
		Args: cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return runWado(flags)
		},
	}

	// Required flags
	command.Flags().StringVar(&flags.studyUID, "study-uid", "", "Study Instance UID (required)")
	command.Flags().StringVar(&flags.seriesUID, "series-uid", "", "Series Instance UID (required)")
	command.Flags().StringVar(&flags.objectUID, "object-uid", "", "SOP Instance UID (required)")
	command.MarkFlagRequired("study-uid")
	command.MarkFlagRequired("series-uid")
	command.MarkFlagRequired("object-uid")

	// Optional flags
	command.Flags().StringVar(&flags.contentType, "content-type", "", "Requested content type (e.g., application/dicom, image/jpeg, image/png)")
	command.Flags().StringVar(&flags.transferSyntax, "transfer-syntax", "", "Transfer syntax UID")
	command.Flags().BoolVar(&flags.anonymize, "anonymize", false, "Anonymize the returned object")
	command.Flags().IntVar(&flags.frameNumber, "frame", 0, "Frame number to retrieve (for multi-frame instances)")
	command.Flags().IntVar(&flags.imageQuality, "quality", 0, "Image quality (1-100, for lossy compression)")
	command.Flags().StringVar(&flags.windowCenter, "window-center", "", "Window center for rendering")
	command.Flags().StringVar(&flags.windowWidth, "window-width", "", "Window width for rendering")
	command.Flags().IntVar(&flags.rows, "rows", 0, "Number of rows for scaling")
	command.Flags().IntVar(&flags.columns, "columns", 0, "Number of columns for scaling")
	command.Flags().StringVar(&flags.region, "region", "", "Region of interest (format: x,y,width,height)")
	command.Flags().StringVarP(&flags.output, "output", "o", "", "Output file path (defaults to stdout)")

	return command
}

func runWado(flags *WadoFlags) error {
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Build WADO-URI parameters
	params := &types.WadoUriParams{
		RequestType: "WADO",
		StudyUID:    flags.studyUID,
		SeriesUID:   flags.seriesUID,
		ObjectUID:   flags.objectUID,
	}

	if flags.contentType != "" {
		params.ContentType = flags.contentType
	}
	if flags.transferSyntax != "" {
		params.TransferSyntax = flags.transferSyntax
	}
	if flags.anonymize {
		params.Anonymize = "yes"
	}
	if flags.frameNumber > 0 {
		params.FrameNumber = flags.frameNumber
	}
	if flags.imageQuality > 0 {
		params.ImageQuality = flags.imageQuality
	}
	if flags.windowCenter != "" {
		params.WindowCenter = flags.windowCenter
	}
	if flags.windowWidth != "" {
		params.WindowWidth = flags.windowWidth
	}
	if flags.rows > 0 {
		params.Rows = flags.rows
	}
	if flags.columns > 0 {
		params.Columns = flags.columns
	}
	if flags.region != "" {
		params.Region = flags.region
	}

	// Execute WADO-URI request
	resp, err := client.WadoUriRetrieve(params)
	if err != nil {
		return fmt.Errorf("failed to retrieve DICOM object: %w", err)
	}
	defer resp.Body.Close()

	// Write output
	var writer io.Writer
	if flags.output != "" {
		file, err := os.Create(flags.output)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer file.Close()
		writer = file
	} else {
		writer = os.Stdout
	}

	n, err := io.Copy(writer, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	if flags.output != "" {
		fmt.Fprintf(os.Stderr, "Written %d bytes to %s\n", n, flags.output)
	}

	return nil
}
