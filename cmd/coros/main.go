package main

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/sebamunozg/coros-api-go/internal/config"
	coros "github.com/sebamunozg/coros-api-go/internal/coros"
	activities "github.com/sebamunozg/coros-api-go/internal/coros/activity"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg, err := config.GetCredentials()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	accessToken, err := coros.Login(ctx, cfg)
	if err != nil {
		log.Fatalf("Error logging in: %v", err)
	}

	// datos de testing
	input := activities.QueryActivitiesInput{
		PageNumber: 1,
		PageSize:   20,
		From:       nil,
		To:         nil,
		ModeList:   "",
	}

	result, err := activities.GetActivities(ctx, cfg, input, accessToken)
	if err != nil {
		log.Fatalf("Error getting activities: %v", err)
	}

	log.Printf("Activities: %+v", result)
}
