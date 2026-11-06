package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sebamunozg/coros-api-go/internal/app"
	"github.com/sebamunozg/coros-api-go/internal/config"
	"github.com/sebamunozg/coros-api-go/internal/coros"
	"github.com/spf13/cobra"
)

var processFitCmd = &cobra.Command{
	Use:   "process-fit",
	Short: "Summarize activities from their FIT files.",
	Long: "Download activities as FIT files, parse them in memory and print one JSON summary " +
		"per activity (NDJSON). Nothing is saved to disk unless --output is given.",
	Args: cobra.NoArgs,
}

type processFlags struct {
	From       string
	To         string
	SportTypes []string
	Output     string
}

func init() {
	var flags processFlags

	processFitCmd.Flags().StringVar(&flags.From, "from", "", "Start date in YYYY-MM-DD format")
	processFitCmd.Flags().StringVar(&flags.To, "to", "", "End date in YYYY-MM-DD format")
	processFitCmd.Flags().StringSliceVar(&flags.SportTypes, "sport-types", []string{}, "Filter by sport types (comma-separated)")
	processFitCmd.Flags().StringVarP(&flags.Output, "output", "o", "", "Write the NDJSON summaries to this file instead of stdout")

	processFitCmd.RunE = func(cmd *cobra.Command, args []string) error {

		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("getting config: %w", err)
		}

		fromDate, err := parseDateFlag(flags.From)
		if err != nil {
			return err
		}
		toDate, err := parseDateFlag(flags.To)
		if err != nil {
			return err
		}

		sportTypes, err := coros.ParseSportTypes(flags.SportTypes)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if flags.Output != "" {
			file, err := os.Create(flags.Output)
			if err != nil {
				return fmt.Errorf("creating output file: %w", err)
			}
			defer file.Close()
			out = file
		}

		application, err := app.NewApplication(cmd.Context(), cfg)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Minute)
		defer cancel()

		result, err := application.ProcessFitActivities(ctx, app.ProcessFitActivitiesOptions{
			From:       fromDate,
			To:         toDate,
			SportTypes: sportTypes,
			Out:        out,
		})
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.ErrOrStderr(), "%d processed, %d failed\n", result.Processed, result.Failed)
		if result.Failed > 0 {
			return fmt.Errorf("%d of %d activities failed", result.Failed, result.Processed+result.Failed)
		}

		return nil
	}
}

func parseDateFlag(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}

	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, fmt.Errorf("invalid date %q, expected YYYY-MM-DD: %w", value, err)
	}

	return &parsed, nil
}
