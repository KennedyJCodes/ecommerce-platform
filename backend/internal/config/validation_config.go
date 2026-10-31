package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ValidateRequiredConfig checks that all mandatory configuration keys are present and non-empty.
// It returns an error if any required key is missing or uses an insecure default.
func ValidateRequiredConfig(config *viper.Viper) error {
	required := []string{
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
	}

	missing := []string{}
	for _, key := range required {
		value := config.GetString(key)
		if value == "" {
			missing = append(missing, key)
		}

		if key == "SECURITY_JWT_JWT_SECRET" && value == "your-secret-key" {
			return fmt.Errorf("missing configuration: SECURITY_JWT_JWT_SECRET cannot use the insecure default value 'your-secret-key'. Please set a strong secret via ENV or .env file")
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing configuration (use ENV vars or .env file): %v", missing)
	}

	return nil
}
