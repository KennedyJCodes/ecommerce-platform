package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// TimingConfig contains the threshold used to identify slow requests.
type TimingConfig struct {
	WarningThreshold time.Duration
}

// DefaultTimingConfig creates the default request logging configuration.
func DefaultTimingConfig() *TimingConfig {
	return &TimingConfig{WarningThreshold: 500 * time.Millisecond}
}

// LoggingMiddleware logs each request once, including its method, path, status,
// duration, and whether it exceeded the configured slow-request threshold.
func LoggingMiddleware(config *TimingConfig) Middleware {
	if config == nil {
		config = DefaultTimingConfig()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := NewResponseWriter(w)

			next.ServeHTTP(rw, r)

			duration := time.Since(start)
			slow := duration > config.WarningThreshold
			level := slog.LevelInfo
			if rw.statusCode >= 400 {
				level = slog.LevelError
			} else if slow {
				level = slog.LevelWarn
			}

			slog.LogAttrs(r.Context(), level, "HTTP request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", rw.statusCode),
				slog.Duration("duration", duration),
				slog.Bool("slow", slow),
			)
		})
	}
}

// ResponseWriter captures the response status code for request logging.
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewResponseWriter wraps an http.ResponseWriter with status tracking.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

// WriteHeader captures and forwards the response status code.
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
