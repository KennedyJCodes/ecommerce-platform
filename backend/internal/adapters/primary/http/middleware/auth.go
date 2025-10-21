// Package middleware provides HTTP middleware components for cross-cutting concerns.
// It includes authentication, logging, and context propagation utilities.
// Middleware functions wrap HTTP handlers to perform tasks before and/or after
// the main request processing.
package middleware

import (
	"context"
	"net/http"
	"strings"

	securityAuth "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/security_auth"
)

// contextKey is a private type used to define keys for context values.
// Using a distinct type prevents collisions with other context keys.
type contextKey string

// userIDContextKey is the key under which the authenticated user's ID is stored
// in the request context. It is unexported to prevent misuse.
const userIDContextKey contextKey = "userID"

// GetUserIDContextKey returns the context key used to retrieve the user ID
// from an HTTP request's context. This allows handlers to extract the
// authenticated user's ID from context:
//
//	userID := r.Context().Value(middleware.GetUserIDContextKey()).(string)
func GetUserIDContextKey() contextKey {
	return userIDContextKey
}

// AuthOptions contains configuration options for the authentication middleware.
// It allows for customization of authentication behavior, particularly which
// paths should be excluded from authentication requirements.
type AuthOptions struct {
	// ExcludedPaths is a slice of URL patterns that will not require authentication.
	// Both exact matches and directory prefixes (ending with '/') are supported.
	ExcludedPaths []string
}

// DefaultAuthOptions creates and returns a new AuthOptions instance with default values.
// The default configuration excludes common public paths like home, authentication pages,
// and static asset directories from requiring authentication.
// Returns a pointer to the newly created AuthOptions.
func DefaultAuthOptions() *AuthOptions {
	return &AuthOptions{
		ExcludedPaths: []string{
			"/",
			"/login",
			"/comments",
			"/register",
			"/css/",
			"/js/",
			"/products",
			"/assets/"},
	}
}

// AuthMiddleware returns an HTTP middleware that enforces authentication using JWT tokens stored in cookies. It wraps an existing http.Handler and performs the following logic:

// 1. If the request path matches any of the patterns in opts.ExcludedPaths, the request is allowed to proceed without authentication.
// 2. Otherwise, the middleware looks for a "token" cookie in the request.
// 3. If the cookie is missing or empty, responds with 401 Unauthorized.
// 4. Parses and validates the JWT token using the security_auth package.
// 5. If token parsing fails (invalid or expired), responds with 401 Unauthorized.
// 6. On successful validation, extracts the UserId from token claims, stores it in the request context under userIDContextKey, and calls the next handler.

// Parameters:
//   - opts: pointer to AuthOptions specifying paths to exclude from auth.

// Returns:
//   - Middleware: a function that transforms an http.Handler into a new one with authentication enforcement.
func AuthMiddleware(options *AuthOptions) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, excludedPath := range options.ExcludedPaths {
				if excludedPath == "/" {
					if r.URL.Path == "/" {
						next.ServeHTTP(w, r)
						return
					}
					continue
				}

				if r.URL.Path == excludedPath {
					next.ServeHTTP(w, r)
					return
				}
				if strings.HasSuffix(excludedPath, "/") && strings.HasPrefix(r.URL.Path, excludedPath) {
					next.ServeHTTP(w, r)
					return
				}
			}
			cookie, err := r.Cookie("token")
			if err != nil || cookie.Value == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString := cookie.Value
			claims, err := securityAuth.ParseTokenWithClaims(tokenString)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			contextWithUser := context.WithValue(ctx, userIDContextKey, claims.UserId)

			next.ServeHTTP(w, r.WithContext(contextWithUser))
		})
	}
}