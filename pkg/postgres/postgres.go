// Package postgres provides a shared pgxpool connection factory
// for all services in the orion-v2 monorepo.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const (
	defaultMaxConns          int32         = 10
	defaultMinConns          int32         = 2
	defaultMaxConnLifetime   time.Duration = time.Hour
	defaultMaxConnIdleTime   time.Duration = 30 * time.Minute
	defaultHealthCheckPeriod time.Duration = time.Minute
)

// Config holds the parameters needed to open a connection pool.
// Extend with MaxConns / MinConns overrides as services need them.
type Config struct {
	DSN string
}

// Connect creates, configures, and validates a pgxpool connection pool.
// The caller is responsible for calling pool.Close() on shutdown.
func Connect(ctx context.Context, cfg Config, log *zap.Logger) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to parse DSN: %w", err)
	}

	poolCfg.MaxConns = defaultMaxConns
	poolCfg.MinConns = defaultMinConns
	poolCfg.MaxConnLifetime = defaultMaxConnLifetime
	poolCfg.MaxConnIdleTime = defaultMaxConnIdleTime
	poolCfg.HealthCheckPeriod = defaultHealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres: ping failed: %w", err)
	}

	log.Info("postgres connected",
		zap.String("host", poolCfg.ConnConfig.Host),
		zap.Uint16("port", poolCfg.ConnConfig.Port),
		zap.String("database", poolCfg.ConnConfig.Database),
	)

	return pool, nil
}
