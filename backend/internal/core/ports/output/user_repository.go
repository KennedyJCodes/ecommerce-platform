// Package output defines persistence contracts for comments and users.
package output

// UserRepository persists and retrieves user credentials.
type UserRepository interface {
	// UserExists checks if a username is already registered in the system.
	// Returns true if the username exists in storage, false otherwise.
	// May return an error for database connectivity issues or storage failures.
	UserExists(username string) (bool, error)

	// GetHashPassword retrieves the hashed password for a registered user.
	// Primarily used during authentication processes. Returns an error if:
		// - User doesn't exist
		// - Password record is corrupted
		// - Storage system failure occurs
	GetHashPassword(username string) (string, error)

	// SaveUser persists a new user record with secure credential storage.
	// Implementations should:
		// - Generate unique salt per user
		// - Hash password using cryptographically secure methods
		// - Prevent duplicate usernames

	// Returns error for:
		// - Duplicate username
		// - Invalid credentials
		// - Storage persistence failures
	SaveUser(username, password, email string) error

	// GetID returns the numeric ID for a given username.
    // Returns:
    //   - int: user ID.
    //   - error: non-nil if user not found or storage error.
	GetID(username string) (int, error)

	EmailExists(email string) (bool, error)
}
