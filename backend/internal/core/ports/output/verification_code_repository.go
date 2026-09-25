// Package output defines output port interfaces for the application.
package output

import (
	"context"
	"time"
)

// VerificationCodeRepository defines the persistence contract for user
// verification code hashes.
type VerificationCodeRepository interface {
	// SaveCodeVerication persists an already-hashed verification code for a user
	// until it expires.
	SaveCodeVerication(ctx context.Context, userID uint64, codeHash string, expiresAt time.Time) error
}
