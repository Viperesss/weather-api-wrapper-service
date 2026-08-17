package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
	"weather-api-wrapper-service/lib/e"
	cache "weather-api-wrapper-service/redis"
)

type Client struct {
	apiKey        string
	weatherClient *http.Client
	cacheClient   *cache.Client
}

func NewClient(cacheClient *cache.Client, apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		weatherClient: &http.Client{
			Timeout: 12 * time.Second,
		},
		cacheClient: cacheClient,
	}
}

func (c *Client) StartServe() error {
	http.HandleFunc("/weather", c.handleWeather)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		return e.Wrap("listen and serve: ", err)
	}
	return nil
}

func (c *Client) handleWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "empty field city", http.StatusBadRequest)
		return
	}

	body, err := c.cacheClient.GetCache(city)
	if err != nil {
		log.Println("failed to get cache:", err)
	}

	if body == nil {
		body, err = c.doRequest(city)
		if err != nil {
			http.Error(w, "failed to get weather", http.StatusBadGateway)
			log.Println(err)
			return
		}

		if err := c.cacheClient.SetCache(city, body); err != nil {
			log.Println("failed to set cache:", err)
		}
	}

	weatherStruct, err := getData(body)
	if err != nil {
		http.Error(w, "failed to get data from request", http.StatusBadGateway)
		log.Println(err)
		return
	}

	dataOutput(w, weatherStruct)
}

func (c *Client) doRequest(city string) ([]byte, error) {
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

	resp, err := c.weatherClient.Do(req)
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

func getData(body []byte) (WeatherReport, error) {
	var ww WeatherReport

	if err := json.Unmarshal(body, &ww); err != nil {
		return WeatherReport{}, e.Wrap("unmarshal body", err)
	}

	return ww, nil
}

func dataOutput(w http.ResponseWriter, weatherStruct WeatherReport) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(weatherStruct); err != nil {
		log.Println("encode response:", err)
	}
}
