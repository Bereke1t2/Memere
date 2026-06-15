// Package logger provides the application-wide structured logger built on
// stdlib log/slog. It is the only place slog.SetDefault is called; all other
// packages obtain a logger via FromContext (request-scoped) or slog.Default()
// (background workers).
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New builds a *slog.Logger tuned for the given environment and level string.
// In production it emits newline-delimited JSON (ELK/Cloud Logging compatible);
// in all other environments it emits human-readable text. It also calls
// slog.SetDefault so that bare slog.Info() calls in workers use this config.
func New(env, levelStr string) *slog.Logger {
	level := parseLevel(levelStr)
	opts := &slog.HandlerOptions{Level: level}

	var h slog.Handler
	if env == "production" {
		h = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		h = slog.NewTextHandler(os.Stdout, opts)
	}
	// Wrap with the redacting handler so sensitive values never appear in logs.
	h = NewRedactingHandler(h)
	l := slog.New(h)
	slog.SetDefault(l)
	return l
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
