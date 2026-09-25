// Package securityAuth provides interfaces and implementations for hashing
// sensitive byte slices in security workflows.
package security_auth

import (
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Hasher defines methods for hashing sensitive byte slices.
type Hasher interface {
	// Hash takes a byte slice and returns its hashed representation.
	Hash(value []byte) (string, error)
}

// BcryptHasher implements the Hasher interface using bcrypt with the DefaultCost.
type BcryptHasher struct{}

// Hash generates a bcrypt hash of the provided bytes using DefaultCost.
// It returns the resulting hash as a string, or an InternalError if hashing fails.
func (h BcryptHasher) Hash(value []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(value, bcrypt.DefaultCost)
	if err != nil {
		return "", errors.NewInternalError("error hashing value")
	}
	return string(hash), nil
}
