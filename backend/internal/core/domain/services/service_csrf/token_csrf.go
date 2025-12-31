package service_csrf

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
)

type CSRFUseCase struct {
	repository output.CSRFRepository
	timeToLive time.Duration
}

func NewCSRFUseCase(repo output.CSRFRepository, timeToLive time.Duration) *CSRFUseCase {
	return &CSRFUseCase{
		repository: repo,
		timeToLive: timeToLive,
	}
}

func (uc *CSRFUseCase) GenerateToken(userID string) (string, error) {
	if userID == "" {
		return "", errors.NewValidationError(errors.ErrEmptyField)
	}

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

func (uc *CSRFUseCase) ValidateToken(tokenValue string, userID string) error {
	token, err := uc.repository.Find(userID)
	if err != nil {
		return errors.NewNotFoundError(errors.ErrCSRFTokenNotFound)
	}

	if token.Value != tokenValue {
        return errors.NewNotFoundError(errors.ErrInvalidCSRFToken)
    }

	if !token.IsValid() {
		uc.repository.Delete(userID)
		return errors.NewNotFoundError(errors.ErrCSRFTokenExpired)
	}

	return nil
}
