// Package output defines persistence contracts for comments and users.
package output

import (
	"context"

	models_auth "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/auth"
)

// UserRepository persists and retrieves user credentials.
type UserRepository interface {
	FindByUserName(ctx context.Context, username string) (*models_auth.User, error)

	// UserExists checks if a username is already registered in the system.
	// Returns true if the username exists in storage, false otherwise.
	// May return an error for database connectivity issues or storage failures.
	UserExists(ctx context.Context, username string) (bool, error)

	// SaveUser persists a new user record with secure credential storage.
	// Implementations should:
		// - Generate unique salt per user
		// - Hash password using cryptographically secure methods
		// - Prevent duplicate usernames

	// Returns error for:
		// - Duplicate username
		// - Invalid credentials
		// - Storage persistence failures
	SaveUser(ctx context.Context, user models_auth.User) (models_auth.User, error)

	EmailExists(ctx context.Context, email string) (bool, error)
}
