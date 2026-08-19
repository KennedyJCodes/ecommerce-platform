// Package http implements HTTP handlers for the ecommerce-platform application.
// This file contains the RefreshHandler, which processes token refresh requests
// by validating the refresh token, revoking the old one, and issuing a new
// access and refresh token pair.
package http

import (
	"net/http"
	"time"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
	httpUtil "github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/http"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/http/cookies"
)

// RefreshHandler handles HTTP requests for token renewal.
// It validates the refresh token stored in a cookie, checks it against the
// blacklist, revokes the old token (rotation), and issues a fresh token pair.
type RefreshHandler struct {
	tokenService  output.TokenService
	blacklistRepo output.TokenBlacklistPort
	isProduction  bool
}

// NewRefreshHandler creates and returns a new RefreshHandler instance.
// Parameters:
//   - tokenService: the service for JWT generation and validation.
//   - blacklistRepo: an implementation of the TokenBlacklistPort for token revocation.
//   - isProduction: whether the app runs in production (sets Secure flag on cookies).
//
// Returns:
//   - *RefreshHandler: ready-to-use handler for the /refresh endpoint.
func NewRefreshHandler(tokenService output.TokenService, blacklistRepo output.TokenBlacklistPort, isProduction bool) *RefreshHandler {
	return &RefreshHandler{tokenService: tokenService, blacklistRepo: blacklistRepo, isProduction: isProduction}
}

// Handle processes HTTP POST requests for token refresh.
// The flow follows a strict security sequence:
//  1. Read the refresh_token cookie from the request.
//  2. Validate the refresh token and its claims.
//  3. Ensure the token is a refresh token.
//  4. Check if the refresh token has been revoked.
//  5. Generate a new access token and refresh token.
//  6. Revoke the old refresh token as part of token rotation.
//  7. Set both new tokens as HttpOnly cookies on the response.
//  8. Respond with a JSON success message.
func (h *RefreshHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpUtil.HandleError(w, errors.NewBadRequestError(errors.ErrMethodNotAllowed))
		return
	}

	cookie, err := r.Cookie(cookies.CookieName("refresh_token"))
	if err != nil || cookie.Value == "" {
		httpUtil.HandleError(w, errors.NewAuthError("Refresh token missing"))
		return
	}

	claims, err := h.tokenService.ValidateToken(cookie.Value)
	if err != nil {
		httpUtil.HandleError(w, errors.NewAuthError("Invalid or expired refresh token"))
		return
	}

	if claims.Type != "refresh" {
    	httpUtil.HandleError(w, errors.NewAuthError("Invalid token type"))
    	return
	}	

	blacklisted, err := h.blacklistRepo.IsBlacklisted(claims.ID)
	if err != nil || blacklisted {
		httpUtil.HandleError(w, errors.NewAuthError("Refresh token revoked"))
		return
	}

	newAccessToken, err := h.tokenService.GenerateToken(claims.UserID, claims.UserName, models.TokenTypeAccess)
	if err != nil {
		httpUtil.HandleError(w, errors.NewInternalError("Error generating access token"))
		return
	}
	
	newRefreshToken, err := h.tokenService.GenerateToken(claims.UserID, claims.UserName, models.TokenTypeRefresh)
	if err != nil {
		httpUtil.HandleError(w, errors.NewInternalError("Error generating refresh token"))
		return
	}

	// Blacklist the old refresh token as part of token rotation.
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl < 0 {
		ttl = 0
	}

	if claims.ID != "" {	
		httpUtil.HandleError(w, errors.NewAuthError("Invalid refresh token"))
		return
	}

	if err := h.blacklistRepo.Add(claims.ID, ttl); err != nil {
		httpUtil.HandleError(w, errors.NewInternalError("Error revoking refresh token"))
		return
	}

	cookies.SetAuthCookie(w, newAccessToken, h.isProduction)
	cookies.SetRefreshCookie(w, newRefreshToken, h.isProduction)

	httpUtil.SendJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Token refreshed successfully",
	})
}
