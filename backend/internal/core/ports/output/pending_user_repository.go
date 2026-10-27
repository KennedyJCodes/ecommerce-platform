// Package output defines output port interfaces for the application.
package output

import (
	"context"
	"time"

	modelsdb "github.com/KennedyJCodes/ecommerce-platform/internal/core/domain/models/database"
)

// PendingUserRepository defines the persistence contract for users awaiting verification.
type PendingUserRepository interface {
	// Save stores a pending user as a Redis Hash with the provided expiration time.
	SavePendingUser(ctx context.Context, user *modelsdb.PendingUser, ttl time.Duration) error

	// GetPendingUserHashCode retrieves the stored verification-code hash for a pending user.
	GetPendingUserHashCode(ctx context.Context, userID string) (string, error)

	// GetPendingUser retrieves all data needed to complete a pending registration.
	GetPendingUser(ctx context.Context, userID string) (*modelsdb.PendingUser, error)

	// DeletePendingUser removes a pending registration after successful verification.
	DeletePendingUser(ctx context.Context, userID string) error

	// RegisterVerificationAttempt atomically consumes one verification attempt and
	// returns the attempt number. It fails if the pending user does not exist or
	// if the attempt limit was exceeded (in which case the pending user is deleted).
	RegisterVerificationAttempt(ctx context.Context, userID string, maxAttempts int) (int, error)
}
