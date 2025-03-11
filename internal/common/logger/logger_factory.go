package logger

import (
	"os"

	"github.com/royroki/LetsGo/internal/config/constants"
)

// LoggerType defines available logging backends.
type LoggerType string

const (
	ZapLoggerType LoggerType = constants.ZapLoggerTypeStr
)

// NewLoggerFactory initializes a logger based on configuration.
func NewLoggerFactory() Logger {
	loggerType := os.Getenv(constants.LoggerTypeEnv)

	switch LoggerType(loggerType) {
	case ZapLoggerType:
		return NewZapLogger()
	default:
		return NewZapLogger() // Default to ZapLogger
	}
}
