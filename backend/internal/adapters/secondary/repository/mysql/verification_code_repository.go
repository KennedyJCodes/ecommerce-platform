// Package repository_mysql provides SQL-based implementations of output ports
// for persisting data in MySQL.
package repository_mysql

import (
	"context"
	"time"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	appErrors "github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
	"github.com/jmoiron/sqlx"
)

// SQLVerificationCodeRepository persists hashed user verification codes.
type SQLVerificationCodeRepository struct {
	db *sqlx.DB
}

var _ output.VerificationCodeRepository = (*SQLVerificationCodeRepository)(nil)

// NewSQLVerificationCodeRepository creates a repository backed by MySQL.
func NewSQLVerificationCodeRepository(db *sqlx.DB) (output.VerificationCodeRepository, error) {
	if db == nil {
		return nil, appErrors.NewInternalError(appErrors.ErrDatabaseConnection)
	}

	return &SQLVerificationCodeRepository{db: db}, nil
}

// SaveCodeVerication stores a verification code hash for a user until its expiration time.
// The codeHash argument must already be hashed by the caller.
func (r *SQLVerificationCodeRepository) SaveCodeVerication(
	ctx context.Context,
	userID uint64,
	codeHash string,
	expiresAt time.Time,
) error {
	const query = `
		INSERT INTO user_verification_codes (user_id, code_hash, expires_at)
		VALUES (?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query, userID, codeHash, expiresAt)
	if err != nil {
		return appErrors.NewInternalError(appErrors.ErrDatabaseInsert).WithError(err)
	}

	return nil
}
