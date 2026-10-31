package config

import "github.com/spf13/viper"

// setDefaults registers the default values for the application configuration.
func setDefaults(config *viper.Viper) {
	// Server configuration defaults.
	config.SetDefault("SERVER_PORT", "8080")
	config.SetDefault("STATIC_DIR", "./../frontend")

	// Logging configuration defaults.
	config.SetDefault("LOG_LEVEL", "INFO")

	// Rate limiting configuration defaults.
	config.SetDefault("RATE_LIMITING_REQUESTS", 10.0)
	config.SetDefault("RATE_LIMITING_BURST", 5)
	config.SetDefault("RATE_LIMITING_CLEANUP_MINUTES", 5)
	config.SetDefault("RATE_LIMITING_EXPIRATION_MINUTES", 30)

	// MySQL database configuration defaults.
	config.SetDefault("DATABASE_HOST", "localhost")
	config.SetDefault("DATABASE_PORT", 3306)
	config.SetDefault("DATABASE_NAME", "ecommerce-mysql")

	// Redis configuration defaults.
	config.SetDefault("REDIS_HOST", "localhost")
	config.SetDefault("REDIS_PORT", 6379)
	config.SetDefault("REDIS_DB", 0)
	config.SetDefault("REDIS_POOL_SIZE", 10)
	config.SetDefault("REDIS_MIN_IDLE_CONNS", 5)
	config.SetDefault("REDIS_MAX_RETRIES", 3)

	// Cookie configuration defaults.
	config.SetDefault("COOKIE_PREFIX", "")

	// Database TLS and SSL configuration defaults.
	config.SetDefault("DATABASE_TLS", false)
	config.SetDefault("SSL_ENABLED", false)
	config.SetDefault("SSL_CERT_FILE", "")
	config.SetDefault("SSL_KEY_FILE", "")
}
