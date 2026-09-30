// Package output defines output port interfaces for the application.
package output

import (
	"time"

	modelsdb "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/database"
)

// PendingUserRepository defines the persistence contract for users awaiting verification.
type PendingUserRepository interface {
	// Save stores a pending user as a Redis Hash with the provided expiration time.
	SavePendingUser(user *modelsdb.PendingUser, ttl time.Duration) error
}
