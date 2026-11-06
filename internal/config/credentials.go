package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"
)

// KeyringService is the service name under which the password is stored.
const KeyringService = "coros-api-go"

// StoredCredentials holds the non-secret values saved by `coros auth`.
// The password lives in the OS keyring.
type StoredCredentials struct {
	URL   string `json:"api_url"`
	Email string `json:"email"`
}

func credentialsPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("finding user config directory: %w", err)
	}

	return filepath.Join(dir, "coros", "credentials.json"), nil
}

// LoadStoredCredentials returns what `coros auth` saved, or the zero value if nothing was saved yet.
func LoadStoredCredentials() (StoredCredentials, error) {
	path, err := credentialsPath()
	if err != nil {
		return StoredCredentials{}, nil
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return StoredCredentials{}, nil
	}
	if err != nil {
		return StoredCredentials{}, fmt.Errorf("reading %s: %w", path, err)
	}

	var stored StoredCredentials
	if err := json.Unmarshal(data, &stored); err != nil {
		return StoredCredentials{}, fmt.Errorf("decoding %s: %w", path, err)
	}

	return stored, nil
}

// SaveCredentials stores the API URL and email on disk and the password in the OS keyring.
func SaveCredentials(apiURL, email, password string) error {
	path, err := credentialsPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(StoredCredentials{URL: apiURL, Email: email}, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	if err := keyring.Set(KeyringService, email, password); err != nil {
		return fmt.Errorf("saving password in the keyring: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}

	return nil
}
