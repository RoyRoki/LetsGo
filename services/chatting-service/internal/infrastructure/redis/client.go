package redis

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/royroki/services/chatting-service/internal/constants"
)

var (
	Rdb redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	// Construct the Redis address from the environment variables
	redisAddr := os.Getenv(constants.RedisHost) + ":" + os.Getenv(constants.RedisPort)

	Rdb = *redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: os.Getenv(constants.RedisPass),
		DB:       0, // use default DB
	})

	// Ping to Redis
	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
	}

	log.Printf("✅ Redis connected: %s", redisAddr)
}
