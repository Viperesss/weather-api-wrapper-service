package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"weather-api-wrapper-service/lib/e"
	"weather-api-wrapper-service/weather"

	"github.com/redis/go-redis/v9"
)

type client struct {
	client *redis.Client
}

func NewClient() *client {
	return &client{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
			Protocol: 2,  // connection protocol
		}),
	}
}

func (c *client) Get(ctx context.Context, city string) (*weather.WeatherReport, error) {
	val, err := c.client.Get(ctx, city).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, e.Wrap("get", err)
	}

	var weatherStruct weather.WeatherReport

	if err := json.Unmarshal(val, &weatherStruct); err != nil {
		return nil, e.Wrap("unmarshal json", err)
	}

	return &weatherStruct, nil
}

func (c *client) Set(ctx context.Context, city string, weatherStruct weather.WeatherReport) error {
	body, err := json.Marshal(weatherStruct)
	if err != nil {
		return e.Wrap("marshal json", err)
	}

	if err := c.client.Set(ctx, city, body, 15*time.Minute).Err(); err != nil {
		return e.Wrap("set", err)
	}

	return nil
}

// func (c *Client) ping() error {
// 	status, err := c.client.Ping(context.Background()).Result()
// 	if err != nil {
// 		return e.Wrap("missing connection to redis", err)
// 	}
// 	log.Println(status)
// 	return nil
// }
