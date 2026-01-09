package dicomweb

import (
	"encoding/json"
	"fmt"

	"github.com/proencaj/gorthanc/types"
	"github.com/spf13/cobra"
)

// QidoFlags holds the flags for the qido command
type QidoFlags struct {
	// Query level
	level string

	// Common query parameters
	limit        int
	offset       int
	includeField string
	fuzzyMatch   bool

	// Study-level filters
	studyUID      string
	patientID     string
	patientName   string
	accessionNum  string
	studyDate     string
	modalitiesIn  string

	// Series-level filters
	seriesUID    string
	modality     string
	seriesNumber string

	// Instance-level filters
	instanceUID string
	sopClassUID string

	// Output
	jsonOutput bool
}

// NewQidoCommand creates the dicomweb qido command
func NewQidoCommand() *cobra.Command {
	flags := &QidoFlags{}

	command := &cobra.Command{
		Use:   "qido",
		Short: "Query DICOM objects using QIDO-RS",
		Long: `Query for DICOM objects using the QIDO-RS (Query based on ID for DICOM Objects by RESTful Services) protocol.
QIDO-RS allows searching for studies, series, and instances based on various DICOM attributes.`,
		Example: `  # Search for all studies
  orthanc dicomweb qido --level studies

  # Search for studies by patient name
  orthanc dicomweb qido --level studies --patient-name "Smith*"

  # Search for studies by date range
  orthanc dicomweb qido --level studies --study-date 20230101-20231231

  # Search for series within a study
  orthanc dicomweb qido --level series --study-uid 1.2.3

  # Search for all CT series
  orthanc dicomweb qido --level series --modality CT

  # Search for instances within a series
  orthanc dicomweb qido --level instances --study-uid 1.2.3 --series-uid 1.2.3.4

  # Paginated results
  orthanc dicomweb qido --level studies --limit 10 --offset 20`,
		Args: cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return runQido(flags)
		},
	}

	// Query level
	command.Flags().StringVar(&flags.level, "level", "studies", "Query level: studies, series, or instances")

	// Pagination
	command.Flags().IntVar(&flags.limit, "limit", 0, "Maximum number of results to return")
	command.Flags().IntVar(&flags.offset, "offset", 0, "Number of results to skip")
	command.Flags().StringVar(&flags.includeField, "include-field", "", "Additional DICOM fields to include in results")
	command.Flags().BoolVar(&flags.fuzzyMatch, "fuzzy", false, "Enable fuzzy matching for string queries")

	// Study-level filters
	command.Flags().StringVar(&flags.studyUID, "study-uid", "", "Study Instance UID")
	command.Flags().StringVar(&flags.patientID, "patient-id", "", "Patient ID")
	command.Flags().StringVar(&flags.patientName, "patient-name", "", "Patient Name (supports wildcards with *)")
	command.Flags().StringVar(&flags.accessionNum, "accession-number", "", "Accession Number")
	command.Flags().StringVar(&flags.studyDate, "study-date", "", "Study Date (format: YYYYMMDD or range YYYYMMDD-YYYYMMDD)")
	command.Flags().StringVar(&flags.modalitiesIn, "modalities-in-study", "", "Modalities in Study")

	// Series-level filters
	command.Flags().StringVar(&flags.seriesUID, "series-uid", "", "Series Instance UID")
	command.Flags().StringVar(&flags.modality, "modality", "", "Modality")
	command.Flags().StringVar(&flags.seriesNumber, "series-number", "", "Series Number")

	// Instance-level filters
	command.Flags().StringVar(&flags.instanceUID, "instance-uid", "", "SOP Instance UID")
	command.Flags().StringVar(&flags.sopClassUID, "sop-class-uid", "", "SOP Class UID")

	// Output
	command.Flags().BoolVar(&flags.jsonOutput, "json", false, "Output in JSON format")

	return command
}

func runQido(flags *QidoFlags) error {
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	var results []map[string]interface{}

	switch flags.level {
	case "studies":
		results, err = runQidoStudies(client, flags)
	case "series":
		results, err = runQidoSeries(client, flags)
	case "instances":
		results, err = runQidoInstances(client, flags)
	default:
		return fmt.Errorf("invalid query level: %s (must be studies, series, or instances)", flags.level)
	}

	if err != nil {
		return err
	}

	// Output results
	jsonOutput := flags.jsonOutput || shouldUseJSON()
	return displayQidoResults(results, jsonOutput)
}

func runQidoStudies(client interface {
	QidoSearchStudies(params *types.QidoStudyQueryParams) ([]map[string]interface{}, error)
}, flags *QidoFlags) ([]map[string]interface{}, error) {
	params := &types.QidoStudyQueryParams{}

	// Common parameters
	if flags.limit > 0 {
		params.Limit = flags.limit
	}
	if flags.offset > 0 {
		params.Offset = flags.offset
	}
	if flags.includeField != "" {
		params.Includefield = flags.includeField
	}
	if flags.fuzzyMatch {
		params.FuzzyMatching = true
	}

	// Study-level filters
	if flags.studyUID != "" {
		params.StudyInstanceUID = flags.studyUID
	}
	if flags.patientID != "" {
		params.PatientID = flags.patientID
	}
	if flags.patientName != "" {
		params.PatientName = flags.patientName
	}
	if flags.accessionNum != "" {
		params.AccessionNumber = flags.accessionNum
	}
	if flags.studyDate != "" {
		params.StudyDate = flags.studyDate
	}
	if flags.modalitiesIn != "" {
		params.ModalitiesInStudy = flags.modalitiesIn
	}

	results, err := client.QidoSearchStudies(params)
	if err != nil {
		return nil, fmt.Errorf("failed to search studies: %w", err)
	}

	return results, nil
}

func runQidoSeries(client interface {
	QidoSearchSeries(studyUID string, params *types.QidoSeriesQueryParams) ([]map[string]interface{}, error)
	QidoSearchAllSeries(params *types.QidoSeriesQueryParams) ([]map[string]interface{}, error)
}, flags *QidoFlags) ([]map[string]interface{}, error) {
	params := &types.QidoSeriesQueryParams{}

	// Common parameters
	if flags.limit > 0 {
		params.Limit = flags.limit
	}
	if flags.offset > 0 {
		params.Offset = flags.offset
	}
	if flags.includeField != "" {
		params.Includefield = flags.includeField
	}
	if flags.fuzzyMatch {
		params.FuzzyMatching = true
	}

	// Series-level filters
	if flags.seriesUID != "" {
		params.SeriesInstanceUID = flags.seriesUID
	}
	if flags.modality != "" {
		params.Modality = flags.modality
	}
	if flags.seriesNumber != "" {
		params.SeriesNumber = flags.seriesNumber
	}

	var results []map[string]interface{}
	var err error

	if flags.studyUID != "" {
		// Search series within a specific study
		results, err = client.QidoSearchSeries(flags.studyUID, params)
	} else {
		// Search all series across all studies
		results, err = client.QidoSearchAllSeries(params)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to search series: %w", err)
	}

	return results, nil
}

func runQidoInstances(client interface {
	QidoSearchInstances(studyUID, seriesUID string, params *types.QidoInstanceQueryParams) ([]map[string]interface{}, error)
	QidoSearchStudyInstances(studyUID string, params *types.QidoInstanceQueryParams) ([]map[string]interface{}, error)
	QidoSearchAllInstances(params *types.QidoInstanceQueryParams) ([]map[string]interface{}, error)
}, flags *QidoFlags) ([]map[string]interface{}, error) {
	params := &types.QidoInstanceQueryParams{}

	// Common parameters
	if flags.limit > 0 {
		params.Limit = flags.limit
	}
	if flags.offset > 0 {
		params.Offset = flags.offset
	}
	if flags.includeField != "" {
		params.Includefield = flags.includeField
	}
	if flags.fuzzyMatch {
		params.FuzzyMatching = true
	}

	// Instance-level filters
	if flags.instanceUID != "" {
		params.SOPInstanceUID = flags.instanceUID
	}
	if flags.sopClassUID != "" {
		params.SOPClassUID = flags.sopClassUID
	}

	var results []map[string]interface{}
	var err error

	switch {
	case flags.studyUID != "" && flags.seriesUID != "":
		// Search instances within a specific series
		results, err = client.QidoSearchInstances(flags.studyUID, flags.seriesUID, params)
	case flags.studyUID != "":
		// Search instances within a specific study
		results, err = client.QidoSearchStudyInstances(flags.studyUID, params)
	default:
		// Search all instances
		results, err = client.QidoSearchAllInstances(params)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to search instances: %w", err)
	}

	return results, nil
}

func displayQidoResults(results []map[string]interface{}, jsonOutput bool) error {
	if len(results) == 0 {
		fmt.Println("No results found")
		return nil
	}

	if jsonOutput {
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	// Default to JSON output for QIDO results since the data is complex DICOM metadata
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}
