// Package input defines service contracts for business logic operations.
// This file contains the TokenService interface for JWT-related operations.
package output

import "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models"

// TokenService defines the contract for JWT generation, validation, and parsing.
// Implementations handle the lifecycle of access tokens and refresh tokens used
// throughout the authentication flow.
type TokenService interface {
	// GenerateJWT creates a signed access JWT for the given user.
	// Returns the token string or an error if signing fails.
	GenerateToken(userId int, userName string, tokenType models.TokenType) (string, error)

	// ValidateRefreshToken parses and validates a refresh token string.
	// Returns the claims if the token is valid and has subject "refresh".
	ValidateToken(tokenString string) (*models.Claims, error)
}
