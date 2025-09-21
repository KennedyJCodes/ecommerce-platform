// Package config provides application configuration management for the sale-watches application.
// It wraps Viper to load configuration from YAML files, environment variables, and defaults.
package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/models"
	"github.com/spf13/viper"
)

// AppConfig holds the Viper instance for application-wide settings.
// It offers typed accessors for different configuration values and validation routines for security-sensitive settings.
type AppConfig struct {
	config *viper.Viper
}

// NewAppConfig initializes and returns a new AppConfig.
// It sets up Viper to read from a YAML file named "config" in the ./internal/config directory, registers default values, enables automatic environment variable overrides, and logs warnings if the config file cannot be read.
func NewAppConfig() *AppConfig {
	config := viper.New()

	// Configuration file settings
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	config.AddConfigPath("./internal/config")

	// Default values for JWT, server port, rate limiting, static directory, and database
	config.SetDefault("security.jwt.jwt_secret", "your-secret-key")

	config.SetDefault("server.port", "8080")
	config.SetDefault("rate_limiting.requests", 10.0)
	config.SetDefault("rate_limiting.cleanup_minutes", 5)

	config.SetDefault("STATIC_DIR", "./../frontend")

	config.SetDefault("database.user", "root")
	config.SetDefault("database.password", "password")
	config.SetDefault("database.host", "localhost")
	config.SetDefault("database.port", 3306)
	config.SetDefault("database.name", "store_watches")

	// Allow environment variables to override settings
	config.AutomaticEnv()

	// Attempt to read the config file; log a warning if it fails
	if err := config.ReadInConfig(); err != nil {
		log.Printf("Warning: Error reading configuration file: %v", err)
		log.Println("Using default values and environment variable")
	}

	return &AppConfig{
		config: config,
	}
}

// GetPort returns the HTTP server port as a string.
// It falls back to "8080" if not set.
func (a *AppConfig) GetPort() string {
	port := a.config.GetString("server.port")
	if port == "" {
		return "8080"
	}
	return port
}

// GetConfig exposes the underlying Viper instance for advanced use cases.
func (a *AppConfig) GetConfig() *viper.Viper {
	return a.config
}

// GetJWTSecret retrieves the JWT secret key from configuration.
func (a *AppConfig) GetJWTSecret() string {
	return a.config.GetString("security.jwt.jwt_secret")
}

// GetRateLimitConfig returns a LimiterConfig populated from rate_limiting settings.
func (a *AppConfig) GetRateLimitConfig() models.LimiterConfig {
	return models.LimiterConfig{
		RequestPerSecond: a.config.GetFloat64("rate_limiting.requests"),
		Burst:            a.config.GetInt("rate_limiting.cleanup_minutes"),
	}
}

// GetStaticDir returns the path to the static files directory.
// It verifies that the configured directory exists, and if not, attempts to resolve an alternate path relative to the executable.
// Logs a warning if neither path exists.
func (a *AppConfig) GetStaticDir() string {
	// Get the value from the configuration
	staticDir := a.config.GetString("STATIC_DIR")

	// If empty, use a default value
	if staticDir == "" {
		staticDir = "./../frontend"
	}

	// Check if the directory exists
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		// Fallback: resolve relative to executable location
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

// ValidateConfig performs sanity checks on critical settings.
// Currently warns if the default JWT secret is used in production
func (a *AppConfig) ValidateConfig() {
	if a.GetJWTSecret() == "your-secret-key" && a.IsProduction() {
		log.Println("WARNING: Using default JWT key in production, this is insecure")
	}
}
