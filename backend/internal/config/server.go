package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	models_security "github.com/KennedyJCodes/ecommerce-platform/internal/core/domain/models/security"
)

// GetPort returns the HTTP server port as a string.
// It falls back to "8080" if not set or empty.
func (a *AppConfig) GetPort() string {
	port := a.config.GetString("SERVER_PORT")
	if port == "" {
		return "8080"
	}
	return port
}

// GetLogLevel returns the configured logging level.
// It falls back to "INFO" when LOG_LEVEL is not set.
func (a *AppConfig) GetLogLevel() string {
	logLevel := a.config.GetString("LOG_LEVEL")
	if logLevel == "" {
		return "INFO"
	}
	return logLevel
}

// GetRateLimitConfig returns a LimiterConfig populated from RATE_LIMITING_* configuration keys.
func (a *AppConfig) GetRateLimitConfig() models_security.LimiterConfig {
	return models_security.LimiterConfig{
		RequestPerSecond:   a.config.GetFloat64("RATE_LIMITING_REQUESTS"),
		Burst:              a.config.GetInt("RATE_LIMITING_BURST"),
		CleanupInterval:    time.Duration(a.config.GetInt("RATE_LIMITING_CLEANUP_MINUTES")) * time.Minute,
		ExpirationDuration: time.Duration(a.config.GetInt("RATE_LIMITING_EXPIRATION_MINUTES")) * time.Minute,
	}
}

// GetStaticDir returns the path to the static files directory.
// It verifies that the configured directory exists, and if not, attempts to resolve an alternate path relative to the executable.
// Logs a warning if neither path exists.
func (a *AppConfig) GetStaticDir() string {
	staticDir := a.config.GetString("STATIC_DIR")

	if staticDir == "" {
		staticDir = "./../frontend"
	}

	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		execPath, err := os.Executable()
		if err == nil {
			execDir := filepath.Dir(execPath)
			altPath := filepath.Join(execDir, "..", "frontend")
			if _, err := os.Stat(altPath); err == nil {
				return altPath
			}
		}
		log.Printf("Warning: Static directory '%s' not found", staticDir)
	}

	return staticDir
}

// IsProduction returns true if the ENV environment variable equals "production".
func (a *AppConfig) IsProduction() bool {
	return a.config.GetString("ENV") == "production"
}
