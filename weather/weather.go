package weather

import (
	"encoding/json"
	"weather-api-wrapper-service/lib/e"
)

func GetData(body []byte) (WeatherReport, error) {
	var ww WeatherReport

	if err := json.Unmarshal(body, &ww); err != nil {
		return WeatherReport{}, e.Wrap("unmarshal body", err)
	}

	return ww, nil
}
