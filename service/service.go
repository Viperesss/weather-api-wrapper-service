package service

import (
	"context"
	"weather-api-wrapper-service/lib/e"
	cache "weather-api-wrapper-service/redis"
	"weather-api-wrapper-service/weatherAPI"
)

// body, err := c.cacheClient.GetCache(city)
// 	if err != nil {
// 		log.Println("failed to get cache:", err)
// 	}

// 	if body == nil {
// 		body, err = c.doRequest(city)
// 		if err != nil {
// 			http.Error(w, "failed to get weather", http.StatusBadGateway)
// 			log.Println(err)
// 			return
// 		}

// 		if err := c.cacheClient.SetCache(city, body); err != nil {
// 			log.Println("failed to set cache:", err)
// 		}
// 	}

func NewService(api *weatherAPI.Client, cache *cache.Client) *Service {
	return &Service{
		api:   api,
		cache: cache,
	}
}

func (s *Service) GetWeather(ctx context.Context, city string) ([]byte, error) {
	body, err := s.cache.Get(ctx, city)
	if err != nil {
		return nil, e.Wrap("get cache", err)
	}

	if body != nil {
		return body, nil
	}

	body, err = s.api.GetWeather(ctx, city)
	if err != nil {
		return nil, e.Wrap("get weather from api", err)
	}

	if err := s.cache.Set(ctx, city, body); err != nil {
		return nil, e.Wrap("set cache", err)
	}

	return body, nil
}
