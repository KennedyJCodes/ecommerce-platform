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

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/security"
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

	// Default values for server, database, Redis, and rate limiting
	config.SetDefault("SERVER_PORT", "8080")
	config.SetDefault("STATIC_DIR", "./../frontend")
	config.SetDefault("RATE_LIMITING_REQUESTS", 10.0)
	config.SetDefault("RATE_LIMITING_BURST", 5)
	config.SetDefault("RATE_LIMITING_CLEANUP_MINUTES", 5)
	config.SetDefault("RATE_LIMITING_EXPIRATION_MINUTES", 30)
	config.SetDefault("DATABASE_HOST", "localhost")
	config.SetDefault("DATABASE_PORT", 3306)
	config.SetDefault("DATABASE_NAME", "ecommerce-mysql")
	config.SetDefault("REDIS_HOST", "localhost")
	config.SetDefault("REDIS_PORT", 6379)
	config.SetDefault("REDIS_DB", 0)
	config.SetDefault("REDIS_POOL_SIZE", 10)
	config.SetDefault("REDIS_MIN_IDLE_CONNS", 5)
	config.SetDefault("REDIS_MAX_RETRIES", 3)
	config.SetDefault("COOKIE_PREFIX", "")
	config.SetDefault("DATABASE_TLS", false)
	config.SetDefault("SSL_ENABLED", false)
	config.SetDefault("SSL_CERT_FILE", "")
	config.SetDefault("SSL_KEY_FILE", "")

	return &AppConfig{
		config: config,
	}
}

// ValidateRequiredConfig checks that all mandatory configuration keys are present and non-empty.
// It terminates the application with a fatal error if any required key is missing or uses an insecure default.
func ValidateRequiredConfig(config *viper.Viper) {
	required := []string{
		"SECURITY_JWT_JWT_SECRET",
		"DATABASE_USER",
		"DATABASE_PASSWORD",
		"REDIS_USERNAME",
		"REDIS_PASSWORD",
		"REDIS_DIAL_TIMEOUT",
		"REDIS_READ_TIMEOUT",
		"REDIS_WRITE_TIMEOUT",
	}

	missing := []string{}
	for _, key := range required {
		value := config.GetString(key)
		if value == "" {
			missing = append(missing, key)
		}

		if key == "SECURITY_JWT_JWT_SECRET" && value == "your-secret-key" {
			log.Fatalf("Missing configuration: SECURITY_JWT_JWT_SECRET cannot use the insecure default value 'your-secret-key'. Please set a strong secret via ENV or .env file.")
		}
	}

	if len(missing) > 0 {
		log.Fatalf("Missing configuration (use ENV vars or .env file): %v", missing)
	}
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
