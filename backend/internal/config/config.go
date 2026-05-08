// Package config provides application configuration management for the ecommerce-platform application.
// It wraps Viper to load configuration from YAML files, environment variables, and defaults.
package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models"
	"github.com/redis/go-redis/v9"
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

	// Allow environment variables to override settings
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Default values for JWT, server port, rate limiting, static directory, and database
	config.SetDefault("server.port", "8080")
	config.SetDefault("STATIC_DIR", "./../frontend")

	config.SetDefault("rate_limiting.requests", 10.0)
	config.SetDefault("rate_limiting.burst", 5)
	config.SetDefault("rate_limiting.cleanup_minutes", 5)

	config.SetDefault("database.host", "localhost")
	config.SetDefault("database.port", 3306)
	config.SetDefault("database.name", "store_watches")

	config.SetDefault("redis.host", "localhost")
	config.SetDefault("redis.port", 6379)
	config.SetDefault("redis.db", 0)
	config.SetDefault("redis.pool_size", 10)
	config.SetDefault("redis.min_idle_conns", 5)
	config.SetDefault("redis.max_retries", 3)

	// Attempt to read the config file; log a warning if it fails
	if err := config.ReadInConfig(); err != nil {
		log.Printf("Warning: Error reading configuration file: %v", err)
		log.Println("Using default values and environment variable")
	}

	return &AppConfig{
		config: config,
	}
}

func ValidateRequiredConfig(config *viper.Viper) {
	required := []string{
		"security.jwt.jwt_secret",

		"database.user",
		"database.password",

		"redis.username",
		"redis.password",
		"redis.dial_timeout",
		"redis.read_timeout",
		"redis.write_timeout",
	}

	missing := []string{}
	for _, key := range required {
		value := config.GetString(key)
		if value == "" {
			missing = append(missing, key)
		}

		if key == "security.jwt.jwt_secret" && value == "your-secret-key" {
			log.Fatalf("Missing configuration: security.jwt.jwt_secret cannot use the insecure default value 'your-secret-key'. Please set a strong secret via ENV or config.yaml.")
		}
	}

	if len(missing) > 0 {
		log.Fatalf("Missing configuration (use ENV vars or config.yaml): %v", missing)
	}
}

func (a *AppConfig) GetString(key string) string {
	return a.config.GetString(key)
}

func (a *AppConfig) GetInt(key string) int {
	return a.config.GetInt(key)
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

func NewRedisClient(cfg models.RedisConfig) *redis.Client {
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

	// Verificar la conexión con timeout
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
	return a.config.GetString("security.jwt.jwt_secret")
}

func (a *AppConfig) GetRedisConfig() models.RedisConfig {
	return models.RedisConfig{
		Host:         a.config.GetString("redis.host"),
		Port:         a.config.GetString("redis.port"),
		Username:     a.config.GetString("redis.username"),
		Password:     a.config.GetString("redis.password"),
		DB:           a.config.GetInt("redis.db"),
		DialTimeout:  a.config.GetDuration("redis.dial_timeout"),
		ReadTimeout:  a.config.GetDuration("redis.read_timeout"),
		WriteTimeout: a.config.GetDuration("redis.write_timeout"),
		PoolSize:     a.config.GetInt("redis.pool_size"),
		MinIdleConns: a.config.GetInt("redis.min_idle_conns"),
		MaxRetries:   a.config.GetInt("redis.max_retries"),
	}
}

// GetRateLimitConfig returns a LimiterConfig populated from rate_limiting settings.
func (a *AppConfig) GetRateLimitConfig() models.LimiterConfig {
	return models.LimiterConfig{
		RequestPerSecond: a.config.GetFloat64("rate_limiting.requests"),
		Burst:            a.config.GetInt("rate_limiting.burst"),
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
