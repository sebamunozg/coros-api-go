package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sebamunozg/coros-api-go/internal/config"
	"github.com/sebamunozg/coros-api-go/internal/coros"
)

type Application struct {
	CorosClient *coros.Client
}

func NewApplication(parentCtx context.Context, cfg config.Config) (*Application, error) {

	client := coros.NewClient(
		cfg.Url,
		cfg.Email,
		cfg.Password,
		"",
		http.DefaultClient,
	)

	ctx, cancel := context.WithTimeout(parentCtx, cfg.RequestTimeout)
	defer cancel()

	if err := client.Login(ctx); err != nil {
		return nil, fmt.Errorf("logging in: %w", err)
	}

	return &Application{
		CorosClient: client,
	}, nil
}
