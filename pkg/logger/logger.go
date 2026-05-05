package logger

import (
	"gopher-order-service/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	var zapCfg zap.Config
	if cfg.App.Env == "production" {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Set log level from config
	level, err := zapcore.ParseLevel(cfg.Logger.Level)
	if err == nil {
		zapCfg.Level = zap.NewAtomicLevelAt(level)
	}

	logger, err := zapCfg.Build()
	if err != nil {
		return nil, err
	}

	zap.ReplaceGlobals(logger) // Set as global logger if needed somewhere implicitly
	return logger, nil
}
