package weatherAPI

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
	"weather-api-wrapper-service/lib/e"
)

type Client struct {
	Client *http.Client
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{
		Client: &http.Client{
			Timeout: 15 * time.Second,
		},
		apiKey: apiKey,
	}
}

func (c *Client) GetWeather(ctx context.Context, city string) ([]byte, error) {
	basePath := "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline"

	cityPath, err := url.JoinPath(basePath, city)
	if err != nil {
		return nil, e.Wrap("join path", err)
	}

	u, err := url.Parse(cityPath)
	if err != nil {
		return nil, e.Wrap("parse req to url", err)
	}

	// contentType=json
	q := u.Query()

	q.Add("key", c.apiKey)
	q.Add("unitGroup", "metric")

	u.RawQuery = q.Encode()

	log.Println("Получен запрос:", u.String())

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap("new request", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, e.Wrap("get request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response status: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap("read response body", err)
	}

	return body, nil
}
