// Package service_auth provides implementations of input port interfaces for authentication.
// It centralizes the core logic for security flows, ensuring consistent handling of credentials and identity tokens across the application.
package service_auth

import (
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
)

// BaseAuthService serves as a foundational structure for authentication use cases.
// Instead of duplicating logic, it provides shared methods for:
//  1. Validating input formats via injected strategy validators.
//  2. Interfacing with persistence for identity verification.
//  3. Orchestrating the generation of session (JWT) and security (CSRF) tokens.
type BaseAuthService struct {
	// UserRepo: output port for user data persistence and existence checks.
	UserRepo output.UserRepository

	// UserNameValidator: strategy to enforce username complexity and format rules.
	UserNameValidator input.Validator

	// PasswordValidator: strategy to enforce password security requirements.
	PasswordValidator input.Validator

	// TokenService: domain service for JWT generation and validation.
	TokenService output.TokenService

	// CSRFService: domain service to manage the lifecycle of CSRF tokens.
	CSRFService output.CSRFService
}

// ValidateUserName evaluates if the provided username meets business requirements.
// Returns a domain-specific ValidationError if the format is rejected.
func (b *BaseAuthService) ValidateUserName(username interface{}) error {
	if err := b.UserNameValidator.Validate(username); err != nil {
		return errors.NewValidationError(errors.ErrInvalidUsername)
	}
	return nil
}

// ValidatePassword evaluates if the provided password meets security standards.
// Returns a domain-specific ValidationError if the requirements are not met.
func (b *BaseAuthService) ValidatePassword(password interface{}) error {
	if err := b.PasswordValidator.Validate(password); err != nil {
		return errors.NewValidationError(errors.ErrInvalidPassword)
	}
	return nil
}

// CheckUserExists verifies the presence of a username in the persistence layer.
// Returns (true, nil) if the user is registered, or handles database errors gracefully.
func (b *BaseAuthService) CheckUserExists(username string) (bool, error) {
	exists, err := b.UserRepo.UserExists(username)
	if err != nil {
		return false, errors.NewInternalError(errors.ErrDatabaseQuery).WithError(err)
	}
	return exists, nil
}

// GenerateTokenPair wraps the security package logic to create an access JWT
// and a refresh JWT. Returns a TokenPair or an InternalError on failure.
func (b *BaseAuthService) GenerateTokenPair(userId int, username string) (*models.TokenPair, error) {
	accessToken, err := b.TokenService.GenerateToken(userId, username, models.TokenTypeAccess)
	if err != nil {
		return nil, errors.NewInternalError(errors.ErrTokenGeneration).WithError(err)
	}

	refreshToken, err := b.TokenService.GenerateToken(userId, username, models.TokenTypeRefresh)
	if err != nil {
		return nil, errors.NewInternalError(errors.ErrTokenGeneration).WithError(err)
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GenerateCSRFToken creates a CSRF token for the given userID.
// It fails silently if CSRFService is not configured, allowing for optional CSRF flows.
// Returns the generated token string, or an empty string and nil if CSRF is disabled.
func (b *BaseAuthService) GenerateCSRFToken(userID string) (string, error) {
	if b.CSRFService == nil {
		return "", nil
	}

	csrfToken, err := b.CSRFService.GenerateToken(userID)
	if err != nil {
		return "", errors.NewInternalError(errors.ErrTokenGeneration).WithError(err)
	}

	return csrfToken, nil
}
