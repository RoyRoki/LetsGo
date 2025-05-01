package main

import (
	"github.com/royroki/matching-service/internal/config"
	"github.com/royroki/matching-service/internal/infrastructure/redis"
)

func main() {
	// Load Config
	config.LoadConfig()

	// Initialize Redis
	redis.InitRedis()
}
