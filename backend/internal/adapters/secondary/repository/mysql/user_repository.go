// Package repository provides SQL-based implementations of output ports for persisting and retrieving user data. It depends on a SQL database connection and pluggable security components for salt generation and password hashing.
package repository_mysql

import (
	"database/sql"
	"log"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
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

// It logs a fatal error if any dependency is nil, ensuring that the repository always has a valid database connection, salt generator, and hasher.
// Returns an output.UserRepository ready for use.
func NewSQLUserRepository(db *sqlx.DB, hasher securityAuth.Hasher) output.UserRepository {
	if db == nil {
		log.Fatal(errors.NewInternalError(errors.ErrDatabaseConnection).Error())
	}

	if hasher == nil {
		log.Fatal(errors.NewInternalError("Hasher not initialized").Error())
	}

	log.Println("NewSQLUserRepository() is running successfully")

	return &SQLUserRepository{
		db:     db,
		hasher: hasher,
	}
}

// UserExists checks whether a user with the given username exists in the database.

// It returns true if a matching record is found, or false otherwise.
// Any SQL errors are wrapped as internal errors.
func (r *SQLUserRepository) UserExists(username string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM User_Registration WHERE UserName = ?)"
	err := r.db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, errors.NewInternalError(errors.ErrDatabaseQuery).WithError(err)
	}

	return exists, nil
}

// GetHashPassword retrieves the hashed password for the specified username.

// If no record is found, returns a NotFoundError. Other SQL errors are wrapped as internal errors.
func (r *SQLUserRepository) GetHashPassword(username string) (string, error) {
	var hashPassword string
	query := "SELECT Password FROM User_Registration WHERE UserName = ?"
	err := r.db.QueryRow(query, username).Scan(&hashPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.NewNotFoundError(errors.ErrUserNotFound)
		}
		return "", errors.NewInternalError(errors.ErrDatabaseQuery).WithError(err)
	}

	return hashPassword, nil
}

// SaveUser inserts a new user into the database with a salted and hashed password.

// It generates a new salt, combines it with the plain password, hashes the result, and executes an INSERT statement. Any generation, hashing, or SQL errors are wrapped as internal errors.
func (r *SQLUserRepository) SaveUser(username, password string) error {
	hash, err := r.hasher.Hash([]byte(password))
	if err != nil {
		return errors.NewInternalError(errors.ErrHashingPassword).WithError(err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return errors.NewInternalError(errors.ErrDatabaseTransaction).WithError(err)
	}

	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO User_Registration (UserName, Password) VALUES (?, ?)", username, hash)
	if err != nil {
		return errors.NewInternalError(errors.ErrDatabaseInsert).WithError(err)
	}

	if err = tx.Commit(); err != nil {
		return errors.NewInternalError(errors.ErrDatabaseCommit).WithError(err)
	}

	return nil
}

// GetID retrieves the unique user ID for a given username from the database.
// It returns a NotFoundError if no matching user is found, or an InternalError if any other database error occurs.

// Parameters:
//   - username: the username string to look up in the User_Registration table.

// Returns:
//   - int: the UserID corresponding to the provided username.
//   - error: non-nil if the user is not found or a database error occurs.
func (r *SQLUserRepository) GetID(username string) (int, error) {
	var id int
	query := "SELECT UserID FROM User_Registration WHERE UserName = ?"

	// Execute the query and scan the single result into id.
	err := r.db.QueryRow(query, username).Scan(&id)
	if err != nil {
		// No matching row indicates user does not exist.
		if err == sql.ErrNoRows {
			return 0, errors.NewNotFoundError(errors.ErrUserNotFound)
		}
		// Wrap other errors as InternalError for upstream handling.
		return 0, errors.NewInternalError(errors.ErrDatabaseQuery).WithError(err)
	}
	// Return the found user ID.
	return id, nil
}
