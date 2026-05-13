// Package logger provides a shared Zap-based structured logger factory
// for all services in the orion-v2 monorepo.
package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds the parameters needed to build a logger.
// Each service maps its own config struct to this type in main.go,
// keeping this package free of any service-specific dependencies.
type Config struct {
	// Env controls the encoder preset:
	//   "local" | "development" → colorized console, caller info
	//   "staging" | "production" → JSON, RFC-3339 timestamps
	Env string

	// Level is the minimum log level: debug | info | warn | error
	Level string

	// Service is injected as the "service" field in every log entry.
	Service string
}

// New builds and returns a *zap.Logger configured for the given environment.
// The caller is responsible for calling logger.Sync() on shutdown.
func New(cfg Config) (*zap.Logger, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	var zapCfg zap.Config

	switch cfg.Env {
	case "local", "development":
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(level)
	default:
		zapCfg = zap.NewProductionConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(level)
		zapCfg.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	}

	zapCfg.InitialFields = map[string]interface{}{
		"service": cfg.Service,
		"env":     cfg.Env,
	}

	log, err := zapCfg.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, fmt.Errorf("logger: build failed: %w", err)
	}

	return log, nil
}

func parseLevel(s string) (zapcore.Level, error) {
	var l zapcore.Level
	if err := l.UnmarshalText([]byte(s)); err != nil {
		return l, fmt.Errorf("logger: unknown level %q (debug|info|warn|error)", s)
	}
	return l, nil
}
