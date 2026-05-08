// Package middleware provides HTTP middleware components for the application.
// It includes protection against Cross-Site Request Forgery (CSRF).
package middleware

import (
	"fmt"
	"net/http"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/services/service_csrf"
)

// CSRFMiddleware coordinates CSRF protection by validating tokens against a dedicated CSRF service. It uses a Double Submit Cookie pattern coupled with server-side validation.
type CSRFMiddleware struct {
	// service_csrf provides the business logic for token validation and lifecycle.
	service_csrf *service_csrf.CSRFUseCase
}

// NewCSRFMiddleware creates a new instance of CSRFMiddleware with the provided CSRF service.
// Returns a pointer to the initialized CSRFMiddleware.
func NewCSRFMiddleware(service_csrf *service_csrf.CSRFUseCase) *CSRFMiddleware {
	return &CSRFMiddleware{service_csrf: service_csrf}
}

// ProtectCR (Protect Collaborative Resources) is a middleware that enforces CSRF  validation for state-changing HTTP methods.

// Implementation details:
//  1. Filters by HTTP Method: Only POST, PUT, DELETE, and PATCH methods are intercepted.
//     Safe methods (GET, HEAD, OPTIONS, TRACE) bypass validation automatically.
//  2. Authentication Check: Retrieves the userID from the request context.
//     Requires the authentication middleware to have run previously.
//  3. Token Extraction:
//     - Mandatory: "csrf_token" cookie.
//     - Complementary: Token from "X-CSRF-Token" header or "csrf_token" form field.
//  4. Double Submit Validation: Ensures the cookie value matches the provided header/form value.
//  5. Service Validation: Calls the CSRF service to verify token integrity and user ownership.

// Security considerations:
//   - This middleware depends on GetUserIDContextKey() to identify the user.
//   - If the user is not authenticated (missing context value), it returns 401 Unauthorized.
//   - Mismatched or missing tokens result in a 403 Forbidden response.
func (m *CSRFMiddleware) ProtectCR(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only check CSRF for non-safe methods (State-changing operations)
		if r.Method != "POST" && r.Method != "PUT" &&
			r.Method != "DELETE" && r.Method != "PATCH" {
			next.ServeHTTP(w, r)
			return
		}

		// Retrieve authenticated user ID from context
		userIDValue := r.Context().Value(GetUserIDContextKey())
		userIDInt, ok := userIDValue.(int)
		if !ok {
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}
		userID := fmt.Sprintf("%d", userIDInt)

		// 1. Get the token from the cookie
		cookie, err := r.Cookie("csrf_token")
		if err != nil {
			http.Error(w, "CSRF token not found", http.StatusForbidden)
			return
		}

		// 2. Get the token from the Header or Form (Double Submit)
		tokenFromHeader := r.Header.Get("X-CSRF-Token")
		if tokenFromHeader == "" {
			tokenFromHeader = r.FormValue("csrf_token")
		}

		// 3. Basic match check (Double Submit Cookie Pattern)
		if cookie.Value != tokenFromHeader {
			http.Error(w, "CSRF token does not match", http.StatusForbidden)
			return
		}

		// 4. Cryptographic/Logical validation via service
		if err := m.service_csrf.ValidateToken(cookie.Value, userID); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
