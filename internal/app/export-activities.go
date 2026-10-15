package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sebamunozg/coros-api-go/internal/config"
	"github.com/sebamunozg/coros-api-go/internal/coros"
)

type ExportActivitiesOptions struct {
	OutputDir  string
	From       *time.Time
	To         *time.Time
	SportTypes []string
	FileType   coros.ExportType
}

type ActivityFile struct {
	LabelId   string
	SportType int
	FileName  string
}

func ExportActivities(options ExportActivitiesOptions) ([]ActivityFile, error) {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg, err := config.GetCredentials()
	if err != nil {
		log.Fatal(err)
	}

	client := coros.NewClient(
		cfg.Url,
		cfg.Email,
		cfg.Password,
		"",
		http.DefaultClient,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Login(ctx); err != nil {
		log.Fatalf("Error logging in: %v", err)
	}

	input := coros.QueryActivitiesInput{
		PageNumber: 1,
		PageSize:   20,
		From:       options.From,
		To:         options.To,
		ModeList:   "",
	}

	result, err := client.GetActivities(ctx, input)
	if err != nil {
		log.Fatalf("Error getting activities: %v", err)
	}

	activities := make([]ActivityFile, 0, len(result.Activities))

	for _, r := range result.Activities {
		activityDate, err := time.Parse("20060102", strconv.Itoa(r.Date))
		if err != nil {
			return nil, err
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

		activityData, err := client.DownloadActivity(ctx, coros.DownloadActivitiesDataInput{
			LabelId:   activity.LabelId,
			SportType: activity.SportType,
			FileType:  options.FileType,
		})

		_, err = client.DownloadActivityFile(ctx, activityData.FileUrl, options.OutputDir, activity.FileName)

		if err != nil {
			log.Printf("Error downloading activity %s: %v", activity.LabelId, err)
			failedDownloads++
			continue
		}
	}

	return activities, nil
}
