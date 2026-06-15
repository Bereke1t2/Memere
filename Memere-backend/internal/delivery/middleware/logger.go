package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Bereke1t2/Memere/memere-backend/pkg/logger"
)

// Logger emits a structured slog access log line per request: method, path,
// status, latency, client IP, and request id. It deliberately logs none of the
// request body, the Authorization header, query strings, or any token (Non-
// Negotiable #8, spec §7.3). It also enriches the request context with a
// request-scoped logger pre-tagged with request_id and method so that usecase
// and repository code can call logger.FromContext(ctx) without extra setup.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// Path only (no raw query string, which could contain sensitive values).
		path := c.Request.URL.Path
		reqID := RequestIDFromContext(c)

		// Attach a request-scoped logger so downstream code can log with context.
		reqLogger := slog.Default().With(
			slog.String("request_id", reqID),
			slog.String("method", c.Request.Method),
		)
		c.Request = c.Request.WithContext(logger.WithContext(c.Request.Context(), reqLogger))

		c.Next()

		reqLogger.Info("request",
			slog.String("path", path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", time.Since(start)),
			slog.String("ip", c.ClientIP()),
		)
	}
}
