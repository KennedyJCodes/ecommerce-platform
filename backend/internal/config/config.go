// Package config provides application configuration management for the ecommerce-platform application.
// It wraps Viper to load configuration from a .env file and OS environment variables, with sensible defaults.
package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/KennedyJCodes/ecommerce-platform/internal/core/domain/models/security"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// AppConfig holds the Viper instance for application-wide settings.
// It offers typed accessors for different configuration values and validation routines for security-sensitive settings.
type AppConfig struct {
	config *viper.Viper
}

// NewAppConfig initializes and returns a new AppConfig.
// It reads a .env file from ./internal/config/.env, registers default values for all settings, and enables automatic
// overrides via OS environment variables. OS environment variables take precedence over the .env file.
func NewAppConfig() *AppConfig {
	config := viper.New()

	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	setDefaults(config)

	return &AppConfig{
		config: config,
	}
}

func (a *AppConfig) GetResendEmailAPIKey() string {
	return a.config.GetString("RESEND_EMAIL_API_KEY")
}

func (a *AppConfig) GetResendEmailFrom() string {
	return a.config.GetString("RESEND_EMAIL_FROM")
}

// GetString returns a string value from the configuration by key.
func (a *AppConfig) GetString(key string) string {
	return a.config.GetString(key)
}

// GetInt returns an integer value from the configuration by key.
func (a *AppConfig) GetInt(key string) int {
	return a.config.GetInt(key)
}

// GetPort returns the HTTP server port as a string.
// It falls back to "8080" if not set or empty.
func (a *AppConfig) GetPort() string {
	port := a.config.GetString("SERVER_PORT")
	if port == "" {
		return "8080"
	}
	return port
}

// NewRedisClient creates and returns a new Redis client connected to the configured Redis instance.
// It verifies the connection with a 5-second timeout and terminates the application on failure.
func NewRedisClient(cfg RedisConfig) *redis.Client {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
	})

	// Verify the connection with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis on %s: %v", addr, err)
	}

	log.Printf("Successfully connected to Redis on %s - Response: %s", addr, pong)
	return client
}

// GetConfig exposes the underlying Viper instance for advanced use cases.
func (a *AppConfig) GetConfig() *viper.Viper {
	return a.config
}

// GetJWTSecret retrieves the JWT secret key from configuration.
func (a *AppConfig) GetJWTSecret() string {
	return a.config.GetString("SECURITY_JWT_JWT_SECRET")
}

// GetRedisConfig returns a RedisConfig struct populated from the REDIS_* configuration keys.
func (a *AppConfig) GetRedisConfig() RedisConfig {
	return RedisConfig{
		Host:         a.config.GetString("REDIS_HOST"),
		Port:         a.config.GetString("REDIS_PORT"),
		Username:     a.config.GetString("REDIS_USERNAME"),
		Password:     a.config.GetString("REDIS_PASSWORD"),
		DB:           a.config.GetInt("REDIS_DB"),
		DialTimeout:  a.config.GetDuration("REDIS_DIAL_TIMEOUT"),
		ReadTimeout:  a.config.GetDuration("REDIS_READ_TIMEOUT"),
		WriteTimeout: a.config.GetDuration("REDIS_WRITE_TIMEOUT"),
		PoolSize:     a.config.GetInt("REDIS_POOL_SIZE"),
		MinIdleConns: a.config.GetInt("REDIS_MIN_IDLE_CONNS"),
		MaxRetries:   a.config.GetInt("REDIS_MAX_RETRIES"),
	}
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

// GetCookiePrefix returns the prefix prepended to all cookie names.
// Used to avoid collisions when multiple applications share a parent domain.
func (a *AppConfig) GetCookiePrefix() string {
	return a.config.GetString("COOKIE_PREFIX")
}

// IsDatabaseTLSEnabled returns whether the MySQL connection should use TLS.
func (a *AppConfig) IsDatabaseTLSEnabled() bool {
	return a.config.GetBool("DATABASE_TLS")
}

// IsSSLEnabled returns whether SSL/TLS should be used.
func (a *AppConfig) IsSSLEnabled() bool {
	return a.config.GetBool("SSL_ENABLED")
}

// GetSSLCertFile returns the path to the SSL certificate file.
func (a *AppConfig) GetSSLCertFile() string {
	return a.config.GetString("SSL_CERT_FILE")
}

// GetSSLKeyFile returns the path to the SSL private key file.
func (a *AppConfig) GetSSLKeyFile() string {
	return a.config.GetString("SSL_KEY_FILE")
}
