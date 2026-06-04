package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Bereke1t2/Memere/memere-backend/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect connects to PostgreSQL using pgxpool and verifies the connection.
func Connect(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := cfg.DB.DSN()

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	poolCfg.MaxConns = cfg.DB.MaxConns
	poolCfg.MaxConnIdleTime = cfg.DB.MaxConnIdleTime
	poolCfg.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
