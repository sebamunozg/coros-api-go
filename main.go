package main

import (
	"log"

	"github.com/joho/godotenv"
	auth "github.com/sebamunozg/coros-api-go/internal/coros"
	activities "github.com/sebamunozg/coros-api-go/internal/coros/activity"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	accessToken, err := auth.Login()
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

	result, err := activities.GetActivities(input, accessToken)
	if err != nil {
		log.Fatalf("Error getting activities: %v", err)
	}

	log.Printf("Activities: %+v", result)
}
