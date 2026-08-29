// Package input defines interfaces for user authentication, registration, and input validation operations.
// These interfaces serve as contracts for service implementations that handle core application logic.
package input

import (
	"context"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/dto/auth"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models"
)

// UserServiceLogin defines the interface for user authentication operations.
// Implementations should verify user credentials and provide JWT tokens for authenticated sessions.
type UserServiceLogin interface {
	// Login authenticates a user using account credentials.
	// Returns a TokenPair containing access and refresh tokens, the CSRF token,
	// or an error if verification fails.
	Login(ctx context.Context, request dto.LoginRequest) (*models.TokenPair, string, error)
}

// UserServiceRegister defines the interface for user registration operations.
// Implementations should handle new user account creation and provide JWT tokens upon successful registration.
type UserServiceRegister interface {
	// Register creates a new user account with provided credentials.
	// Returns a TokenPair containing access and refresh tokens, the CSRF token,
	// or an error if registration fails.
	Register(ctx context.Context, request dto.RegisterAccount) (*models.TokenPair, string, error)
}
