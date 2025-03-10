package config

// DBConfig defines the contract for database configuration methods.
type DBConfig interface {
	GetAddress() string
	GetPassword() string
	GetDBIndex() int
}
