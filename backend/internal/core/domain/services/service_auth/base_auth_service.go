// Package service_auth provides implementations of input port interfaces for authentication.
// It contains shared logic for validating credentials, checking user existence, and
// generating JWT tokens.
package service_auth

import (
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
	securityAuth "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/security_auth"
)

// BaseAuthService offers common authentication operations shared by login and registration services. It delegates credential validation to injected validators, user lookup to the UserRepository, and token creation to the security_auth package.
type BaseAuthService struct {
	// UserRepo provides access to user persistence (e.g., lookup by username).
	UserRepo output.UserRepository

	// UserNameValidator enforces rules on allowed username formats.
	UserNameValidator input.Validator

	// PasswordValidator enforces rules on allowed password formats.
	PasswordValidator input.Validator

	// CSRFService handles CSRF token generation and validation.
	CSRFService input.CSRFService

	// CSRFCookieSetter sets CSRF token cookies in HTTP responses.
	CSRFCookieSetter output.CSRFCookieSetter
}

// ValidateUserName checks the supplied username against the UserNameValidator.
// Returns a ValidationError if the username is invalid.
func (b *BaseAuthService) ValidateUserName(username interface{}) error {
	if err := b.UserNameValidator.Validate(username); err != nil {
		return errors.NewValidationError(errors.ErrInvalidUsername)
	}
	return nil
}

// ValidatePassword checks the supplied password against the PasswordValidator.
// Returns a ValidationError if the password is invalid.
func (b *BaseAuthService) ValidatePassword(password interface{}) error {
	if err := b.PasswordValidator.Validate(password); err != nil {
		return errors.NewValidationError(errors.ErrInvalidPassword)
	}
	return nil
}

// CheckUserExists queries the UserRepository to determine if a user with the given username already exists. Returns (true, nil) if found, (false, nil) if not, or an InternalError if the lookup fails.
func (b *BaseAuthService) CheckUserExists(username string) (bool, error) {
	exists, err := b.UserRepo.UserExists(username)
	if err != nil {
		return false, errors.NewInternalError(errors.ErrDatabaseQuery).WithError(err)
	}
	return exists, nil
}

// GenerateToken creates a signed JWT for the given username using the default JWT service. Returns the token string or an InternalError if token generation fails.
func (b *BaseAuthService) GenerateToken(userId int, username string) (string, error) {
	token, err := securityAuth.GenerateJWT(userId, username)
	if err != nil {
		return "", errors.NewInternalError(errors.ErrTokenGeneration).WithError(err)
	}

	return token, nil
}

// GenerateAndSetCSRFToken generates a CSRF token for the user and sets it as a cookie.
// Returns an error if token generation fails.
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
