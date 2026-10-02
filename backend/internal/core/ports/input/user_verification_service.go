package input

import (
	"context"

	models_auth "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/auth"
)

// UserVerificationService completes a pending registration using its verification code.
type UserVerificationService interface {
	Verify(ctx context.Context, pendingUserID string, code string) (*models_auth.TokenPair, string, error)
}
