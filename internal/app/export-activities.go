package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/sebamunozg/coros-api-go/internal/config"
	"github.com/sebamunozg/coros-api-go/internal/coros"
)

type ExportActivitiesOptions struct {
	From       *time.Time
	To         *time.Time
	SportTypes []string
}

func ExportActivities(options ExportActivitiesOptions) error {
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

	// datos de testing
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

	log.Printf("Activities: %+v", result)

	return nil
}
