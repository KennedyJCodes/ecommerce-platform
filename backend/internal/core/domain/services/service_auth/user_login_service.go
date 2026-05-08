// Package service_auth provides implementations of input port interfaces for authentication services.
package service_auth

import (
	"fmt"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// UserLoginService orchestrates the authentication flow for existing users.
// It embeds BaseAuthService to reuse common security logic, adhering to the input.UserServiceLogin port. It is responsible for the transition from raw credentials to a secure, authenticated session.
type UserLoginService struct {
	BaseAuthService
}

// NewUserLoginService initializes a UserLoginService with all necessary
// infrastructure and domain dependencies.

// Parameters:
//   - userRepo: persistence adapter for user data.
//   - userNameValidator/passwordValidator: business rule engines for input integrity.
//   - csrfService/csrfCookieSetter: components for cross-site request forgery protection.
//
// Returns:
//   - input.UserServiceLogin: the abstracted login service interface.
func NewUserLoginService(userRepo output.UserRepository, userNameValidator, passwordValidator input.Validator, csrfService input.CSRFService, csrfCookieSetter output.CSRFCookieSetter) input.UserServiceLogin {
	return &UserLoginService{
		BaseAuthService: BaseAuthService{
			UserRepo:          userRepo,
			UserNameValidator: userNameValidator,
			PasswordValidator: passwordValidator,
			CSRFService:       csrfService,
			CSRFCookieSetter:  csrfCookieSetter,
		},
	}
}

// Login executes the full authentication protocol.
// The process follows a strict security sequence:
//  1. Input Validation: Ensures data follows domain formats before hitting the DB.
//  2. Identity Verification: Confirms the user exists.
//  3. Credential Challenge: Performs a secure bcrypt comparison against the stored hash.
//  4. Security Upgrading: Generates a new CSRF context for the authenticated session.
//  5. Token Issuance: Signs a JWT to be used for subsequent authorized requests.

// Returns the signed JWT string or a domain-specific error (Unauthorized, NotFound, etc.).
func (l *UserLoginService) Login(account models.Account) (string, error) {
	// 1. Validate format integrity
	if err := l.ValidateUserName(account.UserName); err != nil {
		return "", errors.NewValidationError(errors.ErrInvalidUsername)
	}

	// 2. Identify user in persistence
	exists, err := l.CheckUserExists(account.UserName)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", errors.NewNotFoundError(errors.ErrUserNotFound)
	}

	// 3. Retrieve security credentials
	storedHash, err := l.UserRepo.GetHashPassword(account.UserName)
	if err != nil {
		return "", err
	}

	userId, err := l.UserRepo.GetID(account.UserName)
	if err != nil {
		return "", err
	}

	// 4. Cryptographic verification
	// CompareHashAndPassword handles the complexity of constant-time comparisons to prevent timing attacks.
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(account.Password))
	if err != nil {
		return "", errors.NewAuthError(errors.ErrInvalidCredentials)
	}

	// 5. Establish CSRF Protection for the new session
	userIDStr := fmt.Sprintf("%d", userId)
	if err := l.GenerateAndSetCSRFToken(userIDStr); err != nil {
		return "", err
	}

	// 6. Finalize session via JWT
	token, err := l.GenerateToken(userId, account.UserName)
	if err != nil {
		return "", err
	}
	return token, nil
}
