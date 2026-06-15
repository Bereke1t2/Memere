package cache

import (
	"context"
	"fmt"

	"github.com/Bereke1t2/Memere/memere-backend/config"
	redis "github.com/redis/go-redis/v9"
)

// Connect connects to Redis and verifies the connection. Pool size and timeouts
// are driven by config (Phase 6 Skill 3 §12.2).
func Connect(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr(),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		DialTimeout:  cfg.Redis.DialTimeout,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return client, nil
}
