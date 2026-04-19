package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a new zap logger
func New(level string, format string) (*zap.Logger, error) {
	config := zap.NewProductionConfig()

	// Set log level
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zap.DebugLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	default:
		zapLevel = zap.InfoLevel
	}
	config.Level = zap.NewAtomicLevelAt(zapLevel)

	// Set output format
	if format == "console" {
		config.Encoding = "console"
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config.Encoding = "json"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// Add initial fields
	config.InitialFields = map[string]interface{}{
		"service": "url-shortener",
	}

	return config.Build()
}

// NewNop returns a no-op logger
func NewNop() *zap.Logger {
	return zap.NewNop()
}
