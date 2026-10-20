package service

import (
	"context"
	"fmt"
	"weather-api-wrapper-service/lib/e"
	"weather-api-wrapper-service/weather"
)

func NewService(api Api, cache Cache) *Service {
	return &Service{
		api:   api,
		cache: cache,
	}
}

func NewServiceWithInsights(api Api, cache Cache, insights InsightGenerator) *Service {
	s := NewService(api, cache)
	s.insights = insights
	return s
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

func (s *Service) GetInsight(ctx context.Context, city, businessContext string) (weather.BusinessInsight, error) {
	if s.insights == nil {
		return weather.BusinessInsight{}, e.Wrap("generate business insight", fmt.Errorf("AI insights are not configured"))
	}
	report, err := s.GetWeather(ctx, city)
	if err != nil {
		return weather.BusinessInsight{}, e.Wrap("get forecast for insight", err)
	}
	insight, err := s.insights.GenerateInsight(ctx, report, businessContext)
	if err != nil {
		return weather.BusinessInsight{}, e.Wrap("generate business insight", err)
	}
	return insight, nil
}
