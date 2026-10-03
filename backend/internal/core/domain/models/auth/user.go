// Package models defines core domain entities for the ecommerce-platform application.
// It contains simple data structures that represent business objects throughout the system.
package models_auth

import (
	"time"
)

// User represents a user account in the system.

// Fields:
//   - user_id: the unique identifier for the user.
//   - username: the unique identifier chosen by the user for login purposes.
//   - email: the user's email address; it should be unique and validated.
//   - passwordhash: the hashed version of the user's password; it should be stored securely and never logged or exposed in plain text.
//   - verified_at: the timestamp when the user's email was verified.
//   - created_at: the timestamp when the user account was created.
//   - updated_at: the timestamp when the user account was last updated.
type User struct {
	UserID uint64
	UserName string
	Email string
	Password string
	VerifiedAt time.Time
	CreatedAt time.Time
    UpdatedAt time.Time
}
