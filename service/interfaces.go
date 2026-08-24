package service

import "context"

type Api interface {
	GetWeather(ctx context.Context, city string) ([]byte, error)
}

type Cache interface {
	Get(ctx context.Context, city string) ([]byte, error)
	Set(ctx context.Context, city string, body []byte) error
}

type Service struct {
	api   Api
	cache Cache
}
