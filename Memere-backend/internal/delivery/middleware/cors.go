package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS applies a configurable cross-origin policy. allowedOrigins is the list
// of permitted Origin values; a single "*" entry allows any origin (suitable
// for development — production should pin explicit domains). A preflight
// OPTIONS request is answered with 204.
func CORS(allowedOrigins []string) gin.HandlerFunc {
	allowAll := len(allowedOrigins) == 1 && allowedOrigins[0] == "*"
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
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
