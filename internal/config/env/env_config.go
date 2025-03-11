// config/env_config.go - EnvConfig implementation for loading environment variables

package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/royroki/LetsGo/internal/config"
	"github.com/royroki/LetsGo/internal/config/constants"
)

// EnvConfig holds methods for interacting with environment variables.
type EnvConfig struct {
	RedisAddress  string
	RedisPassword string
	RedisDB       int
	RedisPort     int
	LoggerType    string
}

// NewEnvConfig creates a new EnvConfig instance and loads environment variables.
func NewEnvConfig() config.Config {
	config := &EnvConfig{}
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}
	return config
}

// LoadEnv loads environment variables from a .env file.
func (c *EnvConfig) LoadEnv() error {
	// Get the environment (development, production, test)
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development" // Default to development
	}

	var envFile string
	if env == "development" {
		envFile = ".env.development"
	} else if env == "production" {
		envFile = ".env.production"
	} else if env == "test" {
		envFile = ".env.test"
	} else {
		envFile = ".env" // Default file
	}

	// Load the correct environment file
	if err := godotenv.Load(envFile); err != nil {
		return err
	}

	// Load Redis configurations
	c.RedisAddress = os.Getenv(constants.RedisAddressEnv)
	c.RedisPassword = os.Getenv(constants.RedisPasswordEnv)
	if c.RedisPassword == "" {
		log.Println("No Redis password set, connecting without password")
	}
	c.RedisDB = c.GetInt(constants.RedisDBEnv)
	c.RedisPort = c.GetInt(constants.RedisPortEnv)

	// Load Logger Type
	c.LoggerType = c.Get(constants.LoggerTypeEnv)

	return nil
}

// Get retrieves the value of the environment variable by key.
func (e *EnvConfig) Get(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s not set", key)
	}
	return value
}

// GetBool retrieves a boolean value from the environment variable.
func (e *EnvConfig) GetBool(key string) bool {
	return e.Get(key) == "true"
}

// GetInt retrieves an integer value from the environment variable.
func (e *EnvConfig) GetInt(key string) int {
	value := e.Get(key)
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Error converting %s to integer: %v", key, err)
	}
	return intValue
}
