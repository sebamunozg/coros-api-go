package activities

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/sebamunozg/coros-api-go/internal/config"
)

type CorosResponseBase struct {
	ApiCode string `json:"apiCode"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

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

// func handle(pageNumber int, pageSize int, from *time.Time, to *time.Time, modeList string) (QueryActivitiesOutput, error) {

// 	input := QueryActivitiesInput{
// 		PageSize:   20,
// 		PageNumber: 1,
// 		From:       from,
// 		To:         to,
// 		ModeList:   modeList,
// 	}

// 	result, err := GetActivities(input, accessToken)

// 	if result.Count == 0 {
// 		fmt.Println("No activities found for the given parameters.")
// 		return QueryActivitiesOutput{}, nil
// 	}

// 	if err != nil {
// 		fmt.Println("Error occurred while fetching activities:", err)
// 		return QueryActivitiesOutput{}, err
// 	}

// 	return QueryActivitiesOutput{
// 		Count:      result.Count,
// 		Activities: result.Activities,
// 	}, err
// }

func GetActivities(input QueryActivitiesInput, accessToken string) (QueryActivitiesOutput, error) {

	cfg, err := config.GetCredentials()
	if err != nil {
		log.Fatal(err)
	}

	activities := []Activity{}
	var currentPage = input.PageNumber
	var lastPage = input.PageNumber

	for currentPage <= lastPage {

		// https://teamapi.coros.com/activity/query

		coros_url := cfg.Url
		coros_url += "/activity/query"

		u, err := url.Parse(coros_url)

		query := u.Query()
		query.Set("size", strconv.Itoa(input.PageSize))
		query.Set("pageNumber", strconv.Itoa(currentPage))
		query.Set("modeList", input.ModeList)

		u.RawQuery = query.Encode()

		if input.From != nil {
			query.Set("from", input.From.Format("20060102"))
		}

		if input.To != nil {
			query.Set("to", input.To.Format("20060102"))
		}

		client := &http.Client{}

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			log.Fatalf("Failed to create resource at: %s and the error is: %v\n", coros_url, err)
		}

		req.Header.Add("accessToken", accessToken)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to create resource at: %s and the error is: %v\n", coros_url, err)
		}

		defer resp.Body.Close()

		log.Printf("Query activity response: %s", resp)

		// Crear validación para verificar:
		// Que existen apiCode, message y result.
		// Que result sea "0000".
		// Que los datos tengan la estructura esperada.

		var data QueryActivitiesResponse

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return QueryActivitiesOutput{}, err
		}

		activities = append(activities, data.Data.DataList...)
		lastPage = data.Data.TotalPage

		currentPage++
	}

	return QueryActivitiesOutput{
		Count:      len(activities),
		Activities: activities,
	}, nil
}
