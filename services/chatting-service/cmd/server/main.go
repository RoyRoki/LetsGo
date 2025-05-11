package main

import (
	"github.com/royroki/services/chatting-service/internal/infrastructure/redis"
)

func main() {
	// Initialize Redis
	redis.InitRedis()

}
