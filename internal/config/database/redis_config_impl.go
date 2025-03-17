package database

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/royroki/LetsGo/internal/config"
	"github.com/royroki/LetsGo/internal/config/constants"
)

type RedisConfigImpl struct{}

// NewRedisConfig initializes a Redis configuration instance.
func NewRedisConfig() config.RedisConfigInterface {
	return &RedisConfigImpl{}
}

// GetDBIndex implements config.RedisConfigInterface.
func (r *RedisConfigImpl) GetDBIndex() int {
	indexStr := os.Getenv(constants.RedisDBEnv)
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		log.Fatalf("Invalid REDIS_INDEX: %v", err)
	}
	return index
}

// GetDBAddress retrieves the Redis server address.
func (r *RedisConfigImpl) GetDBAddress() string {
	addr := os.Getenv(constants.RedisAddressEnv)
	if addr == "" {
		log.Fatal("REDIS_ADDRESS is not set")
	}
	return addr
}

// GetDBPassword retrieves the Redis authentication password.
func (r *RedisConfigImpl) GetDBPassword() string {
	return os.Getenv(constants.RedisPasswordEnv)
}

// GetDBPort retrieves the Redis server port.
func (r *RedisConfigImpl) GetDBPort() int {
	portStr := os.Getenv(constants.RedisPortEnv)
	if portStr == "" {
		portStr = constants.RedisDefPortStr // Default Redis port
		log.Println("REDIS_PORT is not set, using default 6379")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid REDIS_PORT: %v", err)
	}

	return port
}

// NewClient initializes and returns a Redis client.
func (r *RedisConfigImpl) NewClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     r.GetDBAddress(),
		Password: r.GetDBPassword(), // Ensure the password is passed
		DB:       r.GetDBIndex(),    // Make sure to use the correct DB index
	})

	return client
}

// Ping tests the connection to the Redis server.
func (r *RedisConfigImpl) Ping() error {
	client := r.NewClient()
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return err
	}

	log.Println("Successfully connected to Redis")
	return nil
}
