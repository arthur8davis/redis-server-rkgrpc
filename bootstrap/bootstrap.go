package bootstrap

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/rookie-ninja/rk-entry/v2/entry"
	"github.com/rookie-ninja/rk-grpc/v2/boot"
	"google.golang.org/grpc"
	"rhyme80/redis-server-rkgrpc/infrastructure/cache"
	"rhyme80/redis-server-rkgrpc/infrastructure/handlergrpc"
	pb "rhyme80/redis-server-rkgrpc/infrastructure/servicegrpc/protosw"
)

func Run(boot []byte) {
	_ = godotenv.Load(".env")

	// Bootstrap basic entries form boot config
	rkentry.BootstrapBuiltInEntryFromYAML(boot)
	rkentry.BootstrapPluginEntryFromYAML(boot)

	// Bootstrap grpc from boot config
	res := rkgrpc.RegisterGrpcEntryYAML(boot)

	// Get gRPCEntry
	grpcEntry := res["redisserver"].(*rkgrpc.GrpcEntry)
	// Register gRPC server
	handler := handlergrpc.New(cache.New(newRedis()))
	grpcEntry.AddRegFuncGrpc(func(server *grpc.Server) {
		pb.RegisterRedisServiceServer(server, handler)
	})
	// Register grpc-gateway func
	grpcEntry.AddRegFuncGw(pb.RegisterRedisServiceHandlerFromEndpoint)

	// Bootstrap grpc entry
	grpcEntry.Bootstrap(context.Background())

	// Wait for shutdown signal
	rkentry.GlobalAppCtx.WaitForShutdownSig()

	// Interrupt gin entry
	grpcEntry.Interrupt(context.Background())

}
