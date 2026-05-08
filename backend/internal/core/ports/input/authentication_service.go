// Package input defines interfaces for user authentication, registration, and input validation operations.
// These interfaces serve as contracts for service implementations that handle core application logic.
package input

import "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models"

// UserServiceLogin defines the interface for user authentication operations.
// Implementations should verify user credentials and provide JWT tokens for authenticated sessions.
type UserServiceLogin interface {
	// Login authenticates a user using account credentials.
	// Returns a JWT token string for successful authentication or an error if verification fails.
	Login(account models.Account) (string, error)
}

// UserServiceRegister defines the interface for user registration operations.
// Implementations should handle new user account creation and provide JWT tokens upon successful registration.
type UserServiceRegister interface {
	// Register creates a new user account with provided credentials.
	// Returns a JWT token string for the newly created account or an error if registration fails.
	Register(account models.Account) (string, error)
}
