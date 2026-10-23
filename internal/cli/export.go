package cli

import (
	"log"
	"strconv"
	"time"

	"github.com/sebamunozg/coros-api-go/internal/app"
	"github.com/sebamunozg/coros-api-go/internal/coros"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export-activities",
	Short: "Export activities",
	Long:  "Export activities from the Coros API.",
	Args:  cobra.NoArgs,
}

type Flags struct {
	OutputDir  string
	FileType   coros.ExportType
	SportTypes []string
	From       string
	To         string
}

type exportTypeValue coros.ExportType

func (f *exportTypeValue) String() string {
	return coros.ExportType(*f).String()
}

func (f *exportTypeValue) Set(value string) error {
	exportType, err := coros.ParseExportType(value)
	if err != nil {
		return err
	}

	*f = exportTypeValue(exportType)
	return nil
}

func (f *exportTypeValue) Type() string {
	return "export-type"
}

func init() {
	var flags Flags

	exportCmd.Flags().StringVar(&flags.From, "from", "", "Start date in YYYY-MM-DD format")
	exportCmd.Flags().StringVar(&flags.To, "to", "", "End date in YYYY-MM-DD format")
	exportCmd.Flags().StringSliceVar(&flags.SportTypes, "sport-types", []string{}, "Filter by sport types (comma-separated)")
	exportCmd.Flags().Var((*exportTypeValue)(&flags.FileType), "export-type", "Export format: fit, tcx or gpx")
	exportCmd.Flags().StringVarP(&flags.OutputDir, "output-dir", "o", "", "Output directory for downloaded files")
	if err := exportCmd.MarkFlagRequired("output-dir"); err != nil {
		log.Fatal(err)
	}

	exportCmd.RunE = func(cmd *cobra.Command, args []string) error {
		var fromDate *time.Time
		var toDate *time.Time

		if flags.From != "" {
			parsed, err := time.Parse("2006-01-02", flags.From)
			if err != nil {
				return err
			}

			fromDate = &parsed
		}

		if flags.To != "" {
			parsed, err := time.Parse("2006-01-02", flags.To)
			if err != nil {
				return err
			}

			toDate = &parsed
		}

		var sportTypeCodes []int
		for _, sportType := range flags.SportTypes {
			code, err := coros.GetSportTypeCode(sportType)
			if err != nil {
				return err
			}

			sportType, err := strconv.Atoi(code)
			if err != nil {
				return err
			}
			sportTypeCodes = append(sportTypeCodes, sportType)
		}

		_, failedDownloads, err := app.ExportActivities(app.ExportActivitiesOptions{
			From:       fromDate,
			To:         toDate,
			FileType:   flags.FileType,
			SportTypes: sportTypeCodes,
			OutputDir:  flags.OutputDir,
		})

		if failedDownloads > 0 {
			log.Printf("Failed to download %d activity files", failedDownloads)
		}

		return err
	}
}
