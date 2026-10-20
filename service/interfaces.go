package service

import (
	"context"
	"weather-api-wrapper-service/weather"
)

type Api interface {
	GetWeather(ctx context.Context, city string) (weather.WeatherReport, error)
}

type Cache interface {
	Get(ctx context.Context, city string) (*weather.WeatherReport, error)
	Set(ctx context.Context, city string, weatherStruct weather.WeatherReport) error
}

type InsightGenerator interface {
	GenerateInsight(ctx context.Context, report weather.WeatherReport, businessContext string) (weather.BusinessInsight, error)
}

type Service struct {
	api      Api
	cache    Cache
	insights InsightGenerator
}
