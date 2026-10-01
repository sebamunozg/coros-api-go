package coros

import "net/http"

type Client struct {
	baseURL     string
	httpClient  *http.Client
	email       string
	password    string
	accessToken string
}

func NewClient(baseURL, email, password, accessToken string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		baseURL:     baseURL,
		httpClient:  httpClient,
		email:       email,
		password:    password,
		accessToken: accessToken,
	}
}
