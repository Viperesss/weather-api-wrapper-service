package cache

import (
	"context"
	"errors"
	"time"
	"weather-api-wrapper-service/lib/e"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
}

func NewClient() *Client {
	return &Client{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
			Protocol: 2,  // connection protocol
		}),
	}
}

func (c *Client) Get(ctx context.Context, city string) ([]byte, error) {
	val, err := c.client.Get(ctx, city).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, e.Wrap("get", err)
	}

	return val, nil
}

func (c *Client) Set(ctx context.Context, city string, body []byte) error {
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
