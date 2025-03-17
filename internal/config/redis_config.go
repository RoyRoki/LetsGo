package config

import "github.com/redis/go-redis/v9"

// RedisConfigInterface extends DBConfig with Redis-specific methods.
type RedisConfigInterface interface {
	DBConfig
	GetDBIndex() int
	NewClient() *redis.Client
}
