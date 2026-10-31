package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestNewAppConfig(t *testing.T) {
	tests := []struct {
		name      string
		env       map[string]string
		wantPort  string
		wantHost  string
		wantRedis string
		wantDB    int
	}{
		{
			name: "environment overrides",
			env: map[string]string{
				"LOG_LEVEL":   "DEBUG",
				"SERVER_PORT": "9090",
				"REDIS_HOST":  "redis.test",
				"REDIS_PORT":  "6380",
			},
			wantPort:  "9090",
			wantHost:  "redis.test",
			wantRedis: "6380",
			wantDB:    0,
		},
		{
			name: "defaults",
			env: map[string]string{
				"SERVER_PORT": "",
				"REDIS_HOST":  "",
				"REDIS_PORT":  "",
			},
			wantPort:  "8080",
			wantHost:  "localhost",
			wantRedis: "6379",
			wantDB:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.env {
				t.Setenv(key, value)
			}

			appConfig := NewAppConfig()
			redisConfig := appConfig.GetRedisConfig()

			if tt.name == "environment overrides" && appConfig.GetLogLevel() != "DEBUG" {
				t.Errorf("GetLogLevel() = %q, want %q", appConfig.GetLogLevel(), "DEBUG")
			}
			if tt.name == "defaults" && appConfig.GetLogLevel() != "INFO" {
				t.Errorf("GetLogLevel() = %q, want %q", appConfig.GetLogLevel(), "INFO")
			}

			if got := appConfig.GetPort(); got != tt.wantPort {
				t.Errorf("GetPort() = %q, want %q", got, tt.wantPort)
			}
			if redisConfig.Host != tt.wantHost {
				t.Errorf("RedisConfig.Host = %q, want %q", redisConfig.Host, tt.wantHost)
			}
			if redisConfig.Port != tt.wantRedis {
				t.Errorf("RedisConfig.Port = %q, want %q", redisConfig.Port, tt.wantRedis)
			}
			if redisConfig.DB != tt.wantDB {
				t.Errorf("RedisConfig.DB = %d, want %d", redisConfig.DB, tt.wantDB)
			}
		})
	}
}

func TestValidateRequiredConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  func() *viper.Viper
		wantErr bool
	}{
		{
			name:    "missing configuration",
			config:  viper.New,
			wantErr: true,
		},
		{
			name: "insecure JWT secret",
			config: func() *viper.Viper {
				config := validRequiredConfig()
				config.Set("SECURITY_JWT_JWT_SECRET", "your-secret-key")
				return config
			},
			wantErr: true,
		},
		{
			name:    "valid configuration",
			config:  validRequiredConfig,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRequiredConfig(tt.config())

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRequiredConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func validRequiredConfig() *viper.Viper {
	config := viper.New()
	for _, key := range []string{
		"LOG_LEVEL",
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
