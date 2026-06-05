package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery recovers from a panic in any downstream handler, logs the stack
// server-side, and returns the standard 500 error envelope without leaking the
// panic value to the client (spec §7.3). It replaces gin's default Recovery so
// the response shape matches the rest of the API.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic recovered [request_id=%s]: %v\n%s", RequestIDFromContext(c), r, debug.Stack())
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    "INTERNAL_SERVER_ERROR",
					"message": "An internal server error occurred",
				})
			}
		}()
		c.Next()
	}
}
