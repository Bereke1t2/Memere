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

// uploadBodyRoutes are the multipart upload routes that legitimately stream a
// large binary body THROUGH the API to the object store (lesson PDFs, direct
// video uploads, course thumbnails). This matters for the Google Drive backend,
// where the bytes are proxied by the backend rather than PUT straight to S3.
// These routes are exempt from the global JSON body limit; each handler enforces
// its own, content-appropriate size cap. Keys are gin route patterns as returned
// by c.FullPath().
var uploadBodyRoutes = map[string]struct{}{
	"/api/v1/lessons/:id/pdf":           {},
	"/api/v1/lessons/:id/videos/upload": {},
	"/api/v1/courses/:id/thumbnail":     {},
}

// BodyLimit rejects requests whose Content-Length exceeds maxBytes with 413.
// It also wraps the body in http.MaxBytesReader so oversized chunked transfers
// are caught during JSON binding, not just via the header.
//
// The multipart upload routes in uploadBodyRoutes are exempt: they proxy large
// binary bodies through the API on purpose and cap size themselves. (Pre-signed
// S3 upload routes never reach the API with a body at all.)
func BodyLimit(maxBytes int64) gin.HandlerFunc {
	msg := fmt.Sprintf("request body must not exceed %d bytes", maxBytes)
	return func(c *gin.Context) {
		if _, exempt := uploadBodyRoutes[c.FullPath()]; exempt {
			c.Next()
			return
		}
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
