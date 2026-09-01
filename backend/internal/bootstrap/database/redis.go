// Package bootstrap_database provides functions to initialize and configure database connections.
// This file specifically handles the retrieval of Redis configuration settings.
package bootstrap_database

import (
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/config"
)

// SetupDatabaseRedis extracts the Redis-specific configuration from the global application config.
// Unlike SQL databases that require an immediate connection string (DSN), this function retrieves a dedicated RedisConfig model which includes the address, password, and DB index needed to initialize the Redis client later in the application lifecycle.

// Parameters:
//   - appConfig: a pointer to the application configuration manager.
//
// Returns:
//   - config.RedisConfig: a struct containing all necessary parameters for Redis connectivity.
//   - error: returns nil by default, but follows the bootstrap signature for consistency.
func SetupDatabaseRedis(appConfig *config.AppConfig) (config.RedisConfig, error) {
	// Delegation to the config provider to extract Redis-specific values.
	redisConfig := appConfig.GetRedisConfig()
	return redisConfig, nil
}
