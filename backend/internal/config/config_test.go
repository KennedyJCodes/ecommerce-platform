package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestNewAppConfigLoadsEnvironmentOverrides(t *testing.T) {
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("REDIS_HOST", "redis.test")
	t.Setenv("REDIS_PORT", "6380")

	appConfig := NewAppConfig()

	if got := appConfig.GetPort(); got != "9090" {
		t.Fatalf("GetPort() = %q, want %q", got, "9090")
	}

	redisConfig := appConfig.GetRedisConfig()
	if redisConfig.Host != "redis.test" {
		t.Errorf("RedisConfig.Host = %q, want %q", redisConfig.Host, "redis.test")
	}
	if redisConfig.Port != "6380" {
		t.Errorf("RedisConfig.Port = %q, want %q", redisConfig.Port, "6380")
	}
}

func TestNewAppConfigLoadsDefaults(t *testing.T) {
	t.Setenv("SERVER_PORT", "")
	t.Setenv("REDIS_HOST", "")
	t.Setenv("REDIS_PORT", "")

	appConfig := NewAppConfig()

	if got := appConfig.GetPort(); got != "8080" {
		t.Fatalf("GetPort() = %q, want default %q", got, "8080")
	}

	redisConfig := appConfig.GetRedisConfig()
	if redisConfig.Host != "localhost" {
		t.Errorf("RedisConfig.Host = %q, want default %q", redisConfig.Host, "localhost")
	}
	if redisConfig.Port != "6379" {
		t.Errorf("RedisConfig.Port = %q, want default %q", redisConfig.Port, "6379")
	}
	if redisConfig.DB != 0 {
		t.Errorf("RedisConfig.DB = %d, want default %d", redisConfig.DB, 0)
	}
}

func TestValidateRequiredConfigReturnsErrorForMissingConfiguration(t *testing.T) {
	config := viper.New()

	err := ValidateRequiredConfig(config)

	if err == nil {
		t.Fatal("ValidateRequiredConfig() returned nil, want an error")
	}
}

func TestValidateRequiredConfigReturnsErrorForInsecureJWTSecret(t *testing.T) {
	config := validRequiredConfig()
	config.Set("SECURITY_JWT_JWT_SECRET", "your-secret-key")

	err := ValidateRequiredConfig(config)

	if err == nil {
		t.Fatal("ValidateRequiredConfig() returned nil, want an error")
	}
}

func TestValidateRequiredConfigReturnsNilForValidConfiguration(t *testing.T) {
	if err := ValidateRequiredConfig(validRequiredConfig()); err != nil {
		t.Fatalf("ValidateRequiredConfig() returned unexpected error: %v", err)
	}
}

func validRequiredConfig() *viper.Viper {
	config := viper.New()
	for _, key := range []string{
		"SECURITY_JWT_JWT_SECRET",
		"DATABASE_USER",
		"DATABASE_PASSWORD",
		"REDIS_USERNAME",
		"REDIS_PASSWORD",
		"REDIS_DIAL_TIMEOUT",
		"REDIS_READ_TIMEOUT",
		"REDIS_WRITE_TIMEOUT",
		"RESEND_EMAIL_API_KEY",
		"RESEND_EMAIL_FROM",
	} {
		config.Set(key, "configured-value")
	}
	return config
}
