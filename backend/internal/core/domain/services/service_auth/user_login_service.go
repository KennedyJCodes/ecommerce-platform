// Package service_auth provides implementations of input port interfaces for authentication services.
package service_auth

import (
	"context"
	"fmt"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/dto/auth"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/auth"
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
func NewUserLoginService(userRepo output.UserRepository, userNameValidator, passwordValidator input.Validator, tokenService output.TokenService, csrfService output.CSRFService) input.UserServiceLogin {
	return &UserLoginService{
		BaseAuthService: BaseAuthService{
			UserRepo:          userRepo,
			UserNameValidator: userNameValidator,
			PasswordValidator: passwordValidator,
			TokenService:      tokenService,
			CSRFService:       csrfService,
		},
	}
}

// Login executes the full authentication protocol.
// The process follows a strict security sequence:
//  1. Input Validation: Ensures data follows domain formats before hitting the DB.
//  2. Identity Verification: Confirms the user exists.
//  3. Credential Challenge: Performs a secure bcrypt comparison against the stored hash.
//  4. Security Upgrading: Generates a new CSRF context for the authenticated session.
//  5. Token Issuance: Signs access and refresh JWTs for subsequent authorized requests.

// Returns a TokenPair containing both tokens and a CSRF token, or a domain-specific error.
func (l *UserLoginService) Login(ctx context.Context, request dto.LoginRequest) (*models_auth.TokenPair, string, error) {
	if err := l.ValidateUserName(request.UserName); err != nil {
		return nil, "", errors.NewValidationError(errors.ErrInvalidUsername)
	}

	user, err := l.UserRepo.FindByUserName(ctx, request.UserName)
	if err != nil {
		return nil, "", errors.NewAuthError(errors.ErrInvalidCredentials)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, "", errors.NewAuthError(errors.ErrInvalidCredentials)
	}

	userIDStr := fmt.Sprintf("%d", user.UserID)
	csrfToken, err := l.GenerateCSRFToken(userIDStr)
	if err != nil {
		return nil, "", err
	}

	tokens, err := l.GenerateTokenPair(int(user.UserID), request.UserName)
	if err != nil {
		return nil, "", err
	}
	return tokens, csrfToken, nil
}
