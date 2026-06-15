package logger

import (
	"context"
	"log/slog"
)

type ctxKey struct{}

// WithContext stores a request-scoped logger (pre-enriched with request_id,
// user_id etc.) in ctx. Middleware calls this; handlers and usecases retrieve
// it with FromContext.
func WithContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

// FromContext returns the request-scoped logger stored by middleware, or
// slog.Default() if none is present (background workers, tests).
func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok && l != nil {
		return l
	}
	return slog.Default()
}
