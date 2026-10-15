package coros

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

type DownloadActivitiesDataInput struct {
	LabelId   string     `json:"labelId"`
	SportType int        `json:"sportType"`
	FileType  ExportType `json:"fileType"`
}

type DownloadActivityDetailData struct {
	FileUrl string `json:"fileUrl"`
}

type DownloadActivityDetailResponse struct {
	CorosResponseBase
	Data DownloadActivityDetailData `json:"data"`
}

func createPostRequest(ctx context.Context, client *http.Client, u *url.URL, accessToken string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accessToken", accessToken)

	return client.Do(req)
}

func validateDownloadActivityDetailResponse(resp DownloadActivityDetailResponse) error {

	if resp.Data.FileUrl == "" {
		return fmt.Errorf("missing fileUrl in Coros response")
	}

	return nil
}

func (c *Client) DownloadActivity(ctx context.Context, input DownloadActivitiesDataInput) (DownloadActivityDetailData, error) {
	coros_url := c.baseURL
	coros_url += "/activity/detail/download"

	u, err := url.Parse(coros_url)
	if err != nil {
		return DownloadActivityDetailData{}, fmt.Errorf("failed to parse URL: %v", err)
	}

	query := u.Query()
	query.Set("labelId", input.LabelId)
	query.Set("sportType", strconv.Itoa(input.SportType))
	query.Set("fileType", strconv.Itoa(int(input.FileType)))

	u.RawQuery = query.Encode()

	resp, err := createPostRequest(ctx, c.httpClient, u, c.accessToken)
	if err != nil {
		return DownloadActivityDetailData{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return DownloadActivityDetailData{}, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	var data DownloadActivityDetailResponse

	decoderErr := json.NewDecoder(resp.Body).Decode(&data)
	closeErr := resp.Body.Close()

	if decoderErr != nil {
		return DownloadActivityDetailData{}, decoderErr
	}

	if closeErr != nil {
		return DownloadActivityDetailData{}, closeErr
	}

	if err := AssertCorosResponseBase(data.CorosResponseBase); err != nil {
		return DownloadActivityDetailData{}, err
	}

	if err := validateDownloadActivityDetailResponse(data); err != nil {
		return DownloadActivityDetailData{}, err
	}

	return data.Data, nil
}

func (c *Client) DownloadActivityFile(ctx context.Context, fileUrl string, directory string, fileName string) ([]byte, error) {
	log.Printf("Downloading %s", fileName)

	req, err := http.NewRequestWithContext(ctx, "GET", fileUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	if err := os.MkdirAll(directory, 0o755); err != nil {
		return nil, fmt.Errorf("creando directorio: %w", err)
	}

	path := filepath.Join(directory, fileName)
	out, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("creando archivo: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return nil, fmt.Errorf("escribiendo archivo: %w", err)
	}
	return nil, nil
}
