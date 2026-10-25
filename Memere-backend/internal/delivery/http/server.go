package http

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	redis "github.com/redis/go-redis/v9"

	"github.com/Bereke1t2/Memere/memere-backend/internal/version"
)

var startTime = time.Now()

// healthHandler is the liveness probe (/healthz). Returns 200 if the process
// is alive regardless of dependency state. Kubernetes restarts the pod only
// when this returns non-200 or times out — i.e., the process is hung/dead.
func healthHandler(_ *pgxpool.Pool, _ *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// readyzHandler is the readiness probe (/readyz). Returns 200 only when the
// process can serve traffic: both PostgreSQL and Redis must be reachable.
// Kubernetes stops routing traffic to the pod while this returns non-200.
func readyzHandler(db *pgxpool.Pool, cache *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		errDB := db.Ping(c.Request.Context())
		errRedis := cache.Ping(c.Request.Context()).Err()

		if errDB != nil || errRedis != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unavailable",
				"db":     errDB == nil,
				"redis":  errRedis == nil,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	}
}

// systemStatsHandler reports internal runtime statistics (Go memory, goroutines,
// GC stats, DB connection pool, and Redis pool). It provides full observability
// during load tests without requiring paid APM or cloud metrics dashboards.
func systemStatsHandler(db *pgxpool.Pool, cache *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		stats := gin.H{
			"status": "ok",
			"uptime": time.Since(startTime).Round(time.Second).String(),
			"runtime": gin.H{
				"goroutines":     runtime.NumGoroutine(),
				"cpus":           runtime.NumCPU(),
				"alloc_mb":       float64(m.Alloc) / (1024 * 1024),
				"total_alloc_mb": float64(m.TotalAlloc) / (1024 * 1024),
				"sys_mb":         float64(m.Sys) / (1024 * 1024),
				"heap_inuse_mb":  float64(m.HeapInuse) / (1024 * 1024),
				"num_gc":         m.NumGC,
			},
		}

		if db != nil {
			dbStat := db.Stat()
			stats["db_pool"] = gin.H{
				"acquired_conns": dbStat.AcquiredConns(),
				"idle_conns":     dbStat.IdleConns(),
				"total_conns":    dbStat.TotalConns(),
				"max_conns":      dbStat.MaxConns(),
			}
		}

		if cache != nil {
			cacheStat := cache.PoolStats()
			stats["redis_pool"] = gin.H{
				"total_conns": cacheStat.TotalConns,
				"idle_conns":  cacheStat.IdleConns,
				"stale_conns": cacheStat.StaleConns,
				"hits":        cacheStat.Hits,
				"misses":      cacheStat.Misses,
				"timeouts":    cacheStat.Timeouts,
			}
		}

		c.JSON(http.StatusOK, stats)
	}
}

// versionHandler reports the build version, commit SHA, and build time injected
// via ldflags at image build time. Useful for verifying which release is running.
func versionHandler() gin.HandlerFunc {
	body := gin.H{
		"version":    version.Version,
		"commit":     version.Commit,
		"build_time": version.BuildTime,
	}
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, body)
	}
}
