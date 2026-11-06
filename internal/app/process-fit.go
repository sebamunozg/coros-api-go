package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/sebamunozg/coros-api-go/internal/coros"
	"github.com/sebamunozg/coros-api-go/internal/fit"
)

type ProcessFitActivitiesOptions struct {
	From       *time.Time
	To         *time.Time
	SportTypes []int
	// Out receives one JSON object (ActivitySummary) per line.
	Out io.Writer
}

type ProcessFitResult struct {
	Processed int
	Failed    int
}

// ProcessFitActivities downloads each matching activity as a FIT file, parses
// it and writes its summary to options.Out as NDJSON. Activities are handled
// one at a time, so memory use does not grow with the number of activities.
// A failing activity is logged and counted, and processing continues.
func (a *Application) ProcessFitActivities(
	ctx context.Context,
	options ProcessFitActivitiesOptions,
) (ProcessFitResult, error) {

	result, err := a.CorosClient.GetActivities(ctx, coros.QueryActivitiesInput{
		PageNumber: 1,
		PageSize:   20,
		From:       options.From,
		To:         options.To,
		SportTypes: options.SportTypes,
	})
	if err != nil {
		return ProcessFitResult{}, fmt.Errorf("getting activities: %w", err)
	}

	enc := json.NewEncoder(options.Out)
	var res ProcessFitResult

	for _, activity := range result.Activities {
		summary, err := a.processFitActivity(ctx, activity)
		if err != nil {
			log.Printf("Error processing activity %s: %v", activity.LabelID, err)
			res.Failed++
			continue
		}

		if err := enc.Encode(summary); err != nil {
			return res, fmt.Errorf("writing summary: %w", err)
		}
		res.Processed++
	}

	return res, nil
}

func (a *Application) processFitActivity(
	ctx context.Context,
	activity coros.DataList,
) (*fit.ActivitySummary, error) {

	detail, err := a.CorosClient.DownloadActivity(ctx, coros.DownloadActivitiesDataInput{
		LabelId:   activity.LabelID,
		SportType: activity.SportType,
		FileType:  coros.ExportTypeFIT,
	})
	if err != nil {
		return nil, fmt.Errorf("requesting download: %w", err)
	}

	body, err := a.CorosClient.OpenActivityFile(ctx, detail.FileUrl)
	if err != nil {
		return nil, fmt.Errorf("downloading FIT file: %w", err)
	}
	defer body.Close()

	summary, err := fit.ParseActivity(body)
	if err != nil {
		return nil, err
	}

	summary.LabelID = activity.LabelID
	return summary, nil
}
