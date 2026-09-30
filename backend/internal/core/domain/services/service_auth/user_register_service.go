package service_auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/dto/auth"
	modelsdb "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/database"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/security/security_auth"
)

type UserRegisterService struct{ BaseAuthService }

func NewUserRegisterService(userRepo output.UserRepository, pendingUserRepository output.PendingUserRepository, userNameValidator, passwordValidator input.Validator, emailValidator input.Validator, hasher security_auth.Hasher, codeVerificationService input.CodeVerificationService, codeVerificationSender output.CodeVerificationSender) input.UserServiceRegister {
	return &UserRegisterService{BaseAuthService: BaseAuthService{
		UserRepo: userRepo, PendingUserRepository: pendingUserRepository,
		UserNameValidator: userNameValidator, PasswordValidator: passwordValidator,
		EmailValidator: emailValidator, Hasher: hasher,
		CodeVerificationService: codeVerificationService, CodeVerificationSender: codeVerificationSender,
	}}
}

func (r *UserRegisterService) Register(ctx context.Context, request dto.RegisterAccount) error {
	if err := r.ValidateUserName(request.UserName); err != nil {
		return err
	}
	if err := r.ValidatePassword(request.Password); err != nil {
		return err
	}
	if err := r.ValidateEmail(request.Email); err != nil {
		return err
	}

	existsUser, err := r.CheckUserExists(ctx, request.UserName)
	if err != nil {
		return err
	}
	if existsUser {
		return errors.NewConflictError(errors.ErrUserAlreadyExists)
	}
	existsEmail, err := r.CheckEmailExists(ctx, request.Email)
	if err != nil {
		return err
	}
	if existsEmail {
		return errors.NewConflictError(errors.ErrEmailAlreadyExists)
	}

	code, err := r.CodeVerificationService.GenerateCodeVerification()
	if err != nil {
		return errors.NewInternalError(errors.ErrGeneratingCodeVerification).WithError(err)
	}
	passwordHash, err := r.HashSensitiveValue([]byte(request.Password))
	if err != nil {
		return errors.NewInternalError(errors.ErrHashingPassword).WithError(err)
	}
	codeHash, err := r.HashSensitiveValue([]byte(code))
	if err != nil {
		return errors.NewInternalError(errors.ErrGeneratingCodeVerification).WithError(err)
	}

	userID, err := newPendingUserID()
	if err != nil {
		return errors.NewInternalError(errors.ErrDatabaseInsert).WithError(err)
	}
	pendingUser := &modelsdb.PendingUser{
		ID: userID, Username: request.UserName, PasswordHash: string(passwordHash),
		Email: request.Email, HashCode: string(codeHash), Attempts: 0,
	}
	if err := r.PendingUserRepository.SavePendingUser(pendingUser, 10*time.Minute); err != nil {
		return err
	}

	return r.SendCodeVerification(request.Email, code)
}

func newPendingUserID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
