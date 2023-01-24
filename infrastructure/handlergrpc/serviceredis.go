package handlergrpc

import (
	"context"
	"fmt"
	"rhyme80/redis-server-rkgrpc/infrastructure/tracing"
	"rhyme80/redis-server-rkgrpc/model"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"rhyme80/redis-server-rkgrpc/infrastructure/cache"
	redisGrpc "rhyme80/redis-server-rkgrpc/infrastructure/servicegrpc/proto"
)

type HandlerGrpc struct {
	redisCache cache.Cache
	tracing    tracing.Tracing

	redisGrpc.UnimplementedRedisServiceServer
}

func New(redisCache cache.Cache, tracing tracing.Tracing) *HandlerGrpc {
	return &HandlerGrpc{redisCache: redisCache, tracing: tracing}
}

func (h *HandlerGrpc) Get(ctx context.Context, get *redisGrpc.RequestGet) (*redisGrpc.ResponseGet, error) {
	if err := validateHeaders(ctx); err != nil {
		return nil, err
	}

	key := get.GetKey()

	timeInit := time.Now()

	stringData := h.redisCache.Get(ctx, key)

	value, err := stringData.Result()
	h.tracing.Logger(ctx, model.DataLogger{
		Operation:    "getCache",
		StartTime:    timeInit,
		EndTime:      time.Now(),
		ElapsedNano:  time.Since(timeInit).Nanoseconds(),
		Where:        "serviceRedis.Get",
		Body:         "key: " + key,
		BodyResponse: value,
		Error:        err,
		IsError:      err != redis.Nil && err != nil,
	})
	if err == redis.Nil || value == "" {
		return &redisGrpc.ResponseGet{Message: "redis: key not found", IsCacheKeyNotFound: true}, nil
	}
	if err != redis.Nil && err != nil {
		return nil, err
	}

	return &redisGrpc.ResponseGet{Value: stringData.Val(), Message: "successful"}, nil
}

func (h *HandlerGrpc) Set(ctx context.Context, set *redisGrpc.RequestSet) (*redisGrpc.ResponseSet, error) {
	if err := validateHeaders(ctx); err != nil {
		return nil, err
	}

	key := set.GetKey()
	value := set.GetValue()
	expiration := set.GetExpirationInSeconds()

	timeInit := time.Now()

	err := h.redisCache.Set(ctx, key, value, time.Duration(expiration)*time.Second).Err()
	h.tracing.Logger(ctx, model.DataLogger{
		Operation:   "setCache",
		StartTime:   timeInit,
		EndTime:     time.Now(),
		ElapsedNano: time.Since(timeInit).Nanoseconds(),
		Where:       "serviceRedis.Set",
		Body:        fmt.Sprintf("key: %s, value: %s, expiration: %d", key, value, expiration),
		Error:       err,
		IsError:     err != redis.Nil && err != nil,
	})
	if err != nil {
		return nil, err
	}

	return &redisGrpc.ResponseSet{Message: "successful"}, nil
}

func (h *HandlerGrpc) Del(ctx context.Context, del *redisGrpc.RequestDel) (*redisGrpc.ResponseDel, error) {
	if err := validateHeaders(ctx); err != nil {
		return nil, err
	}

	key := del.GetKey()

	timeInit := time.Now()

	stringData := h.redisCache.Del(ctx, key)

	_, err := stringData.Result()
	h.tracing.Logger(ctx, model.DataLogger{
		Operation:    "delCache",
		StartTime:    timeInit,
		EndTime:      time.Now(),
		ElapsedNano:  time.Since(timeInit).Nanoseconds(),
		Where:        "serviceRedis.Del",
		Body:         "key: " + key,
		BodyResponse: "",
		Error:        err,
		IsError:      err != redis.Nil && err != nil,
	})
	if err != nil {
		return nil, err
	}

	return &redisGrpc.ResponseDel{Message: "successful"}, nil
}

func (h *HandlerGrpc) Expire(ctx context.Context, expire *redisGrpc.RequestExpire) (*redisGrpc.ResponseExpire, error) {
	if err := validateHeaders(ctx); err != nil {
		return nil, err
	}

	key := expire.GetKey()
	expiration := expire.GetExpirationInSeconds()

	timeInit := time.Now()

	stringData := h.redisCache.Expire(ctx, key, time.Duration(expiration)*time.Second)

	_, err := stringData.Result()
	h.tracing.Logger(ctx, model.DataLogger{
		Operation:    "delCache",
		StartTime:    timeInit,
		EndTime:      time.Now(),
		ElapsedNano:  time.Since(timeInit).Nanoseconds(),
		Where:        "serviceRedis.Expire",
		Body:         "key: " + key,
		BodyResponse: "",
		Error:        err,
		IsError:      err != redis.Nil && err != nil,
	})
	if err != nil {
		return nil, err
	}

	return &redisGrpc.ResponseExpire{Message: "successful"}, nil
}

func (h *HandlerGrpc) HGet(ctx context.Context, get *redisGrpc.RequestHGet) (*redisGrpc.ResponseHGet, error) {
	if err := validateHeaders(ctx); err != nil {
		return nil, err
	}

	hash := get.GetHash()
	key := get.GetKey()

	timeInit := time.Now()

	stringData := h.redisCache.HGet(ctx, hash, key)

	value, err := stringData.Result()
	h.tracing.Logger(ctx, model.DataLogger{
		Operation:    "hGetCache",
		StartTime:    timeInit,
		EndTime:      time.Now(),
		ElapsedNano:  time.Since(timeInit).Nanoseconds(),
		Where:        "serviceRedis.HGet",
		Body:         "hash: " + hash + "key: " + key,
		BodyResponse: value,
		Error:        err,
		IsError:      err != redis.Nil && err != nil,
	})
	if err == redis.Nil || value == "" {
		return &redisGrpc.ResponseHGet{Message: "redis: key not found", IsCacheKeyNotFound: true}, nil
	}
	if err != redis.Nil && err != nil {
		return nil, err
	}

	return &redisGrpc.ResponseHGet{Value: stringData.Val(), Message: "successful"}, nil
}

func (h *HandlerGrpc) HSet(ctx context.Context, set *redisGrpc.RequestHSet) (*redisGrpc.ResponseHSet, error) {
	if err := validateHeaders(ctx); err != nil {
		return nil, err
	}

	hash := set.GetHash()
	key := set.GetKey()
	value := set.GetValue()

	timeInit := time.Now()

	err := h.redisCache.HSet(ctx, hash, key, value).Err()
	h.tracing.Logger(ctx, model.DataLogger{
		Operation:   "hSetCache",
		StartTime:   timeInit,
		EndTime:     time.Now(),
		ElapsedNano: time.Since(timeInit).Nanoseconds(),
		Where:       "serviceRedis.HSet",
		Body:        fmt.Sprintf("hash: %s, key: %s, value: %s", hash, key, value),
		Error:       err,
		IsError:     err != redis.Nil && err != nil,
	})
	if err != nil {
		return nil, err
	}

	return &redisGrpc.ResponseHSet{Message: "successful"}, nil
}

func (h *HandlerGrpc) HDel(ctx context.Context, set *redisGrpc.RequestHDel) (*redisGrpc.ResponseHDel, error) {
	if err := validateHeaders(ctx); err != nil {
		return nil, err
	}

	hash := set.GetHash()
	key := set.GetKey()

	timeInit := time.Now()

	err := h.redisCache.HDel(ctx, hash, key).Err()
	h.tracing.Logger(ctx, model.DataLogger{
		Operation:   "hDelCache",
		StartTime:   timeInit,
		EndTime:     time.Now(),
		ElapsedNano: time.Since(timeInit).Nanoseconds(),
		Where:       "serviceRedis.HDel",
		Body:        fmt.Sprintf("hash: %s, key: %s", hash, key),
		Error:       err,
		IsError:     err != redis.Nil && err != nil,
	})
	if err != nil {
		return nil, err
	}

	return &redisGrpc.ResponseHDel{Message: "successful"}, nil
}

func (h *HandlerGrpc) HealthRedis(ctx context.Context, _ *redisGrpc.RequestHealthRedis) (*redisGrpc.ResponseHealthRedis, error) {
	_, err := h.redisCache.Health(ctx).Result()
	if err != redis.Nil && err != nil {
		return nil, err
	}

	return &redisGrpc.ResponseHealthRedis{Live: true}, nil
}

func validateHeaders(ctx context.Context) error {
	dataMetadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.DataLoss, "failed to get metadata")
	}

	requestId := dataMetadata["request-id"]
	if len(requestId) == 0 {
		return status.Errorf(codes.InvalidArgument, "missing 'request-id' header")
	}
	if strings.Trim(requestId[0], " ") == "" {
		return status.Errorf(codes.InvalidArgument, "empty 'request-id' header")
	}

	hostname := dataMetadata["hostname"]
	if len(hostname) == 0 {
		return status.Errorf(codes.InvalidArgument, "missing 'hostname' header")
	}
	if strings.Trim(hostname[0], " ") == "" {
		return status.Errorf(codes.InvalidArgument, "empty 'hostname' header")
	}

	return nil
}
