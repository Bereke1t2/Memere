package logger

import (
	"context"
	"log/slog"
	"strings"
)

// sensitiveKeys is the deny-list of attribute key substrings that must never
// appear unredacted in any log line (Non-Negotiable #8, spec §7.3).
var sensitiveKeys = []string{
	"password", "passwd", "secret", "token", "authorization",
	"card", "cvv", "pan", "pin", "webhook_secret", "api_key", "apikey",
	"private_key", "access_key", "refresh_token", "fcm_key", "sendgrid",
}

const redacted = "[REDACTED]"

// isSensitive returns true when the attribute key contains any deny-listed
// substring (case-insensitive).
func isSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, s := range sensitiveKeys {
		if strings.Contains(lower, s) {
			return true
		}
	}
	return false
}

// RedactingHandler wraps another slog.Handler and replaces the value of any
// attribute whose key is considered sensitive with the constant "[REDACTED]".
type RedactingHandler struct {
	inner slog.Handler
}

// NewRedactingHandler wraps inner so all sensitive attribute values are scrubbed.
func NewRedactingHandler(inner slog.Handler) *RedactingHandler {
	return &RedactingHandler{inner: inner}
}

func (h *RedactingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *RedactingHandler) Handle(ctx context.Context, r slog.Record) error {
	safe := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	r.Attrs(func(a slog.Attr) bool {
		safe.AddAttrs(scrubAttr(a))
		return true
	})
	return h.inner.Handle(ctx, safe)
}

func (h *RedactingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	scrubbed := make([]slog.Attr, len(attrs))
	for i, a := range attrs {
		scrubbed[i] = scrubAttr(a)
	}
	return &RedactingHandler{inner: h.inner.WithAttrs(scrubbed)}
}

func (h *RedactingHandler) WithGroup(name string) slog.Handler {
	return &RedactingHandler{inner: h.inner.WithGroup(name)}
}

func scrubAttr(a slog.Attr) slog.Attr {
	if isSensitive(a.Key) {
		return slog.String(a.Key, redacted)
	}
	return a
}
