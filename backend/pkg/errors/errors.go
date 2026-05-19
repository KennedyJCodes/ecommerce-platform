// Package errors provides structured HTTP error handling with status codes and error wrapping.
// It implements a consistent error format for web applications and REST APIs, supporting:
// - HTTP status code association
// - Error message standardization
// - Error type checking
// - Error cause wrapping
package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error with HTTP metadata and optional cause.
// Use constructor functions (NewBadRequestError, etc.) to create properly initialized instances.
type AppError struct {
	Code    int
	Message string
	Err     error
}

// Error implements the error interface, formatting the error with code and message.
// Includes wrapped error details if present in the error chain.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("(%d) %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("(%d) %s", e.Code, e.Message)
}

// Unwrap implements the error unwrapping interface for error chain inspection.
// Enables usage with errors.Is() and errors.As() from the standard library.
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithError attaches a root cause error to the AppError instance.
// Enables error chain tracking while maintaining the original AppError context.
// Returns the modified AppError to enable method chaining.
	
// NewConflictError creates 409 Conflict error for resource state conflicts
func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

// Error constructors ---------------------------------------------------------

// NewBadRequestError creates 400 Bad Request error for invalid client requests.
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

// NewConflictError creates 409 Conflict error for resource state conflicts
func NewConflictError(message string) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: message,
	}
}

// NewInternalError creates 500 Internal Server Error for unexpected failures
func NewInternalError(message string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}

// NewAuthError creates 401 Unauthorized error for authentication failures
func NewAuthError(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

// NewTooManyRequestsError creates 429 Too Many Requests for rate limiting
func NewTooManyRequestsError(message string) *AppError {
	return &AppError{
		Code:    http.StatusTooManyRequests,
		Message: message,
	}
}

// NewNotFoundError creates 404 Not Found for missing resources
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

// NewForbiddenError creates 403 Forbidden for unauthorized access attempts
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: message,
	}
}

// NewUnsupportedMediaTypeError creates 415 Unsupported Media Type for invalid Content-Type.
func NewUnsupportedMediaTypeError(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnsupportedMediaType,
		Message: message,
	}
}

// NewValidationError creates 422 Unprocessable Entity for validation failures
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
	}
}

// Error type checkers ---------------------------------------------------------

// IsNotFound checks if error is 404 Not Found type
// Works with both AppError instances and wrapped errors
func IsNotFound(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == http.StatusNotFound
	}
	return false
}

// IsAuthError checks if error is 401 Unauthorized type
// Useful for differentiating authentication failures
func IsAuthError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == http.StatusUnauthorized
	}
	return false
}

// IsValidationError checks if error is 422 Unprocessable Entity type
// Identifies validation failures from business logic
func IsValidationError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == http.StatusUnprocessableEntity
	}
	return false
}

// IsInternalError checks if error is 500 Internal Server Error type
// Helps distinguish system errors from client errors
func IsInternalError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == http.StatusInternalServerError
	}
	return false
}
