package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sebamunozg/coros-api-go/internal/coros"
	"github.com/sebamunozg/coros-api-go/internal/fit"
)

type ProcessFitActivitiesOptions struct {
	From *time.Time
	To   *time.Time
}

func (a *Application) ProcessFitActivities(
	ctx context.Context,
	options ProcessFitActivitiesOptions,
) error {

	input := coros.QueryActivitiesInput{
		PageNumber: 1,
		PageSize:   20,
		From:       options.From,
		To:         options.To,
		SportTypes: []int{100},
		ModeList:   "",
	}

	result, err := a.CorosClient.GetActivities(ctx, input)
	if err != nil {
		return fmt.Errorf("Error getting activities: %w", err)
	}

	activities := make([]ActivityFile, 0, len(result.Activities))

	for _, r := range result.Activities {
		activityDate, err := time.Parse("20060102", strconv.Itoa(r.Date))
		if err != nil {
			log.Printf("Error parsing activity date: %v", err)
			return fmt.Errorf("Error parsing activity date: %w", err)
		}

		activityName := strings.TrimSpace(r.Name)
		activityName = strings.NewReplacer("/", "-", "\\", "-").Replace(activityName)

		activities = append(activities, ActivityFile{
			LabelId:   r.LabelID,
			SportType: r.SportType,
			FileName: fmt.Sprintf(
				"%s %s %s.%s",
				activityDate.Format("2006-01-02"),
				activityName,
				r.LabelID,
				strconv.Itoa(100),
			),
		})
	}

	failedDownloads := 0
	for _, activity := range activities {

		detail, err := a.CorosClient.DownloadActivity(ctx, coros.DownloadActivitiesDataInput{
			LabelId:   activity.LabelId,
			SportType: activity.SportType,
			FileType:  coros.ExportTypeFIT,
		})

		if err != nil {
			log.Printf("Error downloading activity %s: %v", activity.LabelId, err)
			failedDownloads++
			continue
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, detail.FileUrl, nil)
		if err != nil {
			return fmt.Errorf("creating FIT request: %w", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("downloading FIT file: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return fmt.Errorf("unexpected FIT download status: %s", resp.Status)
		}

		processErr := fit.ProcessFIT(resp.Body)

		if processErr != nil {
			log.Printf("Error processing activity file %s: %v", activity.LabelId, processErr)
			failedDownloads++
			continue
		}

	}

	log.Printf("Failed downloads: %d", failedDownloads)
	return nil
}
