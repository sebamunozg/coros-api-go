package main

import (
	"log"

	"github.com/joho/godotenv"
	auth "github.com/sebamunozg/coros-api-go/internal/coros"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	auth.Login()
}
