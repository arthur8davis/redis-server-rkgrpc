package handlergrpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"rhyme80/redis-server-rkgrpc/infrastructure/cache"
	redisGrpc "rhyme80/redis-server-rkgrpc/infrastructure/servicegrpc/protosw"
)

type HandlerGrpc struct {
	redisCache cache.Cache

	redisGrpc.UnimplementedRedisServiceServer
}

func New(redisCache cache.Cache) *HandlerGrpc {
	return &HandlerGrpc{redisCache: redisCache}
}

func (h *HandlerGrpc) Get(ctx context.Context, get *redisGrpc.RequestGet) (*redisGrpc.ResponseGet, error) {
	//if err := validateHeaders(ctx); err != nil {
	//	return nil, err
	//}

	key := get.GetKey()

	stringData := h.redisCache.Get(ctx, key)

	value, err := stringData.Result()
	if err == redis.Nil || value == "" {
		fmt.Println("--------------------")
		fmt.Println(err)
		return &redisGrpc.ResponseGet{Message: "redis: key not found", IsCacheKeyNotFound: true}, nil
	}
	if err != redis.Nil && err != nil {
		fmt.Println("###################")
		fmt.Println(err)
		return nil, err
	}

	return &redisGrpc.ResponseGet{Value: stringData.Val(), Message: "successful"}, nil
}

func (h *HandlerGrpc) Set(ctx context.Context, set *redisGrpc.RequestSet) (*redisGrpc.ResponseSet, error) {
	//if err := validateHeaders(ctx); err != nil {
	//	return nil, err
	//}

	key := set.GetKey()
	value := set.GetValue()
	expiration := set.GetExpirationInSeconds()

	err := h.redisCache.Set(ctx, key, value, time.Duration(expiration)*time.Second).Err()
	if err != nil {
		return nil, err
	}

	return &redisGrpc.ResponseSet{Message: "successful"}, nil
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
