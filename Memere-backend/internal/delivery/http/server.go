package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	redis "github.com/redis/go-redis/v9"

	"github.com/Bereke1t2/Memere/memere-backend/internal/version"
)

// healthHandler reports liveness of the dependencies. Mounted outside /api/v1 so
// orchestrators can probe it without versioning concerns.
func healthHandler(db *pgxpool.Pool, cache *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		errDB := db.Ping(c.Request.Context())
		errRedis := cache.Ping(c.Request.Context()).Err()

		if errDB != nil || errRedis != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "error",
				"db":     errDB == nil,
				"redis":  errRedis == nil,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
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
