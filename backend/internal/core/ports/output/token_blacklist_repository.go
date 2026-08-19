// Package output defines output port interfaces for the application.
// This file contains the TokenBlacklistPort interface for managing revoked JWT tokens.
package output

import "time"

// TokenBlacklistPort defines the contract for persisting and querying revoked JWT tokens.
// Implementations handle the storage and retrieval of blacklisted token identifiers (jti),
// enabling server-side token revocation on logout or security events.
type TokenBlacklistPort interface {
	// Add stores a token identifier in the blacklist with an associated TTL.
	// The token will be considered revoked until the TTL expires.
	// Parameters:
	//   - jti: the unique JWT ID (token identifier) to blacklist.
	//   - ttl: the duration for which the token should remain blacklisted.
	// Returns:
	//   - error: non-nil if the storage operation fails.
	Add(jti string, ttl time.Duration) error

	// IsBlacklisted checks whether a given token identifier has been revoked.
	// Parameters:
	//   - jti: the unique JWT ID to check against the blacklist.
	// Returns:
	//   - bool: true if the token is blacklisted (revoked).
	//   - error: non-nil if the query operation fails.
	IsBlacklisted(jti string) (bool, error)
}
