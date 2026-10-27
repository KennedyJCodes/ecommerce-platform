package service_auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/KennedyJCodes/ecommerce-platform/internal/core/domain/dto/auth"
	modelsdb "github.com/KennedyJCodes/ecommerce-platform/internal/core/domain/models/database"
	"github.com/KennedyJCodes/ecommerce-platform/internal/core/ports/input"
	"github.com/KennedyJCodes/ecommerce-platform/internal/core/ports/output"
	"github.com/KennedyJCodes/ecommerce-platform/pkg/errors"
	"github.com/KennedyJCodes/ecommerce-platform/pkg/security/security_auth"
)

type UserRegisterService struct{ BaseAuthService }

const pendingUserTTL = 10 * time.Minute

// UserRegisterServiceDependencies groups the collaborators required by the
// user registration use case. Keeping this explicit makes the constructor
// easier to read and safer to extend when a new dependency is added.
type UserRegisterServiceDependencies struct {
	UserRepo                output.UserRepository
	PendingUserRepository   output.PendingUserRepository
	UserNameValidator       input.Validator
	PasswordValidator       input.Validator
	EmailValidator          input.Validator
	Hasher                  security_auth.Hasher
	CodeVerificationService input.CodeVerificationService
	CodeVerificationSender  output.CodeVerificationSender
}

func NewUserRegisterService(deps UserRegisterServiceDependencies) input.UserServiceRegister {
	return &UserRegisterService{BaseAuthService: BaseAuthService{
		UserRepo:                deps.UserRepo,
		PendingUserRepository:   deps.PendingUserRepository,
		UserNameValidator:       deps.UserNameValidator,
		PasswordValidator:       deps.PasswordValidator,
		EmailValidator:          deps.EmailValidator,
		Hasher:                  deps.Hasher,
		CodeVerificationService: deps.CodeVerificationService,
		CodeVerificationSender:  deps.CodeVerificationSender,
	}}
}

func (r *UserRegisterService) Register(ctx context.Context, request dto.RegisterAccount) (string, error) {
	if err := r.ValidateUserName(request.UserName); err != nil {
		return "", err
	}
	if err := r.ValidatePassword(request.Password); err != nil {
		return "", err
	}
	if err := r.ValidateEmail(request.Email); err != nil {
		return "", err
	}

	existsUser, err := r.CheckUserExists(ctx, request.UserName)
	if err != nil {
		return "", err
	}
	if existsUser {
		return "", errors.NewConflictError(errors.ErrUserAlreadyExists)
	}
	existsEmail, err := r.CheckEmailExists(ctx, request.Email)
	if err != nil {
		return "", err
	}
	if existsEmail {
		return "", errors.NewConflictError(errors.ErrEmailAlreadyExists)
	}

	code, err := r.CodeVerificationService.GenerateCodeVerification()
	if err != nil {
		return "", errors.NewInternalError(errors.ErrGeneratingCodeVerification).WithError(err)
	}
	passwordHash, err := r.HashSensitiveValue([]byte(request.Password))
	if err != nil {
		return "", errors.NewInternalError(errors.ErrHashingPassword).WithError(err)
	}
	codeHash, err := r.HashSensitiveValue([]byte(code))
	if err != nil {
		return "", errors.NewInternalError(errors.ErrGeneratingCodeVerification).WithError(err)
	}

	userID, err := newPendingUserID()
	if err != nil {
		return "", errors.NewInternalError(errors.ErrDatabaseInsert).WithError(err)
	}
	pendingUser := &modelsdb.PendingUser{
		ID: userID, Username: request.UserName, PasswordHash: string(passwordHash),
		Email: request.Email, HashCode: string(codeHash), Attempts: 0,
	}
	if err := r.PendingUserRepository.SavePendingUser(ctx, pendingUser, pendingUserTTL); err != nil {
		return "", err
	}

	if err := r.SendCodeVerification(request.Email, code); err != nil {
		return "", err
	}

	return userID, nil
}

func newPendingUserID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
