// Package models defines core domain entities and configuration structs for the sale‑watches application.
package models_security

import "time"

// LimiterConfig holds the settings for rate‑limiting behavior.

// The fields are decoded via mapstructure tags when unmarshalling configuration (e.g., with Viper) into this struct. :contentReference[oaicite:0]{index=0}

// RequestPerSecond: the average number of allowed requests per second.
// Burst: the maximum number of requests allowed to exceed the steady‑state rate in a short burst :contentReference[oaicite:1]{index=1}
type LimiterConfig struct {
	// RequestPerSecond specifies the allowed request rate (requests/sec).
	RequestPerSecond float64 `mapstructure:"request_rate"`

	// Burst specifies the maximum burst size over the steady request rate.
	Burst int `mapstructure:"burst"`

	// CleanupInterval specifies how often inactive limiter entries are purged.
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`

	// ExpirationDuration specifies how long an IP limiter can stay idle before removal.
	ExpirationDuration time.Duration `mapstructure:"expiration_duration"`
}
