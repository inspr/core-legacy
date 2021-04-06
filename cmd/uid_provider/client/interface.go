package redisclient

import (
	"context"
)

type RedisClient interface {
	Create(ctx context.Context, data interface{}) interface{}
	Get(ctx context.Context, key string) interface{}
	Update(ctx context.Context, data interface{}) interface{}
	Delete(ctx context.Context, key string) interface{}
}
