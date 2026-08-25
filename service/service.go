package service

import (
	"context"
	"weather-api-wrapper-service/lib/e"
	"weather-api-wrapper-service/weather"
)

func NewService(api Api, cache Cache) *Service {
	return &Service{
		api:   api,
		cache: cache,
	}
}

func (s *Service) GetWeather(ctx context.Context, city string) (weather.WeatherReport, error) {
	weatherCache, err := s.cache.Get(ctx, city)
	if err != nil {
		return weather.WeatherReport{}, e.Wrap("get cache", err)
	}

	if weatherCache != nil {
		return *weatherCache, nil
	}

	weatherStruct, err := s.api.GetWeather(ctx, city)
	if err != nil {
		return weather.WeatherReport{}, e.Wrap("get weather from api", err)
	}

	if err := s.cache.Set(ctx, city, weatherStruct); err != nil {
		return weather.WeatherReport{}, e.Wrap("set cache", err)
	}

	return weatherStruct, nil
}
