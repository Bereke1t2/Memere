package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Bereke1t2/Memere/memere-backend/config"
	"github.com/jackc/pgx/v5"
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
	poolCfg.MinConns = cfg.DB.MinConns
	poolCfg.MaxConnIdleTime = cfg.DB.MaxConnIdleTime
	poolCfg.MaxConnLifetime = cfg.DB.MaxConnLifetime
	poolCfg.HealthCheckPeriod = 1 * time.Minute

	// When connecting through a transaction-mode pooler (Neon's -pooler endpoint /
	// PgBouncer), the server multiplexes many client sessions over few backend
	// connections and cannot keep pgx's implicit prepared statements. Force the
	// simple query protocol so queries don't fail with "prepared statement does
	// not exist". Auto-detected when "-pooler." is in the DSN or DB_POOLED=true.
	if cfg.DB.Pooled || strings.Contains(dsn, "-pooler.") {
		poolCfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	}

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
