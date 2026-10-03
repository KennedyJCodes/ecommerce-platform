// Package models defines core domain entities for the ecommerce-platform application.
package models_security

import "time"

// CSRFToken represents a CSRF token entity in the domain.
// It contains the token value, associated user ID, and expiration information.
type CSRFToken struct {
	Value     string
	UserID    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// IsValid checks if the CSRF token is still valid (not expired).
func (t *CSRFToken) IsValid() bool {
	return time.Now().Before(t.ExpiresAt)
}
