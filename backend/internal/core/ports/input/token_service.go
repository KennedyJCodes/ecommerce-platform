// Package input defines service contracts for business logic operations.
// This file contains the TokenService interface for JWT-related operations.
package input

import "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models"

// TokenService defines the contract for JWT generation, validation, and parsing.
// Implementations handle the lifecycle of access tokens and refresh tokens used
// throughout the authentication flow.
type TokenService interface {
	// GenerateJWT creates a signed access JWT for the given user.
	// Returns the token string or an error if signing fails.
	GenerateJWT(userId int, userName string) (string, error)

	// GenerateRefreshToken creates a signed refresh JWT for the given user.
	// Refresh tokens have a longer expiration than access tokens.
	GenerateRefreshToken(userID int, userName string) (string, error)

	// ValidateRefreshToken parses and validates a refresh token string.
	// Returns the claims if the token is valid and has subject "refresh".
	ValidateRefreshToken(tokenString string) (*models.Claims, error)

	// ParseTokenWithClaims validates a JWT string and extracts its custom claims.
	// Returns the claims if the token signature and expiration are valid.
	ParseTokenWithClaims(tokenString string) (*models.Claims, error)
}
