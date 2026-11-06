package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/zalando/go-keyring"
)

const (
	defaultOutputDir      = "./activities"
	defaultRequestTimeout = 30 * time.Second
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

// GetConfig resolves the configuration. Credentials saved by `coros auth`
// take priority; COROS_API_URL, COROS_EMAIL and COROS_PASSWORD (from the
// environment or an optional .env file) are the fallback.
func GetConfig() (Config, error) {

	// .env is optional.
	_ = godotenv.Load(".env")

	stored, err := LoadStoredCredentials()
	if err != nil {
		return Config{}, err
	}

	apiUrl := firstNonEmpty(stored.URL, os.Getenv("COROS_API_URL"))
	if apiUrl == "" {
		return Config{}, errors.New("no API URL configured: run `coros auth` or set COROS_API_URL")
	}

	email := firstNonEmpty(stored.Email, os.Getenv("COROS_EMAIL"))
	if email == "" {
		return Config{}, errors.New("no email configured: run `coros auth` or set COROS_EMAIL")
	}

	password, err := keyring.Get(KeyringService, email)
	if err != nil || password == "" {
		password = os.Getenv("COROS_PASSWORD")
	}
	if password == "" {
		return Config{}, errors.New("no password configured: run `coros auth` or set COROS_PASSWORD")
	}

	outputDir, timeout, err := loadFileConfig(".config/coros/config.json")
	if err != nil {
		return Config{}, err
	}

	return Config{
		Url:            apiUrl,
		Email:          email,
		Password:       password,
		OutputDir:      outputDir,
		RequestTimeout: timeout,
	}, nil
}

// loadFileConfig reads the optional config.json. A missing file yields the defaults.
func loadFileConfig(path string) (string, time.Duration, error) {
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return defaultOutputDir, defaultRequestTimeout, nil
	}
	if err != nil {
		return "", 0, fmt.Errorf("opening config.json: %w", err)
	}
	defer file.Close()

	var fileCfg fileConfig
	if err := json.NewDecoder(file).Decode(&fileCfg); err != nil {
		return "", 0, fmt.Errorf("decoding config.json: %w", err)
	}

	if fileCfg.OutputDir == "" {
		return "", 0, errors.New("output_dir is empty in config.json")
	}

	if fileCfg.RequestTimeoutSecs <= 0 {
		return "", 0, errors.New("request_timeout_seconds must be greater than zero")
	}

	return fileCfg.OutputDir, time.Duration(fileCfg.RequestTimeoutSecs) * time.Second, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
