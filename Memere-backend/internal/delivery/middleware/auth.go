package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/course"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/jwt"
)

// bearerToken extracts the token from an "Authorization: Bearer <token>" header,
// returning "" when the header is absent or malformed.
func bearerToken(c *gin.Context) string {
	h := c.GetHeader("Authorization")
	const prefix = "Bearer "
	if len(h) > len(prefix) && strings.EqualFold(h[:len(prefix)], prefix) {
		return strings.TrimSpace(h[len(prefix):])
	}
	return ""
}

// actorFromToken verifies an access token and builds the caller's identity.
func actorFromToken(mgr *jwt.Manager, token string) (*course.Actor, error) {
	claims, err := mgr.Verify(token, jwt.TokenTypeAccess)
	if err != nil {
		return nil, err
	}
	return &course.Actor{UserID: claims.UserID, Role: entity.Role(claims.Role)}, nil
}

// RequireAuth rejects any request without a valid access token (401), otherwise
// loads the caller's identity into the context.
func RequireAuth(mgr *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := bearerToken(c)
		if token == "" {
			abortUnauthorized(c)
			return
		}
		actor, err := actorFromToken(mgr, token)
		if err != nil {
			abortUnauthorized(c)
			return
		}
		setActor(c, actor)
		c.Next()
	}
}

// OptionalAuth attaches the caller's identity when a valid access token is
// present but never fails the request: anonymous callers continue with no
// actor. Used on public reads whose visibility adapts to the viewer.
func OptionalAuth(mgr *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if token := bearerToken(c); token != "" {
			if actor, err := actorFromToken(mgr, token); err == nil {
				setActor(c, actor)
			}
		}
		c.Next()
	}
}

// RequireRole gates a route by the caller's role (403 otherwise). It must run
// after RequireAuth. Resource-level ownership stays in the usecase layer — this
// is only the coarse role gate (spec §7.2).
func RequireRole(roles ...entity.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		actor, ok := ActorFromContext(c)
		if !ok || actor == nil {
			abortUnauthorized(c)
			return
		}
		if !actor.Role.HasRole(roles...) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    "FORBIDDEN",
				"message": "insufficient role",
			})
			return
		}
		c.Next()
	}
}

// abortUnauthorized writes the standard opaque 401 envelope.
func abortUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"code":    "UNAUTHORIZED",
		"message": "authentication required",
	})
}
