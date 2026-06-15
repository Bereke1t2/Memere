package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"

	"github.com/Bereke1t2/Memere/memere-backend/internal/infrastructure/tracing"
)

// Tracing starts an OpenTelemetry span for each request and propagates
// TraceContext / Baggage from incoming W3C headers. The span is named after
// the route template (e.g. "GET /api/v1/courses/:id") so every request to the
// same endpoint shares one span name — keeping the trace namespace manageable.
//
// Attributes follow the OTel HTTP semantic conventions (semconv v1.30).
func Tracing() gin.HandlerFunc {
	tracer := tracing.Tracer("memere.http")

	return func(c *gin.Context) {
		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}
		spanName := c.Request.Method + " " + route

		ctx, span := tracer.Start(c.Request.Context(), spanName)
		defer span.End()

		// Attach the enriched context so downstream usecases can add child spans.
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		span.SetAttributes(
			semconv.HTTPRequestMethodKey.String(c.Request.Method),
			semconv.HTTPRouteKey.String(route),
			semconv.HTTPResponseStatusCodeKey.Int(c.Writer.Status()),
			attribute.String("request_id", RequestIDFromContext(c)),
		)
	}
}
