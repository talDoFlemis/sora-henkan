package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/taldoflemis/sora-henkan/settings"
)

// NewPool creates a new PostgreSQL connection pool with the provided settings
func NewPool(ctx context.Context, cfg settings.DatabaseSettings) (*pgxpool.Pool, error) {
	// Build connection string from settings
	connString := cfg.BuildConnectionString()

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	config.ConnConfig.Tracer = otelpgx.NewTracer()

	// Configure connection pool
	config.MaxConns = int32(cfg.MaxOpenConns)
	config.MinConns = int32(cfg.MaxIdleConns)
	config.MaxConnLifetime = time.Duration(cfg.ConnMaxLifetimeMinutes) * time.Minute
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	// Create pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	if err := otelpgx.RecordStats(pool); err != nil {
		return nil, fmt.Errorf("unable to record database stats: %w", err)
	}

	return pool, nil
}
