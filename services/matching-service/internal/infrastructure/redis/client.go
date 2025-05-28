package redis

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/royroki/letsgo/services/matching-service/internal/constants"
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

func GetOrCreateIndex(userID string) (int64, error) {
	indexStr, err := Rdb.HGet(Ctx, "user:index", userID).Result()
	if err == redis.Nil {
		// Doesn't exist, create new index
		index, err := Rdb.Incr(Ctx, "user:index:counter").Result()
		if err != nil {
			return 0, err
		}

		pipe := Rdb.TxPipeline()
		pipe.HSet(Ctx, "user:index", userID, index)
		pipe.HSet(Ctx, "index:user", index, userID)
		_, err = pipe.Exec(Ctx)
		if err != nil {
			return 0, err
		}

		return index, nil
	} else if err != nil {
		return 0, err
	}

	// Convert index string to int64
	index, err := strconv.ParseInt(indexStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return index, nil
}
