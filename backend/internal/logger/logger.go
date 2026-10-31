// Package logger provides the application's structured logging primitives.
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New creates a logger configured with the requested log level and output
// format. Production logs are emitted as JSON; every other environment uses
// slog's human-readable text format.
//
// An unknown or empty log level defaults to slog.LevelInfo.
func New(logLevel, environment string) *slog.Logger {
	options := &slog.HandlerOptions{
		Level: parseLevel(logLevel),
	}

	if strings.EqualFold(strings.TrimSpace(environment), "production") {
		return slog.New(slog.NewJSONHandler(os.Stdout, options))
	}

	return slog.New(slog.NewTextHandler(os.Stdout, options))
}

func parseLevel(logLevel string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(logLevel)) {
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
