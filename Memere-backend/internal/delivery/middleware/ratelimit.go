package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"

	infmetrics "github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/metrics"
)

// rateLimiter is a Redis-backed fixed-window counter keyed by an arbitrary key.
// The first hit in a window starts the TTL; subsequent hits INCR it. On
// overflow it returns 429 + Retry-After and emits the rate_limit_block_total
// metric for the given limiter label. Fail-open on Redis errors.
func rateLimiter(
	client *goredis.Client,
	keyFn func(*gin.Context) string,
	limiterLabel string,
	limit int,
	window time.Duration,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keyFn(c)
		ctx := c.Request.Context()

		count, err := client.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}
		if count == 1 {
			_ = client.Expire(ctx, key, window).Err()
		}

		if count > int64(limit) {
			retryAfter := window
			if ttl, err := client.TTL(ctx, key).Result(); err == nil && ttl > 0 {
				retryAfter = ttl
			}
			infmetrics.ObserveRateLimitBlock(limiterLabel)
			c.Writer.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    "TOO_MANY_REQUESTS",
				"message": "rate limit exceeded, please retry later",
			})
			return
		}
		c.Next()
	}
}

// ipKey builds a Redis key scoped by prefix and client IP.
func ipKey(prefix string) func(*gin.Context) string {
	return func(c *gin.Context) string { return prefix + ":ip:" + c.ClientIP() }
}

// userKey builds a Redis key scoped by prefix and the authenticated user ID
// (falls back to IP if unauthenticated, so the function is safe everywhere).
func userKey(prefix string) func(*gin.Context) string {
	return func(c *gin.Context) string {
		actor, ok := ActorFromContext(c)
		if ok && actor != nil {
			return prefix + ":uid:" + actor.UserID.String()
		}
		return prefix + ":ip:" + c.ClientIP()
	}
}

// RateLimit is the global per-IP limiter applied to every route: `rpm` requests
// per minute. It runs early in the middleware chain.
func RateLimit(client *goredis.Client, rpm int) gin.HandlerFunc {
	return rateLimiter(client, ipKey("rl:global"), "global", rpm, time.Minute)
}

// LoginRateLimit is the stricter per-IP limiter for auth routes (register,
// login, refresh). Default: 5 attempts per 15-minute window (spec §7.3).
func LoginRateLimit(client *goredis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return rateLimiter(client, ipKey("rl:auth"), "auth", limit, window)
}

// UserRateLimit is a per-authenticated-user limiter for expensive authenticated
// routes (e.g. payment initiation). Falls back to IP for anonymous callers.
// Default: 30 requests per minute per user.
func UserRateLimit(client *goredis.Client, limit int) gin.HandlerFunc {
	return rateLimiter(client, userKey("rl:user"), "user", limit, time.Minute)
}

// ProviderWebhookRateLimit limits inbound provider webhooks by provider name
// (extracted from the URL path param :provider). Prevents a misbehaving
// provider from flooding the queue. Default: 120 per minute per provider.
func ProviderWebhookRateLimit(client *goredis.Client, limit int) gin.HandlerFunc {
	return rateLimiter(client, func(c *gin.Context) string {
		return "rl:webhook:" + c.Param("provider")
	}, "provider", limit, time.Minute)
}
