package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
)

// rateLimiter is a Redis-backed fixed-window limiter keyed by client IP. The
// first request in a window sets the counter's TTL; subsequent requests INCR it.
// When the count exceeds the limit within the window the request is rejected
// with 429 and a Retry-After header (spec §7.3).
func rateLimiter(client *goredis.Client, keyPrefix string, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keyPrefix + ":" + c.ClientIP()
		ctx := c.Request.Context()

		count, err := client.Incr(ctx, key).Result()
		if err != nil {
			// Fail open: a limiter outage must not take the API down. The request
			// proceeds unthrottled rather than erroring.
			c.Next()
			return
		}
		if count == 1 {
			// First hit in this window — start the expiry.
			_ = client.Expire(ctx, key, window).Err()
		}

		if count > int64(limit) {
			retryAfter := window
			if ttl, err := client.TTL(ctx, key).Result(); err == nil && ttl > 0 {
				retryAfter = ttl
			}
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

// RateLimit is the global per-IP limiter: `rpm` requests per minute.
func RateLimit(client *goredis.Client, rpm int) gin.HandlerFunc {
	return rateLimiter(client, "rl:global", rpm, time.Minute)
}

// LoginRateLimit is the stricter limiter for POST /auth/login (default 5
// attempts per 15 minutes per IP, spec §7.3) to blunt credential stuffing.
func LoginRateLimit(client *goredis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return rateLimiter(client, "rl:login", limit, window)
}
