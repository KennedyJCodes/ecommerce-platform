// Package output defines output port interfaces for the application.
package output

// CSRFCookieSetter defines an interface for setting CSRF token cookies.
// This allows services to set CSRF cookies without directly depending on HTTP types.
type CSRFCookieSetter interface {
	// SetCSRFCookie sets a CSRF token cookie with the given token value.
	// The implementation should handle secure cookie settings based on the environment.
	SetCSRFCookie(token string)
}
