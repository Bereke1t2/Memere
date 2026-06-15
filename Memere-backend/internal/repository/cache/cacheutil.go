// Package cache provides cache-aside helpers for the Redis layer (§3.3, §12.2).
// It wraps hot read paths with a short-TTL Redis layer; misses fall through to
// the provided loader (typically a Postgres repository call). All cache
// operations fail open: a Redis error causes a transparent DB fallback so cache
// unavailability never degrades the API.
//
// Keys and TTLs:
//
//	courses:list:{hash}     60 s  — invalidated on any course mutation
//	courses:detail:{id}     60 s  — invalidated on update / publish / delete
//	users:id:{id}           5 m   — invalidated on profile update
//	sub:plans               10 m  — invalidated on plan config change (rare)
package cache

import (
	"context"
	"encoding/json"
	"time"

	goredis "github.com/redis/go-redis/v9"

	infmetrics "github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/metrics"
)

// GetOrSet is a generic cache-aside helper: it returns the cached JSON value if
// present; otherwise it calls load, caches the result, and returns it. Errors
// from Redis are logged and treated as misses so the DB always acts as fallback.
//
// class is the metric label ("course_list", "course_detail", "user_id", …).
func GetOrSet[T any](
	ctx context.Context,
	rdb *goredis.Client,
	key, class string,
	ttl time.Duration,
	load func(ctx context.Context) (T, error),
) (T, error) {
	if b, err := rdb.Get(ctx, key).Bytes(); err == nil {
		var v T
		if json.Unmarshal(b, &v) == nil {
			infmetrics.ObserveCacheHit(class)
			return v, nil
		}
	}

	infmetrics.ObserveCacheMiss(class)
	v, err := load(ctx)
	if err != nil {
		return v, err
	}
	if b, e := json.Marshal(v); e == nil {
		_ = rdb.Set(ctx, key, b, ttl).Err()
	}
	return v, nil
}

// Del removes one or more cache keys (write-through invalidation). Errors are
// silently dropped — a stale cache entry will expire on its own TTL.
func Del(ctx context.Context, rdb *goredis.Client, keys ...string) {
	if len(keys) == 0 {
		return
	}
	_ = rdb.Del(ctx, keys...).Err()
}
