package config

import (
	"log"

	"github.com/spf13/viper"
)

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
			log.Fatalf("Missing configuration: SECURITY_JWT_JWT_SECRET cannot use the insecure default value 'your-secret-key'. Please set a strong secret via ENV or .env file.")
		}
	}

	if len(missing) > 0 {
		log.Fatalf("Missing configuration (use ENV vars or .env file): %v", missing)
	}
}
