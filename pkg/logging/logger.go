package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new zap.Logger.
func NewLogger(serviceName string) (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	
	// Custom encoder config
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.CallerKey = "caller"
	
	// If in local dev, we might want console encoder, but for sandbox we stick to production-like JSON
	if os.Getenv("LOG_FORMAT") == "console" {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := config.Build(zap.AddCaller(), zap.Fields(
		zap.String("service", serviceName),
		zap.String("env", os.Getenv("ENV")),
	))
	if err != nil {
		return nil, err
	}

	return logger, nil
}
