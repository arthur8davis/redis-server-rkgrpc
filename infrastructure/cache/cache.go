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

func (c Cache) Set(ctx context.Context, key, value string, expiration time.Duration) *redis.StatusCmd {
	return c.client.Set(ctx, key, value, expiration)
}

func (c Cache) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.client.Get(ctx, key)
}

func (c Cache) Del(ctx context.Context, key string) *redis.IntCmd {
	return c.client.Del(ctx, key)
}

func (c Cache) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return c.client.Expire(ctx, key, expiration)
}

func (c Cache) HSet(ctx context.Context, hash, key, value string) *redis.IntCmd {
	return c.client.HSet(ctx, hash, key, value)
}

func (c Cache) HGet(ctx context.Context, hash, key string) *redis.StringCmd {
	return c.client.HGet(ctx, hash, key)
}

func (c Cache) HDel(ctx context.Context, hash, key string) *redis.IntCmd {
	return c.client.HDel(ctx, hash, key)
}

func (c Cache) Health(ctx context.Context) *redis.StatusCmd {
	return c.client.Ping(ctx)
}
