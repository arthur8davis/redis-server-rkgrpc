package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	client *redis.Client
}

func New(client *redis.Client) Cache {
	return Cache{client: client}
}

func (r Cache) Set(ctx context.Context, key, value string, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(ctx, key, value, expiration)
}

func (r Cache) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.client.Get(ctx, key)
}

func (r Cache) Health(ctx context.Context) *redis.StatusCmd {
	return r.client.Ping(ctx)
}
