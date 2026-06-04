package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	redis "github.com/redis/go-redis/v9"
)

func NewServer(db *pgxpool.Pool, cache *redis.Client) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		log.Printf("[%s] %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	router.GET("/health", func(c *gin.Context) {
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

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	return router
}
