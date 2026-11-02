package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Url            string
	Email          string
	Password       string
	OutputDir      string
	RequestTimeout time.Duration
}

type fileConfig struct {
	OutputDir          string `json:"output_dir"`
	RequestTimeoutSecs int    `json:"request_timeout_seconds"`
}

func GetConfig() (Config, error) {

	apiUrl := os.Getenv("COROS_API_URL")
	if apiUrl == "" {
		return Config{}, fmt.Errorf("COROS_API_URL is not set")
	}

	email := os.Getenv("COROS_EMAIL")
	if email == "" {
		return Config{}, fmt.Errorf("COROS_EMAIL is not set")
	}

	password := os.Getenv("COROS_PASSWORD")
	if password == "" {
		return Config{}, fmt.Errorf("COROS_PASSWORD is not set")
	}

	file, err := os.Open(".config/coros/config.json")
	if err != nil {
		return Config{}, fmt.Errorf("opening config.json: %w", err)
	}
	defer file.Close()

	var fileCfg fileConfig
	if err := json.NewDecoder(file).Decode(&fileCfg); err != nil {
		return Config{}, fmt.Errorf("decoding config.json: %w", err)
	}

	if fileCfg.OutputDir == "" {
		return Config{}, fmt.Errorf("output_dir is empty in config.json")
	}

	if fileCfg.RequestTimeoutSecs <= 0 {
		return Config{}, fmt.Errorf(
			"request_timeout_seconds must be greater than zero",
		)
	}

	cfg := Config{
		Url:            os.Getenv("COROS_API_URL"),
		Email:          os.Getenv("COROS_EMAIL"),
		Password:       os.Getenv("COROS_PASSWORD"),
		OutputDir:      fileCfg.OutputDir,
		RequestTimeout: time.Duration(fileCfg.RequestTimeoutSecs) * time.Second,
	}

	return cfg, nil
}
