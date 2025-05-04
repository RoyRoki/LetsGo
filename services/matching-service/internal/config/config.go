package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/royroki/letsgo/services/matching-service/internal/constants"
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
		cfg = &Config{
			RedisHost:     getEnv(constants.RedisHost, "redis"),
			RedisPort:     getEnv(constants.RedisPort, "6379"),
			RedisPassword: getEnv(constants.RedisPass, "redis@123"),
			GRPCPort:      getEnv(constants.GRPCPORT, "9090"),
			MatchTimeout:  getEnvAsInt("MATCH_TIMEOUT_SECONDS", 30),
		}
	})
	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		log.Printf("ENV: %s = %s\n", key, value)
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
