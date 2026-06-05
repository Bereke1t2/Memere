// Package middleware holds the Gin middleware stack for the delivery layer:
// recovery, request-id, structured logging, CORS, Redis-backed rate limiting,
// and JWT authentication/authorization (spec §3.2, §7). Handlers read the
// authenticated caller via ActorFromContext.
package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/Bereke1t2/Memere/memere-backend/internal/usecase/course"
)

// Context keys. Using typed constants avoids collisions with other writers of
// the gin context.
const (
	actorKey     = "actor"
	requestIDKey = "request_id"
)

// setActor stores the authenticated caller on the context for handlers to read.
func setActor(c *gin.Context, actor *course.Actor) {
	c.Set(actorKey, actor)
}

// ActorFromContext returns the authenticated caller, or (nil, false) for an
// anonymous request. Handlers pass the returned *course.Actor straight into the
// usecase layer, where ownership/visibility is enforced.
func ActorFromContext(c *gin.Context) (*course.Actor, bool) {
	v, ok := c.Get(actorKey)
	if !ok {
		return nil, false
	}
	actor, ok := v.(*course.Actor)
	return actor, ok
}

// RequestIDFromContext returns the request id assigned by the RequestID
// middleware, or "" if none.
func RequestIDFromContext(c *gin.Context) string {
	if v, ok := c.Get(requestIDKey); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
