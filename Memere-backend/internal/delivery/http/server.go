package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	redis "github.com/redis/go-redis/v9"
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
