// Package logger provides a common logging interface for the application.
package logger

import "context"

// Logger defines the contract for a logging system.
type Logger interface {
	Info(ctx context.Context, msg string, fields map[string]interface{})
	Warn(ctx context.Context, msg string, fields map[string]interface{})
	Error(ctx context.Context, msg string, err error, fields map[string]interface{})
	Debug(ctx context.Context, msg string, fields map[string]interface{})
}
