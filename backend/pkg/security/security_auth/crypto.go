// Package securityAuth provides interfaces and implementations for password hashing
// and salt generation, supporting secure authentication workflows in the sale-watches application.
package security_auth

import (
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Hasher defines methods for hashing password byte slices.
type Hasher interface {
	// Hash takes a password as a byte slice and returns its hashed representation or an error if hashing fails.
	Hash(password []byte) (string, error)
}

// BcryptHasher implements the Hasher interface using bcrypt with the DefaultCost.
type BcryptHasher struct {}

// Hash generates a bcrypt hash of the provided password bytes using DefaultCost.
// It returns the resulting hash as a string, or an InternalError if hashing fails.
func (h BcryptHasher) Hash(password []byte) (string, error) {
	var hashPassword, err = bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", errors.NewInternalError("error hashing the password")
	}
	return string(hashPassword), nil
}
