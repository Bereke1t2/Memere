package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	redis "github.com/redis/go-redis/v9"

	"github.com/Bereke1t2/Memere/memere-backend/internal/version"
)

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
