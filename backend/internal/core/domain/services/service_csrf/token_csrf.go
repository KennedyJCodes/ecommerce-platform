// Package service_csrf implements the security logic for Cross-Site Request Forgery protection.
package service_csrf

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/KennedyJCodes/ecommerce-platform/internal/core/domain/models/security"
	"github.com/KennedyJCodes/ecommerce-platform/internal/core/ports/output"
	"github.com/KennedyJCodes/ecommerce-platform/pkg/errors"
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
func (uc *CSRFUseCase) GenerateToken(ctx context.Context, userID string) (string, error) {
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
	token := &models_security.CSRFToken{
		Value:     tokenValue,
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: now.Add(uc.timeToLive),
	}

	if err := uc.repository.Save(ctx, token); err != nil {
		return "", errors.NewInternalError(errors.ErrTokenGeneration).WithError(err)
	}

	return tokenValue, nil
}

// ValidateToken verifies if the provided token matches the stored value for the user.
// It also checks for expiration; if the token is expired, it is removed from the repository.
// On success, the token is consumed (deleted) to enforce single-use policy.
func (uc *CSRFUseCase) ValidateToken(ctx context.Context, tokenValue string, userID string) error {
	token, err := uc.repository.Find(ctx, userID)
	if err != nil {
		return errors.NewNotFoundError(errors.ErrCSRFTokenNotFound)
	}

	if token.Value != tokenValue {
		return errors.NewNotFoundError(errors.ErrInvalidCSRFToken)
	}

	// Check if token has exceeded its TTL
	if !token.IsValid() {
		_ = uc.repository.Delete(ctx, userID)
		return errors.NewNotFoundError(errors.ErrCSRFTokenExpired)
	}

	// Consume the token (single-use policy)
	if err := uc.repository.Delete(ctx, userID); err != nil {
		return errors.NewInternalError(errors.ErrCSRFDelete).WithError(err)
	}

	return nil
}
