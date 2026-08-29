// Package service_auth provides implementations of the core input ports for authentication flows.
package service_auth

import (
	"context"
	"fmt"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/dto/auth"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models"
	models_auth "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/auth"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/security/security_auth"
)

// UserRegisterService handles the logic for creating new user identities.
// It leverages BaseAuthService for shared validation and token orchestration, ensuring that new accounts meet the same security standards as existing ones.
type UserRegisterService struct {
	BaseAuthService
}

// NewUserRegisterService constructs a UserRegisterService with the necessary collaborators for persistence, validation, and security session management.

// Parameters:
//   - userRepo: output port for user persistence.
//   - userNameValidator/passwordValidator: components to enforce domain constraints.
//   - csrfService/csrfCookieSetter: infrastructure for CSRF protection.
//
// Returns:
//   - input.UserServiceRegister: the registration service interface.
func NewUserRegisterService(userRepo output.UserRepository, userNameValidator, passwordValidator input.Validator, emailValidator input.Validator, tokenService output.TokenService, csrfService output.CSRFService, passwordHasher security_auth.Hasher) input.UserServiceRegister {
	return &UserRegisterService{
		BaseAuthService: BaseAuthService{
			UserRepo:          userRepo,
			UserNameValidator: userNameValidator,
			PasswordValidator: passwordValidator,
			EmailValidator:    emailValidator,
			TokenService:      tokenService,
			CSRFService:       csrfService,
			Hasher:            passwordHasher,
		},
	}
}

// Register orchestrates the onboarding of a new user into the system.

// The registration flow is executed as an atomic-like business process:
//  1. Format Validation: Verifies username and password against complexity rules.
//  2. Collision Check: Prevents duplicate accounts by verifying username uniqueness.
//  3. Secure Persistence: Delegates the hashing and storage to the repository layer.
//  4. Identity Resolution: Retrieves the newly generated unique ID.
//  5. Security Context: Initializes the CSRF state for the new user.
//  6. Session Issuance: Signs access and refresh JWTs for immediate authentication.
//
// Returns a TokenPair, a CSRF token on success, or an error.
func (r *UserRegisterService) Register(ctx context.Context, request dto.RegisterAccount) (*models.TokenPair, string, error) {
	if err := r.ValidateUserName(request.UserName); err != nil {
		return nil, "", err
	}
	
	if err := r.ValidatePassword(request.Password); err != nil {
		return nil, "", err
	}

	if err := r.ValidateEmail(request.Email); err != nil {
		return nil, "", err
	}
	
	existsEmail, err := r.CheckEmailExists(ctx, request.Email)
	if err != nil {
		return nil, "", err
	}
	if existsEmail {
		return nil, "", errors.NewConflictError(errors.ErrEmailAlreadyExists)
	}
	
	existsUser, err := r.CheckUserExists(ctx, request.UserName)
	if err != nil {
		return nil, "", err
	}
	if existsUser {
		return nil, "", errors.NewConflictError(errors.ErrUserAlreadyExists)
	}
	
	hash, err := r.Hasher.Hash([]byte(request.Password))
	if err != nil {
		return nil, "", errors.NewInternalError(errors.ErrHashingPassword).WithError(err)
	}

	newUser := models_auth.User{
		UserName:     request.UserName,
		Email:        request.Email,
		PasswordHash: string(hash),
	}

	// The repository is expected to handle the cryptographic hashing before storage.
	savedUser, err := r.UserRepo.SaveUser(ctx, newUser); 
	if err != nil {
		return nil, "", err
	}

	userIDStr := fmt.Sprintf("%d", savedUser.UserID)
	csrfToken, err := r.GenerateCSRFToken(userIDStr)
	if err != nil {
		return nil, "", err
	}

	tokens, err := r.GenerateTokenPair(int(savedUser.UserID), request.UserName)
	if err != nil {
		return nil, "", err
	}
	return tokens, csrfToken, nil
}
