package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS applies a configurable cross-origin policy.
//
//   - In development (allowAll = true, i.e. allowedOrigins == ["*"]) any origin
//     is permitted — convenient for local Flutter/web dev.
//   - In production the allowedOrigins list must contain explicit domains. If it
//     is still "*" the middleware logs a warning and closes the CORS response
//     (§2.5 fail-closed rule): cross-origin requests are rejected until proper
//     origins are configured.
//
// Preflight OPTIONS requests are answered with 204. Max-Age caches preflight
// results in the browser for 10 minutes to reduce OPTIONS round-trips.
func CORS(allowedOrigins []string, isProduction bool) gin.HandlerFunc {
	allowAll := len(allowedOrigins) == 1 && allowedOrigins[0] == "*"

	if allowAll && isProduction {
		slog.Warn("CORS_ALLOWED_ORIGINS is '*' in production — cross-origin requests will be rejected until explicit origins are configured")
		allowAll = false // fail closed
	}

	allowed := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowed[o] = true
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		switch {
		case allowAll:
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		case origin != "" && allowed[origin]:
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Add("Vary", "Origin")
		default:
			// Origin not allowed — no ACAO header; browser will block the request.
			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")
		c.Writer.Header().Set("Access-Control-Max-Age", "600") // 10 minutes

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
