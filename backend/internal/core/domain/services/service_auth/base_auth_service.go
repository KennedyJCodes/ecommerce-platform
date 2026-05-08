// Package service_auth provides implementations of input port interfaces for authentication.
// It centralizes the core logic for security flows, ensuring consistent handling of credentials and identity tokens across the application.
package service_auth

import (
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
	securityAuth "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/security_auth"
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

	// CSRFService: domain service to manage the lifecycle of CSRF tokens.
	CSRFService input.CSRFService

	// CSRFCookieSetter: output port to bridge domain tokens with HTTP transport cookies.
	CSRFCookieSetter output.CSRFCookieSetter
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

// GenerateToken wraps the security package logic to create an identity JWT.
// It maps technical token generation errors into domain-understandable InternalErrors.
func (b *BaseAuthService) GenerateToken(userId int, username string) (string, error) {
	token, err := securityAuth.GenerateJWT(userId, username)
	if err != nil {
		return "", errors.NewInternalError(errors.ErrTokenGeneration).WithError(err)
	}

	return token, nil
}

// GenerateAndSetCSRFToken orchestrates the creation of a CSRF token and its subsequent delivery to the client via a cookie setter.
// It fails silently if services are not injected, allowing for optional CSRF flows.
func (b *BaseAuthService) GenerateAndSetCSRFToken(userID string) error {
	if b.CSRFService == nil || b.CSRFCookieSetter == nil {
		return nil
	}

	csrfToken, err := b.CSRFService.GenerateToken(userID)
	if err != nil {
		return errors.NewInternalError(errors.ErrTokenGeneration).WithError(err)
	}

	b.CSRFCookieSetter.SetCSRFCookie(csrfToken)
	return nil
}
