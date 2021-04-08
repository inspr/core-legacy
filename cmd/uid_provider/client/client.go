package redisclient

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	client *redis.Client
}

func NewRedisClient() *Client {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

	return &Client{
		client: redis.NewClient(&redis.Options{
			Addr:     host + ":" + port,
			Password: password,
			DB:       0, // use default DB
		}),
	}
}

func (c *Client) Create(ctx context.Context, data interface{}) interface{} {
	return nil
}

func (c *Client) Get(ctx context.Context, key string) interface{} {
	return nil
}

func (c *Client) Update(ctx context.Context, data interface{}) interface{} {
	return nil
}

func (c *Client) Delete(ctx context.Context, key string) interface{} {
	return nil
}
