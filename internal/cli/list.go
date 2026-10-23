package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/sebamunozg/coros-api-go/internal/app"
	h "github.com/sebamunozg/coros-api-go/internal/helpers"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list-activities",
	Short: "List activities",
	Long:  "List activities from the Coros API.",
	Args:  cobra.NoArgs,
}

type listFlags struct {
	From string
	To   string
}

func init() {
	var flags listFlags

	listCmd.Flags().StringVar(&flags.From, "from", "", "Start date in YYYY-MM-DD format")
	listCmd.Flags().StringVar(&flags.To, "to", "", "End date in YYYY-MM-DD format")

	listCmd.RunE = func(cmd *cobra.Command, args []string) error {
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

		activities, err := application.ListActivities(ctx, app.ListActivitiesOptions{
			From: fromDate,
			To:   toDate,
		})
		if err != nil {
			return err
		}

		for _, activity := range activities {
			fmt.Fprintf(
				cmd.OutOrStdout(),
				"%s - %s - %s - sport type: %d\n",
				h.FormatActivityDate(activity.Date),
				activity.Name,
				activity.LabelId,
				activity.SportType,
			)
		}

		return err
	}
}
