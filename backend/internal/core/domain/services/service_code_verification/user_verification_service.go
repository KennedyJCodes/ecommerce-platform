package service_code_verification

import (
	"context"
	"strconv"
	"time"

	models_auth "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/auth"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	appErrors "github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/security/security_auth"
)

// UserVerificationService validates a pending user's code and completes registration.
type UserVerificationService struct {
	userRepository        output.UserRepository
	pendingUserRepository output.PendingUserRepository
	tokenService          output.TokenService
	csrfService           output.CSRFService
	hasher                security_auth.Hasher
}

var _ input.UserVerificationService = (*UserVerificationService)(nil)

// NewUserVerificationService creates the service used by the /verify endpoint.
func NewUserVerificationService(
	userRepository output.UserRepository,
	pendingUserRepository output.PendingUserRepository,
	tokenService output.TokenService,
	csrfService output.CSRFService,
	hasher security_auth.Hasher,
) input.UserVerificationService {
	return &UserVerificationService{
		userRepository:        userRepository,
		pendingUserRepository: pendingUserRepository,
		tokenService:          tokenService,
		csrfService:           csrfService,
		hasher:                hasher,
	}
}

// Verify compares the submitted code, saves the verified user, and creates its session tokens.
func (s *UserVerificationService) Verify(ctx context.Context, pendingUserID string, code string) (*models_auth.TokenPair, string, error) {
	if pendingUserID == "" || code == "" {
		return nil, "", appErrors.NewValidationError(appErrors.ErrInvalidRequest)
	}

	hashCode, err := s.pendingUserRepository.GetPendingUserHashCode(pendingUserID)
	if err != nil {
		return nil, "", err
	}
	if err := s.hasher.Compare([]byte(code), hashCode); err != nil {
		return nil, "", appErrors.NewAuthError(appErrors.ErrInvalidCredentials)
	}

	pendingUser, err := s.pendingUserRepository.GetPendingUser(pendingUserID)
	if err != nil {
		return nil, "", err
	}

	user, err := s.userRepository.SaveUser(ctx, models_auth.User{
		UserName:   pendingUser.Username,
		Email:      pendingUser.Email,
		Password:   pendingUser.PasswordHash,
		VerifiedAt: time.Now(),
	})
	if err != nil {
		return nil, "", err
	}

	if err := s.pendingUserRepository.DeletePendingUser(pendingUserID); err != nil {
		return nil, "", err
	}

	accessToken, err := s.tokenService.GenerateToken(int(user.UserID), user.UserName, models_auth.TokenTypeAccess)
	if err != nil {
		return nil, "", appErrors.NewInternalError(appErrors.ErrTokenGeneration).WithError(err)
	}
	refreshToken, err := s.tokenService.GenerateToken(int(user.UserID), user.UserName, models_auth.TokenTypeRefresh)
	if err != nil {
		return nil, "", appErrors.NewInternalError(appErrors.ErrTokenGeneration).WithError(err)
	}

	csrfToken, err := s.csrfService.GenerateToken(strconv.FormatUint(user.UserID, 10))
	if err != nil {
		return nil, "", appErrors.NewInternalError(appErrors.ErrTokenGeneration).WithError(err)
	}

	return &models_auth.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, csrfToken, nil
}
