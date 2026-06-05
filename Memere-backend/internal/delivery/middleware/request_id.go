package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// requestIDHeader is the canonical header used to carry the request id in and
// out of the service.
const requestIDHeader = "X-Request-ID"

// RequestID ensures every request has a correlation id: it honors an incoming
// X-Request-ID, otherwise generates one, stores it on the context, and echoes
// it in the response header so clients and logs can correlate a request.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(requestIDHeader)
		if id == "" {
			id = uuid.NewString()
		}
		c.Set(requestIDKey, id)
		c.Writer.Header().Set(requestIDHeader, id)
		c.Next()
	}
}
