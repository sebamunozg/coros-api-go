package config

import (
	"fmt"
	"os"
)

type Credentials struct {
	Url      string
	Email    string
	Password string
}

func GetCredentials() (Credentials, error) {
	cfg := Credentials{
		Url:      os.Getenv("COROS_API_URL"),
		Email:    os.Getenv("COROS_EMAIL"),
		Password: os.Getenv("COROS_PASSWORD"),
	}

	if cfg.Url == "" {
		return cfg, fmt.Errorf("COROS_API_URL is not set")
	}
	if cfg.Email == "" {
		return cfg, fmt.Errorf("COROS_EMAIL is not set")
	}
	if cfg.Password == "" {
		return cfg, fmt.Errorf("COROS_PASSWORD is not set")
	}

	return cfg, nil
}
