package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds the hardening response headers required by §7.3 and
// §1.5. All routes receive these headers; the CSP is strict because the API
// only serves JSON (no inline scripts, no frames, no embedded assets).
//
// HSTS assumes a TLS-terminating edge (spec §3.2 / Non-Negotiable #6); in
// local dev the header is present but harmless over HTTP.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		h.Set("X-XSS-Protection", "0") // modern browsers: CSP is the mitigation; legacy header causes issues
		c.Next()
	}
}

// BodyLimit rejects requests whose Content-Length exceeds maxBytes with 413.
// It also wraps the body in http.MaxBytesReader so oversized chunked transfers
// are caught during JSON binding, not just via the header.
// Pre-signed upload routes bypass this via the S3 redirect — they never proxy
// the binary body through the API server.
func BodyLimit(maxBytes int64) gin.HandlerFunc {
	msg := fmt.Sprintf("request body must not exceed %d bytes", maxBytes)
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxBytes {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"code":    "REQUEST_TOO_LARGE",
				"message": msg,
			})
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()

		// MaxBytesReader sets a flag on the response writer when the limit is
		// exceeded during reading (e.g. chunked transfer). Catch it here.
		if c.Writer.Written() {
			return
		}
		if err := c.Errors.Last(); err != nil {
			if isBodyTooLarge(err.Err) {
				c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
					"code":    "REQUEST_TOO_LARGE",
					"message": msg,
				})
			}
		}
	}
}

func isBodyTooLarge(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "http: request body too large"
}
