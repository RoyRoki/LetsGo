package main

import (
	"github.com/royroki/letsgo/services/matching-service/internal/config"
	"github.com/royroki/letsgo/services/matching-service/internal/infrastructure/grpc"
	"github.com/royroki/letsgo/services/matching-service/internal/infrastructure/redis"
)

func main() {
	// Load Config
	config.LoadConfig()

	// Initialize Redis
	redis.InitRedis()

	// Init gRPC Server
	grpc.NewGRPCServer()
}
