package app

import (
	"context"
	"fmt"
	"time"

	"github.com/sebamunozg/coros-api-go/internal/coros"
)

type ListActivitiesOptions struct {
	From *time.Time
	To   *time.Time
}

func (a *Application) ListActivities(
	ctx context.Context,
	options ListActivitiesOptions,
) ([]coros.Activity, error) {

	result, err := a.CorosClient.GetActivities(ctx, coros.QueryActivitiesInput{
		PageNumber: 1,
		PageSize:   20,
		From:       options.From,
		To:         options.To,
		// SportTypes: options.SportTypes,
		ModeList: "",
	})

	if err != nil {
		return nil, fmt.Errorf("getting activities: %w", err)
	}

	return result.Activities, nil
}
