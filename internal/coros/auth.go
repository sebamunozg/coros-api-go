package coros

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sebamunozg/coros-api-go/internal/config"
	"github.com/sebamunozg/coros-api-go/internal/core"
	"github.com/sebamunozg/coros-api-go/internal/crypto"
)

type Credentials struct {
	Url      string
	Email    string
	Password string
}

type LoginBody struct {
	Account     string `json:"account"`
	AccountType int    `json:"accountType"`
	Pwd         string `json:"pwd"`
}

func Login(ctx context.Context, cfg config.Credentials) (string, error) {

	payload := LoginBody{
		Account:     cfg.Email,
		AccountType: 2,
		Pwd:         crypto.CreateMD5Hash(cfg.Password),
	}

	b, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to Serialize to JSON from native Go struct type: %v", err)
	}

	body := bytes.NewBuffer(b)

	coros_url := cfg.Url
	coros_url += "/account/login"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, coros_url, body)
	if err != nil {
		log.Printf("Failed to create resource at: %s and the error is: %v\n", coros_url, err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Failed to read response body: %v", err)
	}

	accessToken, found := core.FindAccessToken(data)
	if !found {
		return "", fmt.Errorf("accessToken not found or is not a string")
	}

	return accessToken, nil
}
