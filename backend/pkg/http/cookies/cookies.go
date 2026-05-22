// Package http provides cookie configuration utilities with secure defaults and flexible options.
// It implements a functional options pattern for creating and managing HTTP cookies securely, with special handling for authentication cookies and production environment considerations.
package cookies

import (
	"net/http"
	"time"
)

// cookiePrefix is prepended to all cookie names to avoid collisions when
// multiple applications share the same parent domain. Configured once at
// startup via SetCookiePrefix.
var cookiePrefix string

// SetCookiePrefix configures the global prefix for all cookie names.
// Should be called once at application startup. An empty prefix is valid.
func SetCookiePrefix(prefix string) {
	cookiePrefix = prefix
}

// prefixedName returns the cookie name with the global prefix prepended.
func prefixedName(name string) string {
	return cookiePrefix + name
}

// CookieName returns the full cookie name including the configured prefix.
// Useful for middlewares and handlers that need to read cookies by name.
func CookieName(name string) string {
	return prefixedName(name)
}

// CookieConfig defines the complete configuration for HTTP cookie properties.
// Used as a template for creating standardized http.Cookie instances.
type CookieConfig struct {
	Name     string        // Cookie name (required)
	Value    string        // Cookie value (empty for session cookies)
	MaxAge   time.Duration // Duration until cookie expiration
	HttpOnly bool          // Restrict cookie access to HTTP only
	Path     string        // URL path scope for cookie
	Secure   bool          // Require HTTPS transport
	SameSite http.SameSite // SameSite policy enforcement
}

// CookieOption defines functional options for modifying CookieConfig instances.
type CookieOption func(*CookieConfig)

// WithValue sets the cookie's value string.
// Typically used for storing tokens or session identifiers.
func WithValue(value string) CookieOption {
	return func(c *CookieConfig) {
		c.Value = value
	}
}

// WithMaxAge sets the cookie's expiration duration.
// Negative values will create session cookies that expire when the browser closes.
func WithMaxAge(duration time.Duration) CookieOption {
	return func(c *CookieConfig) {
		c.MaxAge = duration
	}
}

// WithHttpOnly controls JavaScript access to the cookie.
// Recommended true for security-sensitive cookies.
func WithHttpOnly(httpOnly bool) CookieOption {
	return func(c *CookieConfig) {
		c.HttpOnly = httpOnly
	}
}

// WithPath defines the URL path scope for cookie transmission.
// Defaults to "/" for whole domain accessibility.
func WithPath(path string) CookieOption {
	return func(c *CookieConfig) {
		c.Path = path
	}
}

// WithSecure enforces HTTPS-only cookie transmission.
// Should always be true in production environments.
func WithSecure(secure bool) CookieOption {
	return func(c *CookieConfig) {
		c.Secure = secure
	}
}

// WithSameSite sets the SameSite policy for cross-site requests.
// Defaults to Lax mode for balanced security and functionality.
func WithSameSite(sameSite http.SameSite) CookieOption {
	return func(c *CookieConfig) {
		c.SameSite = sameSite
	}
}

// NewCookieConfig creates a new CookieConfig with secure defaults:
// - Path: "/"
// - MaxAge: 24 hours
// - HttpOnly: true
// - Secure: false
// - SameSite: Lax
// Options are applied in sequence to override defaults.
func NewCookieConfig(name string, options ...CookieOption) CookieConfig {
	config := CookieConfig{
		Name:     name,
		Path:     "/",
		MaxAge:   24 * time.Hour,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	for _, option := range options {
		option(&config)
	}

	return config
}

// NewAuthCookieConfig creates pre-configured authentication cookie settings:
// - Name: "token"
// - Value: Provided JWT/access token
// - MaxAge: 12 hours
// - HttpOnly: true
// - SameSite: Lax
// - Secure: Enabled in production environments
// Additional options can override these defaults.
func NewAuthCookieConfig(token string, isProduction bool, options ...CookieOption) CookieConfig {
	defaultOptions := []CookieOption{
		WithValue(token),
		WithMaxAge(12 * time.Hour),
		WithPath("/"),
		WithHttpOnly(true),
		WithSameSite(http.SameSiteLaxMode),
	}

	if isProduction {
		defaultOptions = append(defaultOptions, WithSecure(true))
	}

	allOptions := append(defaultOptions, options...)

	return NewCookieConfig(prefixedName("token"), allOptions...)
}

// SetCookie writes a cookie to the HTTP response using configuration.
// Handles expiration timing conversion from Duration to Expires/MaxAge.
func SetCookie(w http.ResponseWriter, config CookieConfig) {
	cookie := http.Cookie{
		Name:     config.Name,
		Value:    config.Value,
		HttpOnly: config.HttpOnly,
		Path:     config.Path,
		Secure:   config.Secure,
		SameSite: config.SameSite,
	}

	if config.MaxAge < 0 {
		cookie.MaxAge = -1
		cookie.Expires = time.Time{}
	} else {
		cookie.Expires = time.Now().Add(config.MaxAge)
		cookie.MaxAge = 0
	}

	http.SetCookie(w, &cookie)
}

// SetAuthCookie helper combines NewAuthCookieConfig and SetCookie for authentication workflows.
// Enforces secure settings based on production environment flag.
func SetAuthCookie(w http.ResponseWriter, token string, isProduction bool, options ...CookieOption) {
	config := NewAuthCookieConfig(token, isProduction, options...)
	SetCookie(w, config)
}

// ClearCookie invalidates a cookie by setting empty value and immediate expiration.
// Uses path "/" to ensure proper invalidation across all paths.
// The isProduction flag controls the Secure flag to match the original auth cookie.
func ClearCookie(w http.ResponseWriter, name string, isProduction bool) {
	config := CookieConfig{
		Name:     prefixedName(name),
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteLaxMode,
	}
	SetCookie(w, config)
}

// SetRefreshCookie sets a refresh token cookie with secure defaults.
// The cookie is HttpOnly, SameSite=Strict, has a 7-day expiration, and
// Secure flag based on production environment.
func SetRefreshCookie(w http.ResponseWriter, token string, isProduction bool) {
	config := NewCookieConfig(prefixedName("refresh_token"),
		WithValue(token),
		WithMaxAge(7*24*time.Hour),
		WithHttpOnly(true),
		WithSecure(isProduction),
		WithSameSite(http.SameSiteStrictMode),
	)
	SetCookie(w, config)
}

// SetCSRFCookie sets a CSRF token cookie with secure defaults.
// The cookie is HttpOnly, has SameSite=Strict, and Secure flag based on production environment.
func SetCSRFCookie(w http.ResponseWriter, token string, isProduction bool) {
	config := NewCookieConfig(prefixedName("csrf_token"),
		WithValue(token),
		WithMaxAge(24*time.Hour),
		WithHttpOnly(true),
		WithSameSite(http.SameSiteStrictMode),
		WithSecure(isProduction),
	)
	SetCookie(w, config)
}
