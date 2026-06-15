package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/metrics"
)

// Metrics records Prometheus RED metrics (Rate / Errors / Duration) for every
// HTTP request. It must be placed AFTER RequestID so the route template is
// already settled on the context, and BEFORE handlers so latency captures the
// full handler time.
//
// The "route" label uses gin's FullPath() (e.g. "/api/v1/courses/:id") rather
// than the actual path to keep series cardinality bounded.
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unknown"
		}
		status := statusClass(c.Writer.Status())

		metrics.HTTPRequests.WithLabelValues(c.Request.Method, route, status).Inc()
		metrics.HTTPDuration.WithLabelValues(c.Request.Method, route).Observe(time.Since(start).Seconds())
	}
}

// statusClass converts a numeric HTTP status code to the conventional string
// class ("2xx", "4xx", "5xx") so label cardinality stays at five values.
func statusClass(code int) string {
	switch {
	case code < 200:
		return "1xx"
	case code < 300:
		return "2xx"
	case code < 400:
		return "3xx"
	case code < 500:
		return "4xx"
	default:
		return "5xx"
	}
}

// statusClassFromInt is an alias used in tests.
func statusClassFromInt(code int) string { return statusClass(code) }
