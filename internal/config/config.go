// config/config.go - Config interface for loading and accessing configuration values

package config

// Config defines the contract for configuration-related methods.
type Config interface {
	LoadEnv() error          // Loads environment variables (e.g., from .env file)
	Get(key string) string   // Retrieves a string value for the given key
	GetInt(key string) int   // Retrieves an integer value for the given key
	GetBool(key string) bool // Retrieves a boolean value for the given key
}

// DBConfig defines the contract for database configuration.
type DBConfig interface {
	GetDBAddress() string
	GetDBPassword() string
	GetDBPort() int
	Ping() error
}
