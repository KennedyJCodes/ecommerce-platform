// Package middleware provides HTTP middleware components for the application.
// It includes CORS handling, authentication, logging, and other cross-cutting concerns.
package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

// CORSConfig defines the configuration options for Cross-Origin Resource Sharing (CORS).
//
// Implementation details:
//   - All fields are exported for direct modification
//   - AllowedOrigins supports wildcard "*" for development
//   - AllowedMethods should include OPTIONS for preflight support
//   - AllowedHeaders should include any custom headers used
//   - AllowCredentials affects cookie and authentication handling
//   - MaxAge is used only for preflight responses
//
// Validation rules:
//   - When AllowCredentials is true, AllowedOrigins cannot be ["*"]
//   - AllowedMethods should not be empty
//   - MaxAge should be positive for preflight caching
type CORSConfig struct {
	// AllowedOrigins is a list of allowed origins (use "*" for all)
	AllowedOrigins []string
	// AllowedMethods is a list of allowed HTTP methods
	AllowedMethods []string
	// AllowedHeaders is a list of allowed custom headers
	AllowedHeaders []string
	// AllowCredentials indicates if credentials are allowed
	AllowCredentials bool
	// ExposedHeaders is a list of headers to expose
	ExposedHeaders []string
	// MaxAge is the maximum age of the preflight cache
	MaxAge int
}

func PublicCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
		MaxAge:           3600,
	}
}

func PrivateCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"X-CSRF-Token"},
		MaxAge:           86400,
	}
}

// CORSMiddleware creates a new middleware handler that implements CORS.
// It checks if the request origin is allowed and adds appropriate CORS headers
// to the response.
//
// Implementation details:
//  1. Extracts the Origin header from the request
//  2. Checks if the origin is in the allowed origins list
//  3. If origin is allowed:
//     - Sets Access-Control-Allow-Origin header
//     - Sets Access-Control-Allow-Credentials if enabled
//     - For OPTIONS requests (preflight):
//     * Sets Access-Control-Allow-Methods
//     * Sets Access-Control-Allow-Headers
//     * Sets Access-Control-Max-Age
//     * Returns 200 OK immediately
//     - For regular requests:
//     * Sets Access-Control-Expose-Headers if configured
//  4. Passes control to the next handler
//
// Security considerations:
//   - When AllowCredentials is true, "*" cannot be used for AllowedOrigins
//   - Preflight requests are handled separately to optimize performance
//   - Headers are set only for allowed origins
func CORSMiddleware(config *CORSConfig) Middleware {
	if config.AllowCredentials {
		for _, origin := range config.AllowedOrigins {
			if origin == "*" {
				panic("CORS: no se puede usar AllowCredentials con comodín de origen '*'")
			}
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if the origin is allowed
			var allowedOrigin string
			allowAll := false
			for _, o := range config.AllowedOrigins {
				if o == "*" {
					allowAll = true
					break
				}
				if o == origin {
					allowedOrigin = origin
					break
				}
			}

			if allowAll && !config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else if allowedOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
				if config.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			} else {
				next.ServeHTTP(w, r)
				return
			}

			// For OPTIONS requests (preflight)
			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))

				if config.MaxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
				}

				// Respond immediately to preflight requests
				w.WriteHeader(http.StatusNoContent)
				return
			}

			// Expose headers if configured
			if len(config.ExposedHeaders) > 0 {
				w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ", "))
			}

			// Continue to the next middleware function
			next.ServeHTTP(w, r)
		})
	}
}
