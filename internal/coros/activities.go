package coros

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type QueryActivitiesInput struct {
	PageSize   int        `json:"pageSize"`
	PageNumber int        `json:"pageNumber"`
	From       *time.Time `json:"from"`
	To         *time.Time `json:"to"`
	ModeList   string     `json:"modeList"`
}

type Activity struct {
	Date      int     `json:"date"`
	LabelId   string  `json:"labelId"`
	Name      *string `json:"name"`
	SportType int     `json:"sportType"`
}

type QueryActivitiesData struct {
	Count      int        `json:"count"`
	DataList   []Activity `json:"datalist"`
	PageNumber int        `json:"pageNumber"`
	TotalPage  int        `json:"totalPage"`
}

type QueryActivitiesResponse struct {
	CorosResponseBase
	Data QueryActivitiesData `json:"data"`
}

type QueryActivitiesOutput struct {
	Count      int        `json:"count"`
	Activities []Activity `json:"activities"`
}

func CreateRequest(ctx context.Context, client *http.Client, u *url.URL, accessToken string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accessToken", accessToken)

	return client.Do(req)
}

func validateQueryActivitiesResponse(resp QueryActivitiesResponse) error {

	if resp.Data.TotalPage < 1 {
		return fmt.Errorf("invalid totalPage in Coros response: %d", resp.Data.TotalPage)
	}

	if resp.Data.PageNumber < 1 {
		return fmt.Errorf("invalid pageNumber in Coros response: %d", resp.Data.PageNumber)
	}

	if resp.Data.Count < 0 {
		return fmt.Errorf("invalid count in Coros response: %d", resp.Data.Count)
	}

	return nil
}

func (c *Client) GetActivities(ctx context.Context, input QueryActivitiesInput) (QueryActivitiesOutput, error) {

	activities := []Activity{}
	var currentPage = input.PageNumber
	var lastPage = input.PageNumber

	for currentPage <= lastPage {

		coros_url := c.baseURL
		coros_url += "/activity/query"

		u, err := url.Parse(coros_url)

		query := u.Query()
		query.Set("size", strconv.Itoa(input.PageSize))
		query.Set("pageNumber", strconv.Itoa(currentPage))
		query.Set("modeList", input.ModeList)

		if input.From != nil {
			query.Set("from", input.From.Format("20060102"))
		}

		if input.To != nil {
			query.Set("to", input.To.Format("20060102"))
		}

		u.RawQuery = query.Encode()

		resp, err := CreateRequest(ctx, c.httpClient, u, c.accessToken)
		if err != nil {
			return QueryActivitiesOutput{}, err
		}

		if resp.StatusCode != http.StatusOK {
			return QueryActivitiesOutput{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		var data QueryActivitiesResponse

		decoderErr := json.NewDecoder(resp.Body).Decode(&data)
		closeErr := resp.Body.Close()

		if decoderErr != nil {
			return QueryActivitiesOutput{}, decoderErr
		}

		if closeErr != nil {
			return QueryActivitiesOutput{}, closeErr
		}

		if err := AssertCorosResponseBase(data.CorosResponseBase); err != nil {
			return QueryActivitiesOutput{}, err
		}

		if err := validateQueryActivitiesResponse(data); err != nil {
			return QueryActivitiesOutput{}, err
		}

		log.Printf("Query activity data: %+v", data)

		activities = append(activities, data.Data.DataList...)
		lastPage = data.Data.TotalPage

		currentPage++
	}

	return QueryActivitiesOutput{
		Count:      len(activities),
		Activities: activities,
	}, nil
}
