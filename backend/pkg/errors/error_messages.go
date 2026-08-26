// Package errors defines common error messages used throughout the application.
// It provides consistent error messaging for various domains including authentication,
// database operations, validation, and API rate limiting.
package errors

// Common error messages
const (
	// Authentication errors
	ErrHashingPassword    = "Error hashing password"
	ErrInvalidCredentials = "Invalid credentials"
	ErrUserNotFound       = "User not found"
	ErrUserAlreadyExists  = "The user already exists"
	ErrInvalidUsername    = "Invalid username"
	ErrInvalidPassword    = "Invalid password"
	ErrTokenGeneration    = "Error generating token"
	ErrTokenValidation    = "Invalid or expired token"
	ErrInvalidCSRFToken   = "Invalid CSRF token"
	ErrCSRFTokenExpired   = "CSRF token expired"
	ErrCSRFTokenNotFound  = "CSRF token not found"
	ErrInvalidEmail       = "Invalid email address"
	ErrEmailAlreadyExists  = "The email address is already registered"

	// Database errors
	ErrDatabaseTransaction = "Error starting database transaction"
	ErrDatabaseRollback    = "Error rolling back database transaction"
	ErrDatabaseCommit      = "Error committing database transaction"
	ErrDatabaseConnection  = "Database connection error"
	ErrDatabaseQuery       = "Error executing query"
	ErrDatabaseInsert      = "Error inserting into the database"
	ErrDatabaseUpdate      = "Error updating the database"
	ErrDatabaseDelete      = "Error deleting from database"

	// Validation errors
	ErrEmptyField        = "The field cannot be empty"
	ErrInvalidFormat     = "Invalid format"
	ErrInvalidLength     = "Invalid length"
	ErrInvalidCharacters = "Characters not allowed"

	// Comment operations errors
	ErrCommentNotFound = "Comment not found"
	ErrCommentCreation = "Error creating comment"
	ErrCommentUpdate   = "Error updating comment"
	ErrCommentDelete   = "Error deleting comment"

	// Rate limiting errors
	ErrTooManyRequests   = "Too many requests"
	ErrRateLimitExceeded = "Rate limit exceeded"

	// General API errors
	ErrInternalServer       = "Internal Server Error"
	ErrMethodNotAllowed     = "Disallowed method"
	ErrInvalidRequest       = "Invalid request"
	ErrRequestBodyTooLarge  = "Request body too large"
	ErrUnsupportedMediaType = "Content-Type must be application/json"
	ErrUnauthorized         = "Unauthorized"
	ErrForbidden            = "Prohibited access"
)
