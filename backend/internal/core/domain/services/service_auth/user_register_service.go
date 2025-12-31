// Package service_auth provides implementations of the core input ports for authentication flows.
package service_auth

import (
	"fmt"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
)

// UserRegisterService implements the input.UserServiceRegister interface.

// It uses BaseAuthService for shared validation and token generation logic, and orchestrates the full registration flow.
type UserRegisterService struct {
	BaseAuthService
}

// NewUserRegisterService constructs a UserRegisterService with its dependencies.

// Parameters:
//   - userRepo: repository for persisting and querying user data.
//   - userNameValidator: enforces rules on username formats.
//   - passwordValidator: enforces rules on password strength.

//   - csrfService: service for CSRF token generation (input.CSRFService)
//   - csrfCookieSetter: setter for CSRF token cookies (output.CSRFCookieSetter)
//
// Returns:
//   - input.UserServiceRegister: the initialized registration service.
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

// Register processes a new user registration.

// It performs the following steps in order:
//   1. ValidateUserName – ensures the username meets formatting rules.
//   2. ValidatePassword – ensures the password meets strength rules.
//   3. CheckUserExists – returns a ConflictError if the username is already taken.
//   4. SaveUser       – persists the new username and password (with salt and hash).
//   5. GenerateToken – issues a JWT token for the newly created user.

// Parameters:
//   - account: models.Account containing UserName and Password.

// Returns:
//   - string: a signed JWT token upon successful registration.
//   - error: non‑nil if any validation, conflict, or persistence error occurs.
func (r *UserRegisterService) Register(account models.Account) (string, error) {
	// 1. Validate username format
	if err := r.ValidateUserName(account.UserName); err != nil {
		return "", errors.NewValidationError(errors.ErrInvalidUsername)
	}

	// 2. Validate password strength
	if err := r.ValidatePassword(account.Password); err != nil {
		return "", errors.NewValidationError(errors.ErrInvalidPassword)
	}

	// 3. Ensure the user does not already exist
	exists, err := r.CheckUserExists(account.UserName)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errors.NewConflictError(errors.ErrUserAlreadyExists)
	}

	// 4. Persist the new user with salted+hashed password
	if err := r.UserRepo.SaveUser(account.UserName, account.Password); err != nil {
		return "", err
	}

	userId, err := r.UserRepo.GetID(account.UserName)
	if err != nil {
		return "", err
	}

	// 5. Generate and set CSRF token
	userIDStr := fmt.Sprintf("%d", userId)
	if err := r.GenerateAndSetCSRFToken(userIDStr); err != nil {
		return "", err
	}

	// 6. Issue a JWT token for the new user
	return r.GenerateToken(userId, account.UserName)
}
