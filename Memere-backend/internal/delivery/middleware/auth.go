package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
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
func actorFromToken(mgr *jwt.Manager, token string) (*course.Actor, *jwt.Claims, error) {
	claims, err := mgr.Verify(token, jwt.TokenTypeAccess)
	if err != nil {
		return nil, nil, err
	}
	return &course.Actor{UserID: claims.UserID, Role: entity.Role(claims.Role)}, claims, nil
}

// TokenChecker is the subset of SessionRepository needed by RequireAuth to
// consult the access-token denylist. Injected at router construction so the
// middleware has no import cycle into the auth usecase.
type TokenChecker interface {
	IsTokenDenied(ctx context.Context, jti string) (bool, error)
}

// RequireAuth rejects any request without a valid access token (401), otherwise
// loads the caller's identity into the context. It also checks the JTI denylist
// so that revoked tokens (logout/suspend) are rejected before their expiry.
func RequireAuth(mgr *jwt.Manager, denylist ...repository.SessionRepository) gin.HandlerFunc {
	var checker TokenChecker
	if len(denylist) > 0 && denylist[0] != nil {
		checker = denylist[0]
	}
	return func(c *gin.Context) {
		token := bearerToken(c)
		if token == "" {
			abortUnauthorized(c)
			return
		}
		actor, claims, err := actorFromToken(mgr, token)
		if err != nil {
			abortUnauthorized(c)
			return
		}
		// Denylist check: reject revoked JTIs (logout / suspend).
		if checker != nil {
			denied, err := checker.IsTokenDenied(c.Request.Context(), claims.ID)
			if err == nil && denied {
				abortUnauthorized(c)
				return
			}
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
			if actor, _, err := actorFromToken(mgr, token); err == nil {
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
