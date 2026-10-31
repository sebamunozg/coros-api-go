package cli

import (
	"context"
	"time"

	"github.com/sebamunozg/coros-api-go/internal/app"
	"github.com/spf13/cobra"
)

var processFitCmd = &cobra.Command{
	Use:   "process-fit",
	Short: "Process last FIT files.",
	Long:  "Process FIT files and extract activities from the Coros API.",
	Args:  cobra.NoArgs,
}

type processFlags struct {
	From string
	To   string
}

func init() {
	var flags processFlags

	processFitCmd.Flags().StringVar(&flags.From, "from", "", "From date")
	processFitCmd.Flags().StringVar(&flags.To, "to", "", "To date")

	processFitCmd.RunE = func(cmd *cobra.Command, args []string) error {

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

		application, err := app.NewApplication()
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(cmd.Context(), 10*time.Second)
		defer cancel()

		err = application.ProcessFitActivities(ctx, app.ProcessFitActivitiesOptions{
			From: fromDate,
			To:   toDate,
		})
		if err != nil {
			return err
		}

		return err
	}

}
