// Package repository provides SQL-based implementations of output ports for persisting and retrieving user data. It depends on a SQL database connection and pluggable security components for salt generation and password hashing.
package repository_mysql

import (
	"context"
	"database/sql"
	"errors"

	models_auth "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/auth"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	errorsApp "github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
	securityAuth "github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/security/security_auth"
	"github.com/jmoiron/sqlx"
)

// SQLUserRepository implements the UserRepository interface using a SQL database.

// It requires a *sqlx.DB for database operations, a Generator for creating salts, and a Hasher for hashing passwords.
type SQLUserRepository struct {
	db     *sqlx.DB
	hasher securityAuth.Hasher
}

// NewSQLUserRepository creates a new SQLUserRepository instance.

// It validates constructor dependencies and returns an error instead of
// terminating the process, leaving startup decisions to the composition root.
func NewSQLUserRepository(db *sqlx.DB, hasher securityAuth.Hasher) (output.UserRepository, error) {
	if db == nil {
		return nil, errorsApp.NewInternalError(errorsApp.ErrDatabaseConnection)
	}

	if hasher == nil {
		return nil, errorsApp.NewInternalError("Hasher not initialized")
	}

	return &SQLUserRepository{
		db:     db,
		hasher: hasher,
	}, nil
}

func (r *SQLUserRepository) FindByUserName(ctx context.Context, username string) (*models_auth.User, error) {
	var user models_auth.User
    query := `SELECT user_id, username, email, password, created_at, updated_at 
            FROM user_registration WHERE username = ?`

    err := r.db.QueryRowContext(ctx, query, username).Scan(
        &user.UserID, &user.UserName, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
    )
    if errors.Is(err, sql.ErrNoRows) {
        return &models_auth.User{}, errorsApp.NewBadRequestError(errorsApp.ErrUserNotFound)
    }
    if err != nil {
        return &models_auth.User{}, err
    }
    return &user, nil
}

// UserExists checks whether a user with the given username exists in the database.
// It returns true if a matching record is found, or false otherwise.
// Any SQL errors are wrapped as internal errors.
func (r *SQLUserRepository) UserExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM user_registration WHERE username = ?)"
	err := r.db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, errorsApp.NewInternalError(errorsApp.ErrDatabaseQuery).WithError(err)
	}

	return exists, nil
}

func (r *SQLUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM user_registration WHERE email = ?)"
	err := r.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, errorsApp.NewInternalError(errorsApp.ErrDatabaseQuery).WithError(err)
	}
	
	return exists, nil
}

// SaveUser inserts a new user into the database with a salted and hashed password.
// It generates a new salt, combines it with the plain password, hashes the result, and executes an INSERT statement. Any generation, hashing, or SQL errors are wrapped as internal errors.
func (r *SQLUserRepository) SaveUser(ctx context.Context, user models_auth.User) (models_auth.User, error) {
		tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return models_auth.User{}, errorsApp.NewInternalError(errorsApp.ErrDatabaseTransaction).WithError(err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx,
		"INSERT INTO user_registration (username, password, email) VALUES (?, ?, ?)",
		user.UserName, user.Password, user.Email,
	)
	if err != nil {
		return models_auth.User{}, errorsApp.NewInternalError(errorsApp.ErrDatabaseInsert).WithError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models_auth.User{}, errorsApp.NewInternalError(errorsApp.ErrDatabaseInsert).WithError(err)
	}
	user.UserID = uint64(id)

	if err = tx.Commit(); err != nil {
		return models_auth.User{}, errorsApp.NewInternalError(errorsApp.ErrDatabaseCommit).WithError(err)
	}

	return user, nil
}