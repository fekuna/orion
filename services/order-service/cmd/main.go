package main

import (
	"log"

	"github.com/fekuna/orion/pkg/logger"
	"github.com/fekuna/orion/services/product-service/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	zapLog, err := logger.New(logger.Config{
		Env:     cfg.App.Env,
		Level:   cfg.Log.Level,
		Service: cfg.App.Name,
	})
	if err != nil {
		log.Fatalf("failed to build logger: %v", err)
	}
	defer zapLog.Sync()

}
