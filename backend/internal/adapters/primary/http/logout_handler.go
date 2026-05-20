// Package http implements HTTP handlers for the ecommerce-platform application.
// This file contains the LogoutHandler, which processes logout requests by revoking
// the current JWT token and clearing the authentication cookie.
package http

import (
	"net/http"
	"time"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/primary/http/middleware"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
	httpUtil "github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/http"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/http/cookies"
)

// LogoutHandler handles HTTP requests for user logout.
// It revokes the current JWT by adding its unique identifier (jti) to the token blacklist,
// preventing further use of the same token even before its natural expiration.
type LogoutHandler struct {
	// blacklistRepo: output port for persisting revoked JWT identifiers.
	blacklistRepo output.TokenBlacklistPort
	// isProduction: controls the Secure flag on the cleared cookie.
	isProduction bool
}

// NewLogoutHandler creates and returns a new LogoutHandler instance.
// Parameters:
//   - blacklistRepo: an implementation of the TokenBlacklistPort for token revocation.
//   - isProduction: whether the app runs in production (sets Secure flag on cleared cookie).
// Returns:
//   - *LogoutHandler: ready-to-use handler for the logout endpoint.
func NewLogoutHandler(blacklistRepo output.TokenBlacklistPort, isProduction bool) *LogoutHandler {
	return &LogoutHandler{blacklistRepo: blacklistRepo, isProduction: isProduction}
}

// Handle processes HTTP POST requests for user logout.
// It retrieves the authenticated user's JWT claims from the request context,
// adds the token's jti to the blacklist with its remaining TTL, and clears
// the authentication cookie from the client.
//
// Steps:
//  1. Retrieve the JWT claims from the request context (set by AuthMiddleware).
//  2. Calculate the remaining TTL until the token's original expiration.
//  3. Add the token's unique ID (jti) to the blacklist repository.
//  4. Clear the "token" cookie on the client.
//  5. Respond with a JSON success message.
//
// Parameters:
//   - w: http.ResponseWriter to write HTTP response headers and body.
//   - r: *http.Request containing the authenticated user's context.
func (h *LogoutHandler) Handle(w http.ResponseWriter, r *http.Request) {
	claimsValue := r.Context().Value(middleware.GetClaimsContextKey())
	claims, ok := claimsValue.(*models.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ttl := 5 * time.Minute
	if claims.ExpiresAt != nil {
		ttl = time.Until(claims.ExpiresAt.Time)
		if ttl < 0 {
			ttl = 0
		}
	}

	if claims.ID != "" {
		if err := h.blacklistRepo.Add(claims.ID, ttl); err != nil {
			httpUtil.HandleError(w, errors.NewInternalError("Error logging out"))
			return
		}
	}

	cookies.ClearCookie(w, "token", h.isProduction)
	httpUtil.SendJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}
