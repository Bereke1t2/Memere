package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

// internalTokenHeader is the header Cloud Scheduler (or any internal caller)
// sets to authenticate to the /internal/jobs/* endpoints. A custom header keeps
// these machine-to-machine routes distinct from user bearer auth, so nothing on
// the normal auth path can accidentally satisfy them.
const internalTokenHeader = "X-Internal-Job-Token"

// InternalToken guards the internal job endpoints with a shared secret compared
// in constant time. The secret is injected via INTERNAL_JOB_TOKEN (never
// hardcoded) and the same value is configured on the Cloud Scheduler job's
// header. The caller may present it either in the X-Internal-Job-Token header or
// as "Authorization: Bearer <token>" (Cloud Scheduler's OIDC-free custom-header
// mode uses the former; the Bearer form is accepted for convenience).
//
// A mismatch or missing token returns an opaque 401 — the same envelope the rest
// of the API uses — so probing reveals nothing about the expected value.
func InternalToken(expected string) gin.HandlerFunc {
	expectedBytes := []byte(expected)
	return func(c *gin.Context) {
		got := c.GetHeader(internalTokenHeader)
		if got == "" {
			got = bearerToken(c)
		}
		// Constant-time compare; the length guard avoids leaking length via timing.
		if len(got) != len(expectedBytes) ||
			subtle.ConstantTimeCompare([]byte(got), expectedBytes) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "authentication required",
			})
			return
		}
		c.Next()
	}
}
