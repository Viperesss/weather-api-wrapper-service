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
	"weather-api-wrapper-service/weather"
)

type client struct {
	Client *http.Client
	apiKey string
}

func NewClient(apiKey string) *client {
	return &client{
		Client: &http.Client{
			Timeout: 15 * time.Second,
		},
		apiKey: apiKey,
	}
}

func (c *client) GetWeather(ctx context.Context, city string) (weather.WeatherReport, error) {
	basePath := "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline"

	cityPath, err := url.JoinPath(basePath, city)
	if err != nil {
		return weather.WeatherReport{}, e.Wrap("join path", err)
	}

	u, err := url.Parse(cityPath)
	if err != nil {
		return weather.WeatherReport{}, e.Wrap("parse req to url", err)
	}

	// contentType=json
	q := u.Query()

	q.Add("key", c.apiKey)
	q.Add("unitGroup", "metric")

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return weather.WeatherReport{}, e.Wrap("new request", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return weather.WeatherReport{}, e.Wrap("get request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return weather.WeatherReport{}, fmt.Errorf("bad response status: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return weather.WeatherReport{}, e.Wrap("read response body", err)
	}

	weatherStruct, err := weather.GetData(body)
	if err != nil {
		return weather.WeatherReport{}, e.Wrap("get data", err)
	}

	log.Println("Выполнен запрос во внешний API:", city)

	return weatherStruct, nil
}
