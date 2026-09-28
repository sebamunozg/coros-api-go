package coros

import "fmt"

type CorosResponseBase struct {
	ApiCode string `json:"apiCode"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func AssertCorosResponseBase(resp CorosResponseBase) error {
	if resp.ApiCode == "" {
		return fmt.Errorf("missing apiCode in Coros response")
	}
	if resp.Message == "" {
		return fmt.Errorf("missing message in Coros response")
	}
	if resp.Result == "" {
		return fmt.Errorf("missing result in Coros response")
	}

	if resp.Result != "0000" {
		return fmt.Errorf("Coros API error: %s - %s", resp.Result, resp.Message)
	}

	return nil
}
