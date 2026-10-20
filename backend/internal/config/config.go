// Package config provides application configuration management for the ecommerce-platform application.
// It wraps Viper to load configuration from a .env file and OS environment variables, with sensible defaults.
package config

import (
	"strings"

	"github.com/spf13/viper"
)

// AppConfig holds the Viper instance for application-wide settings.
// It offers typed accessors for different configuration values and validation routines for security-sensitive settings.
type AppConfig struct {
	config *viper.Viper
}

// NewAppConfig initializes and returns a new AppConfig.
// It registers default values for all settings and enables automatic overrides via OS environment variables.
func NewAppConfig() *AppConfig {
	config := viper.New()

	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	setDefaults(config)

	return &AppConfig{
		config: config,
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

// GetConfig exposes the underlying Viper instance for advanced use cases.
func (a *AppConfig) GetConfig() *viper.Viper {
	return a.config
}
