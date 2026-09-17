package auth

import (
	"bytes"
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

func Login() {

	cfg, err := config.GetCredentials()
	if err != nil {
		log.Fatal(err)
	}

	payload := LoginBody{
		Account:     cfg.Email,
		AccountType: 2,
		Pwd:         crypto.CreateMD5Hash(cfg.Password),
	}

	b, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to Serialize to JSON from native Go struct type: %v", err)
	}

	body := bytes.NewBuffer(b)

	fmt.Println(body.String())

	coros_url := cfg.Url
	coros_url += "/account/login"

	resp, err := http.Post(coros_url, "application/json; charset=utf-8", body)
	if err != nil {
		log.Fatalf("Failed to create resource at: %s and the error is: %v\n", coros_url, err)
	}
	defer resp.Body.Close()

	log.Printf("Status received from server is: %s", resp.Status)
	log.Printf("StatusCode received from server is: %d", resp.StatusCode)
	log.Printf("Content Type received from Server is: %s", resp.Header["Content-Type"][0])

	data := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	accessToken, found := core.FindAccessToken(data)
	if !found {
		log.Println("accessToken no existe o no es un string")
		return
	}

	fmt.Printf("accessToken: %s\n", accessToken)

	for key, value := range data {
		fmt.Printf("%s: %v\n", key, value)
	}
}
