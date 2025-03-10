package config

import (
	"log"
	"os"
	"strconv"

	config "github.com/royroki/LetsGo/config/database"
	"github.com/royroki/LetsGo/internal/config/constants"
)

type RedisConfigImpl struct{}

func NewRedisConfig() config.DBConfig {
	return &RedisConfigImpl{}
}

func (r *RedisConfigImpl) GetAddress() string {
	addr := os.Getenv(constants.RedisAddressEnv)
	if addr == "" {
		log.Fatal("REDIS_ADDRESS environment variable is not set")
	}
	return addr
}

func (r *RedisConfigImpl) GetPassword() string {
	return os.Getenv(constants.RedisPasswordEnv)
}

func (r *RedisConfigImpl) GetDBIndex() int {
	db := os.Getenv(constants.RedisDBEnv)
	if db == "" {
		return 0 // Default to database 0 if not set
	}
	dbIndex, err := strconv.Atoi(db)
	if err != nil {
		log.Fatalf("Error parsing REDIS_DB value: %v", err)
	}
	return dbIndex
}
