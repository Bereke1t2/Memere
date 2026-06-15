package logger

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestRedactingHandler_ScrubsSensitiveKeys(t *testing.T) {
	sensitive := []string{
		"password", "token", "authorization", "api_key", "webhook_secret",
		"refresh_token", "private_key", "card_number",
	}
	safe := []string{"user_id", "course_id", "status", "latency"}

	var buf bytes.Buffer
	inner := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	l := slog.New(NewRedactingHandler(inner))

	// Log all sensitive and safe fields in one record.
	args := make([]any, 0, len(sensitive)+len(safe)*2)
	for _, k := range sensitive {
		args = append(args, k, "super-secret-value")
	}
	for _, k := range safe {
		args = append(args, k, "some-value")
	}
	l.Info("test", args...)

	output := buf.String()

	for _, k := range sensitive {
		if strings.Contains(output, "super-secret-value") {
			// The raw value should never appear for any sensitive key.
			t.Errorf("sensitive key %q leaked its value into the log output", k)
		}
		if !strings.Contains(output, "[REDACTED]") {
			t.Errorf("expected [REDACTED] placeholder for key %q", k)
		}
	}
	for _, k := range safe {
		if !strings.Contains(output, k+"=some-value") && !strings.Contains(output, k+`="some-value"`) {
			t.Errorf("safe key %q value was unexpectedly scrubbed or missing", k)
		}
	}
}

func TestRedactingHandler_WithAttrs_Scrubs(t *testing.T) {
	var buf bytes.Buffer
	inner := slog.NewTextHandler(&buf, nil)
	h := NewRedactingHandler(inner)
	h2 := h.WithAttrs([]slog.Attr{slog.String("password", "letmein"), slog.String("path", "/api")})

	// Verify the returned handler is still a RedactingHandler.
	if _, ok := h2.(*RedactingHandler); !ok {
		t.Error("WithAttrs should return a *RedactingHandler")
	}
	// Log via the child handler and confirm the password is redacted.
	l2 := slog.New(h2)
	l2.Info("attr-test")
	output := buf.String()
	if strings.Contains(output, "letmein") {
		t.Error("WithAttrs leaked password value")
	}
	_ = context.Background()
}

func TestIsSensitive(t *testing.T) {
	cases := []struct{ key string; want bool }{
		{"password", true},
		{"user_password", true},
		{"Token", true},
		{"Authorization", true},
		{"api_key", true},
		{"user_id", false},
		{"course_id", false},
		{"status", false},
	}
	for _, c := range cases {
		if got := isSensitive(c.key); got != c.want {
			t.Errorf("isSensitive(%q) = %v, want %v", c.key, got, c.want)
		}
	}
}
