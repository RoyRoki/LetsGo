package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/royroki/matching-service/internal/constants"
)

type Config struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
	GRPCPort      string
	MatchTimeout  int
}

var (
	cfg  *Config
	once sync.Once
)

func LoadConfig() *Config {
	once.Do(func() {
		// Load .env file
		// if err := godotenv.Load(); err != nil {
		// 	log.Println("⚠️ No .env file found, using system environment variables")
		// }
		cfg = &Config{
			RedisHost:     getEnv(constants.RedisHost, "redis"),
			RedisPort:     getEnv(constants.RedisPort, "6379"),
			RedisPassword: getEnv("REDIS_PASSWORD", ""),
			GRPCPort:      getEnv("GRPC_PORT", "50051"),
			MatchTimeout:  getEnvAsInt("MATCH_TIMEOUT_SECONDS", 30),
		}
	})
	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, defaultVal int) int {
	if valueStr, ok := os.LookupEnv(key); ok {
		var value int
		_, err := fmt.Sscanf(valueStr, "%d", &value)
		if err != nil {
			return value
		}
	}
	return defaultVal
}
