package service_code_verification_test

import (
	"context"
	"errors"
	"testing"
	"time"

	models_auth "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/auth"
	modelsdb "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/database"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/services/service_code_verification"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/security/security_auth"
)

// ---------- Mocks ----------

type pendingUserRepositoryMock struct {
	pendingUser *modelsdb.PendingUser
	hashCode    string
	deleted     bool

	hashCodeErr error
	getUserErr  error
	deleteErr   error
}

func (m *pendingUserRepositoryMock) SavePendingUser(*modelsdb.PendingUser, time.Duration) error {
	return nil
}
func (m *pendingUserRepositoryMock) GetPendingUserHashCode(string) (string, error) {
	if m.hashCodeErr != nil {
		return "", m.hashCodeErr
	}
	return m.hashCode, nil
}
func (m *pendingUserRepositoryMock) GetPendingUser(string) (*modelsdb.PendingUser, error) {
	if m.getUserErr != nil {
		return nil, m.getUserErr
	}
	return m.pendingUser, nil
}
func (m *pendingUserRepositoryMock) DeletePendingUser(string) error {
	m.deleted = true
	return m.deleteErr
}

type userRepositoryMock struct {
	savedUser models_auth.User
	err       error
}

func (m *userRepositoryMock) FindByUserName(context.Context, string) (*models_auth.User, error) {
	return nil, nil
}
func (m *userRepositoryMock) UserExists(context.Context, string) (bool, error)  { return false, nil }
func (m *userRepositoryMock) EmailExists(context.Context, string) (bool, error) { return false, nil }
func (m *userRepositoryMock) SaveUser(_ context.Context, user models_auth.User) (models_auth.User, error) {
	m.savedUser = user
	if m.err != nil {
		return models_auth.User{}, m.err
	}
	user.UserID = 42
	return user, nil
}

type tokenServiceMock struct {
	tokens []models_auth.TokenType
}

func (m *tokenServiceMock) GenerateToken(_ int, _ string, tokenType models_auth.TokenType) (string, error) {
	m.tokens = append(m.tokens, tokenType)
	if tokenType == models_auth.TokenTypeAccess {
		return "access-token", nil
	}
	return "refresh-token", nil
}
func (m *tokenServiceMock) ValidateToken(string) (*models_auth.Claims, error) { return nil, nil }

type csrfServiceMock struct{ userID string }

func (m *csrfServiceMock) GenerateToken(userID string) (string, error) {
	m.userID = userID
	return "csrf-token", nil
}
func (m *csrfServiceMock) ValidateToken(string, string) error { return nil }

// ---------- Helpers ----------

func newVerificationService(t *testing.T, pending *pendingUserRepositoryMock, users *userRepositoryMock, tokens *tokenServiceMock, csrf *csrfServiceMock) *service_code_verification.UserVerificationService {
	t.Helper()
	hasher := security_auth.BcryptHasher{}
	service := service_code_verification.NewUserVerificationService(users, pending, tokens, csrf, hasher)
	return service.(*service_code_verification.UserVerificationService)
}

func hashOf(t *testing.T, code string) string {
	t.Helper()
	h, err := (security_auth.BcryptHasher{}).Hash([]byte(code))
	if err != nil {
		t.Fatal(err)
	}
	return h
}

// ---------- Test ----------

func TestUserVerificationService_Verify(t *testing.T) {
	validHash := hashOf(t, "123456")

	tests := []struct {
		name        string
		code        string
		hashCodeErr error 
		getUserErr  error
		saveErr     error
		wantErr     bool
		wantSaved   bool
		wantDeleted bool
	}{
		{
			name:        "success",
			code:        "123456",
			wantSaved:   true,
			wantDeleted: true,
		},
		{
			name:    "incorrect code",
			code:    "000000",
			wantErr: true,
		},
		{
			name:        "hash code not found or expired",
			code:        "123456",
			hashCodeErr: errors.New("pending code not found"),
			wantErr:     true,
		},
		{
			name:       "get pending user fails",
			code:       "123456",
			getUserErr: errors.New("pending user not found"),
			wantErr:    true,
		},
		{
			name:      "save error",
			code:      "123456",
			saveErr:   errors.New("database unavailable"),
			wantErr:   true,
			wantSaved: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pending := &pendingUserRepositoryMock{
				hashCode:    validHash,
				hashCodeErr: tt.hashCodeErr,
				getUserErr:  tt.getUserErr,
				pendingUser: &modelsdb.PendingUser{
					Username:     "alice",
					Email:        "alice@example.com",
					PasswordHash: "password-hash",
				},
			}
			users := &userRepositoryMock{err: tt.saveErr}
			tokens := &tokenServiceMock{}
			csrf := &csrfServiceMock{}
			svc := newVerificationService(t, pending, users, tokens, csrf)

			gotTokens, gotCSRF, err := svc.Verify(context.Background(), "pending-1", tt.code)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
			if saved := users.savedUser.UserName != ""; saved != tt.wantSaved {
				t.Errorf("user saved = %v, want %v", saved, tt.wantSaved)
			}
			if pending.deleted != tt.wantDeleted {
				t.Errorf("pending deleted = %v, want %v", pending.deleted, tt.wantDeleted)
			}
			if tt.wantErr {
				return
			}


			if gotTokens.AccessToken != "access-token" || gotTokens.RefreshToken != "refresh-token" || gotCSRF != "csrf-token" {
				t.Errorf("unexpected tokens: %+v, %q", gotTokens, gotCSRF)
			}
			if users.savedUser.Email != "alice@example.com" || users.savedUser.Password != "password-hash" || users.savedUser.VerifiedAt.IsZero() {
				t.Errorf("saved user = %+v", users.savedUser)
			}
			if csrf.userID != "42" || len(tokens.tokens) != 2 {
				t.Errorf("csrfUser=%q tokens=%v", csrf.userID, tokens.tokens)
			}
		})
	}
}

var _ output.PendingUserRepository = (*pendingUserRepositoryMock)(nil)