package app

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/sebamunozg/coros-api-go/internal/coros"
)

type ExportActivitiesOptions struct {
	OutputDir  string
	From       *time.Time
	To         *time.Time
	SportTypes []int
	FileType   coros.ExportType
}

type ActivityFile struct {
	LabelId   string
	SportType int
	FileName  string
}

func (a *Application) ExportActivities(
	ctx context.Context,
	options ExportActivitiesOptions,
) ([]ActivityFile, int, error) {

	input := coros.QueryActivitiesInput{
		PageNumber: 1,
		PageSize:   20,
		From:       options.From,
		To:         options.To,
		SportTypes: options.SportTypes,
		ModeList:   "",
	}

	result, err := a.CorosClient.GetActivities(ctx, input)
	if err != nil {
		return nil, 0, fmt.Errorf("Error getting activities: %w", err)
	}

	activities := make([]ActivityFile, 0, len(result.Activities))

	for _, r := range result.Activities {
		activityDate, err := time.Parse("20060102", strconv.Itoa(r.Date))
		if err != nil {
			return nil, 0, fmt.Errorf("Error parsing activity date: %w", err)
		}

		activityName := strings.TrimSpace(r.Name)
		activityName = strings.NewReplacer("/", "-", "\\", "-").Replace(activityName)

		activities = append(activities, ActivityFile{
			LabelId:   r.LabelId,
			SportType: r.SportType,
			FileName: fmt.Sprintf(
				"%s %s %s.%s",
				activityDate.Format("2006-01-02"),
				activityName,
				r.LabelId,
				options.FileType.String(),
			),
		})
	}

	failedDownloads := 0
	for _, activity := range activities {

		activityData, err := a.CorosClient.DownloadActivity(ctx, coros.DownloadActivitiesDataInput{
			LabelId:   activity.LabelId,
			SportType: activity.SportType,
			FileType:  options.FileType,
		})

		if err != nil {
			log.Printf("Error downloading activity %s: %v", activity.LabelId, err)
			failedDownloads++
			continue
		}

		_, err = a.CorosClient.DownloadActivityFile(
			ctx,
			activityData.FileUrl,
			options.OutputDir,
			activity.FileName,
		)

		if err != nil {
			log.Printf("Error downloading activity file %s: %v", activity.LabelId, err)
			failedDownloads++
			continue
		}
	}

	return activities, failedDownloads, nil
}
