package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/sebamunozg/coros-api-go/internal/config"
	"github.com/sebamunozg/coros-api-go/internal/coros"
)

type Application struct {
	CorosClient *coros.Client
}

func NewApplication() (*Application, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("loading .env file: %w", err)
	}

	cfg, err := config.GetCredentials()
	if err != nil {
		return nil, fmt.Errorf("getting credentials: %w", err)
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
		return nil, fmt.Errorf("logging in: %w", err)
	}

	return &Application{
		CorosClient: client,
	}, nil
}
