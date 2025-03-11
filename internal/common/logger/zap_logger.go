package logger

import (
	"context"
	"log"
	"os"

	"github.com/royroki/LetsGo/internal/config/constants"
	"go.uber.org/zap"
)

// ZapLogger is an implementation of Logger using Uber's Zap.
type ZapLogger struct {
	logger *zap.Logger
}

// NewZapLogger creates a new instance of ZapLogger dynamically based on APP_ENV.
func NewZapLogger() *ZapLogger {
	env := os.Getenv(constants.AppEnv)
	var zapLogger *zap.Logger
	var err error

	if env == constants.DevelopmentStr {
		zapLogger, err = zap.NewDevelopment()
	} else {
		zapLogger, err = zap.NewProduction()
	}

	if err != nil {
		log.Fatalf("Failed to initialize Zap logger: %v", err)
	}

	return &ZapLogger{logger: zapLogger}
}

// Info logs an informational message.
func (l *ZapLogger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	l.logger.Sugar().Infow(msg, "fields", fields)
}

// Warn logs a warning message.
func (l *ZapLogger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	l.logger.Sugar().Warnw(msg, "fields", fields)
}

// Error logs an error message.
func (l *ZapLogger) Error(ctx context.Context, msg string, err error, fields map[string]interface{}) {
	l.logger.Sugar().Errorw(msg, "error", err, "fields", fields)
}

// Debug logs a debug message.
func (l *ZapLogger) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	l.logger.Sugar().Debugw(msg, "fields", fields)
}
