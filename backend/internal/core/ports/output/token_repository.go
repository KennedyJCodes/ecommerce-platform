// Package output defines persistence contracts for tokens.
package output

// TokenRepository manages token storage and retrieval operations.
type TokenRepository interface {
	// GetValidToken retrieves a valid token from storage.
	// Returns:
	//   - string: the token if it exists and is valid
	//   - bool: true if token exists and is valid, false otherwise
	//   - error: non-nil if storage operation fails
	GetValidToken() (string, bool, error)
}