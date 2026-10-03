// Package middleware provides HTTP middleware components for cross-cutting concerns.
// It includes authentication, logging, and context propagation utilities.
// Middleware functions wrap HTTP handlers to perform tasks before or after
// the main request processing.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/auth"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/http/cookies"
)

// contextKey is a private type used to define keys for context values.
// Using a distinct type prevents collisions with other context keys.
type contextKey string

// userIDContextKey is the key under which the authenticated user's ID is stored
// in the request context. It is unexported to prevent direct use outside this package.
const userIDContextKey contextKey = "userID"

// claimsContextKey is the key under which the validated JWT claims are stored
// in the request context.
const claimsContextKey contextKey = "claims"

// GetUserIDContextKey returns the context key used to retrieve the authenticated
// user's ID from the request context.
//
// Example:
//
//	userID := r.Context().Value(middleware.GetUserIDContextKey())
func GetUserIDContextKey() contextKey {
	return userIDContextKey
}

// GetClaimsContextKey returns the context key used to retrieve the validated
// JWT claims from the request context.
func GetClaimsContextKey() contextKey {
	return claimsContextKey
}

// AuthOptions contains the configuration required by the authentication middleware.
//
// ExcludedPaths contains paths that do not require authentication through this
// middleware. TokenService validates the JWT, while BlacklistRepo optionally
// checks whether the token has been revoked.
type AuthOptions struct {
	// ExcludedPaths contains URL paths that do not require authentication.
	// Exact paths are supported, as well as path prefixes ending with "/".
	ExcludedPaths []string
	// TokenService validates the JWT and returns its claims.
	TokenService output.TokenService
	// BlacklistRepo checks whether a validated JWT has been revoked.
	BlacklistRepo output.TokenBlacklistPort
}

// DefaultAuthOptions creates the default authentication middleware configuration.
//
// The default configuration excludes public routes and static asset directories
// from access-token authentication.
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

// AuthMiddleware creates middleware that authenticates requests using an access
// JWT stored in a cookie.
//
// For protected routes, the middleware:
//  1. Retrieves the token from the authentication cookie.
//  2. Validates the JWT and its claims.
//  3. Ensures the token is an access token.
//  4. Checks whether the token has been revoked.
//  5. Stores the authenticated user's ID and validated claims in the request context.
//
// Requests matching ExcludedPaths bypass authentication.
//
// If the cookie is missing, the token is invalid or expired, the token is not
// an access token, or the token has been revoked, the middleware responds with
// HTTP 401 Unauthorized. If the blacklist repository returns an internal error,
// the middleware responds with HTTP 500 Internal Server Error.
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
			cookie, err := r.Cookie(cookies.CookieName("token"))
			if err != nil || cookie.Value == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString := cookie.Value
			claims, err := options.TokenService.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			if claims.Type != string(models_auth.TokenTypeAccess) {
    			http.Error(w, "Invalid token type", http.StatusUnauthorized)
   				return
			}

			// Check whether the validated token has been revoked.
			if options.BlacklistRepo != nil {
				blacklisted, err := options.BlacklistRepo.IsBlacklisted(claims.ID)
				if err != nil {
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}

				if blacklisted {
					http.Error(w, "Token revoked", http.StatusUnauthorized)
					return
				}
			}

			ctx := r.Context()
			contextWithUser := context.WithValue(ctx, userIDContextKey, claims.UserID)
			contextWithClaims := context.WithValue(contextWithUser, claimsContextKey, claims)

			next.ServeHTTP(w, r.WithContext(contextWithClaims))
		})
	}
}