// Package service_auth provides implementations of the core input ports for authentication flows.
package service_auth

import (
	"fmt"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
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
func NewUserRegisterService(userRepo output.UserRepository, userNameValidator, passwordValidator input.Validator, csrfService input.CSRFService, csrfCookieSetter output.CSRFCookieSetter) input.UserServiceRegister {
	return &UserRegisterService{
		BaseAuthService: BaseAuthService{
			UserRepo:          userRepo,
			UserNameValidator: userNameValidator,
			PasswordValidator: passwordValidator,
			CSRFService:       csrfService,
			CSRFCookieSetter:  csrfCookieSetter,
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
//  6. Session Issuance: Returns a signed JWT for immediate authentication.
//
// Returns a JWT token string on success, or an error if any step in the flow fails.
func (r *UserRegisterService) Register(account models.Account) (string, error) {
	// 1. Validate input integrity (Format and Strength)
	if err := r.ValidateUserName(account.UserName); err != nil {
		return "", errors.NewValidationError(errors.ErrInvalidUsername)
	}

	if err := r.ValidatePassword(account.Password); err != nil {
		return "", errors.NewValidationError(errors.ErrInvalidPassword)
	}

	// 2. Ensure username uniqueness (Conflict Check)
	exists, err := r.CheckUserExists(account.UserName)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errors.NewConflictError(errors.ErrUserAlreadyExists)
	}

	// 3. Persist the new identity
	// The repository is expected to handle the cryptographic hashing before storage.
	if err := r.UserRepo.SaveUser(account.UserName, account.Password); err != nil {
		return "", err
	}

	// 4. Resolve the persistent ID
	userId, err := r.UserRepo.GetID(account.UserName)
	if err != nil {
		return "", err
	}

	// 5. Establish immediate security session (CSRF)
	userIDStr := fmt.Sprintf("%d", userId)
	if err := r.GenerateAndSetCSRFToken(userIDStr); err != nil {
		return "", err
	}

	// 6. Generate access token
	return r.GenerateToken(userId, account.UserName)
}
