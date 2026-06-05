package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger emits a structured access log line per request: method, path, status,
// latency, client IP and request id. It deliberately logs none of the request
// body, the Authorization header, query strings, or any token — those can carry
// passwords or credentials (Non-Negotiable #8, spec §7.3).
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// Path only (no raw query string, which could contain sensitive values).
		path := c.Request.URL.Path

		c.Next()

		log.Printf("method=%s path=%s status=%d latency=%s ip=%s request_id=%s",
			c.Request.Method,
			path,
			c.Writer.Status(),
			time.Since(start),
			c.ClientIP(),
			RequestIDFromContext(c),
		)
	}
}
