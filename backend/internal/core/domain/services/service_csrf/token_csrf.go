// Package service_csrf implements the security logic for Cross-Site Request Forgery protection.
package service_csrf

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
)

// CSRFUseCase orchestrates the lifecycle of anti-forgery tokens.
// It manages secure token generation, persistence with TTL (Time To Live), and integrity validation to prevent unauthorized cross-origin requests.
type CSRFUseCase struct {
	repository output.CSRFRepository
	timeToLive time.Duration
}

// NewCSRFUseCase initializes the service with a repository for token storage and a duration defining the validity window of each token.
func NewCSRFUseCase(repo output.CSRFRepository, timeToLive time.Duration) *CSRFUseCase {
	return &CSRFUseCase{
		repository: repo,
		timeToLive: timeToLive,
	}
}

// GenerateToken creates a cryptographically secure random token for a specific user.
// The token is encoded in Base64 (URL Safe) and persisted with an expiration timestamp.
func (uc *CSRFUseCase) GenerateToken(userID string) (string, error) {
	if userID == "" {
		return "", errors.NewValidationError(errors.ErrEmptyField)
	}

	// Uses a cryptographically secure random number generator (CSPRNG)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", errors.NewInternalError(errors.ErrTokenGeneration).WithError(err)
	}

	tokenValue := base64.RawURLEncoding.EncodeToString(bytes)

	now := time.Now()
	token := &models.CSRFToken{
		Value:     tokenValue,
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: now.Add(uc.timeToLive),
	}

	if err := uc.repository.Save(token); err != nil {
		return "", errors.NewInternalError(errors.ErrTokenGeneration).WithError(err)
	}

	return tokenValue, nil
}

// ValidateToken verifies if the provided token matches the stored value for the user.
// It also checks for expiration; if the token is expired, it is removed from the repository.
func (uc *CSRFUseCase) ValidateToken(tokenValue string, userID string) error {
	token, err := uc.repository.Find(userID)
	if err != nil {
		return errors.NewNotFoundError(errors.ErrCSRFTokenNotFound)
	}

	if token.Value != tokenValue {
		return errors.NewNotFoundError(errors.ErrInvalidCSRFToken)
	}

	// Check if token has exceeded its TTL
	if !token.IsValid() {
		uc.repository.Delete(userID) // Cleanup expired token
		return errors.NewNotFoundError(errors.ErrCSRFTokenExpired)
	}

	return nil
}
