package cli

import (
	"time"

	"github.com/sebamunozg/coros-api-go/internal/app"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export-activities",
	Short: "Export activities",
	Long:  "Export activities from the Coros API.",
	Args:  cobra.NoArgs,
}

type FileTypeFlag []string
type Flags struct {
	OutputDir  string
	FileTypes  FileTypeFlag
	SportTypes []string
	From       string
	To         string
}

func init() {
	var flags Flags

	exportCmd.Flags().StringVar(&flags.From, "from", "", "Start date in YYYY-MM-DD format")
	exportCmd.Flags().StringVar(&flags.To, "to", "", "End date in YYYY-MM-DD format")
	exportCmd.Flags().StringSliceVar(&flags.SportTypes, "sport-types", []string{}, "Filter by sport types (comma-separated)")

	exportCmd.RunE = func(cmd *cobra.Command, args []string) error {
		var fromDate *time.Time
		var toDate *time.Time

		if flags.From == "" {
			from := time.Now().AddDate(0, -1, 0)
			fromDate = &from
		}

		if flags.To == "" {
			to := time.Now()
			toDate = &to
		}

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

		return app.ExportActivities(app.ExportActivitiesOptions{
			From: fromDate,
			To:   toDate,
		})
	}
}
