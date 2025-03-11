package config

import "github.com/go-redis/redis/v8"

// RedisConfigInterface extends DBConfig with Redis-specific methods.
type RedisConfigInterface interface {
	DBConfig
	GetDBIndex() int
	NewClient() *redis.Client
}
