package cache

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/Bereke1t2/Memere/memere-backend/config"
	redis "github.com/redis/go-redis/v9"
)

// Connect connects to Redis and verifies the connection. Pool size and timeouts
// are driven by config (Phase 6 Skill 3 §12.2).
//
// Two connection styles are supported:
//
//   - REDIS_URL set: a full redis:// or rediss:// connection string (as handed
//     out by managed providers like Upstash). A rediss:// scheme enables TLS.
//     Host/Port/Password/DB/REDIS_TLS are ignored; our configured pool size and
//     timeouts are still applied on top of the parsed options.
//   - discrete vars: Host/Port/Password/DB, with TLS enabled when REDIS_TLS=true
//     (required by Upstash and most managed Redis; leave false for local Redis).
func Connect(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	opts, err := redisOptions(cfg)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return client, nil
}

// redisOptions builds the go-redis client options from config, preferring a full
// REDIS_URL when present and otherwise assembling discrete settings.
func redisOptions(cfg *config.Config) (*redis.Options, error) {
	if cfg.Redis.URL != "" {
		opts, err := redis.ParseURL(cfg.Redis.URL)
		if err != nil {
			return nil, fmt.Errorf("parse REDIS_URL: %w", err)
		}
		// Our env-configured pool/timeout knobs win over any URL query params so
		// tuning stays in one place. TLS/auth/addr come from the URL itself.
		opts.PoolSize = cfg.Redis.PoolSize
		opts.DialTimeout = cfg.Redis.DialTimeout
		opts.ReadTimeout = cfg.Redis.ReadTimeout
		opts.WriteTimeout = cfg.Redis.WriteTimeout
		return opts, nil
	}

	opts := &redis.Options{
		Addr:         cfg.Redis.Addr(),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		DialTimeout:  cfg.Redis.DialTimeout,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
	}
	if cfg.Redis.TLS {
		// Managed Redis (Upstash) presents a valid public certificate, so no
		// custom CA is needed; verify against the configured host name.
		opts.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: cfg.Redis.Host,
		}
	}
	return opts, nil
}
