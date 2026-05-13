package main

import (
	"context"
	"log"

	"go.uber.org/zap"

	"github.com/fekuna/orion-v2/pkg/logger"
	"github.com/fekuna/orion-v2/pkg/postgres"
	"github.com/fekuna/orion-v2/services/product-service/internal/config"
	"github.com/fekuna/orion-v2/services/product-service/internal/product"
	"github.com/fekuna/orion-v2/services/product-service/internal/server"
)

func main() {
	// 1. Config — must be first; everything else depends on it.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2. Logger — init early so all subsequent messages are structured.
	//    The pkg/logger.Config decouples the shared factory from this
	//    service's internal config struct.
	zapLog, err := logger.New(logger.Config{
		Env:     cfg.App.Env,
		Level:   cfg.Log.Level,
		Service: cfg.App.Name,
	})
	if err != nil {
		log.Fatalf("failed to build logger: %v", err)
	}
	defer zapLog.Sync() //nolint:errcheck

	// 3. Infrastructure — shared external resources via pkg/.
	ctx := context.Background()

	db, err := postgres.Connect(ctx, postgres.Config{
		DSN: cfg.Database.DSN,
	}, zapLog)
	if err != nil {
		zapLog.Fatal("failed to connect to postgres", zap.Error(err))
	}
	defer db.Close()

	// 4. Modules — wire each feature: repo → usecase → handler.
	//    Add new modules here as the service grows.
	productRepo := product.NewPostgresRepo(db)
	productUC := product.NewUseCase(productRepo)
	productHandler := product.NewHandler(productUC, zapLog)

	// 5. Server — create and register all routes.
	srv := server.New(cfg, zapLog)
	v1 := srv.Group("/api/v1")
	productHandler.RegisterRoutes(v1.Group("/products"))

	// 6. Start — blocks until SIGINT/SIGTERM then graceful shutdown.
	if err := srv.Start(); err != nil {
		zapLog.Fatal("server exited with error", zap.Error(err))
	}
}
