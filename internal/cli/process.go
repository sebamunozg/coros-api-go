package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/sebamunozg/coros-api-go/internal/app"
	"github.com/sebamunozg/coros-api-go/internal/config"
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

		// Cargar variables de configuración
		if err := godotenv.Load(".env"); err != nil {
			return fmt.Errorf("loading .env file: %w", err)
		}

		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("getting config: %w", err)
		}

		application, err := app.NewApplication(cmd.Context(), cfg)
		if err != nil {
			return err
		}

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
