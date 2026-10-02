package public_handlers

import (
	"net/http"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/dto/auth"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
	httpUtil "github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/http"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/http/cookies"
)

// VerifyHandler handles verification-code submissions for pending registrations.
type VerifyHandler struct {
	userVerificationService input.UserVerificationService
	isProduction            bool
}

// NewVerifyHandler creates the verification endpoint handler.
func NewVerifyHandler(svc input.UserVerificationService, isProduction bool) *VerifyHandler {
	if svc == nil {
		panic("NewVerifyHandler: userVerificationService is nil")
	}
	return &VerifyHandler{userVerificationService: svc, isProduction: isProduction}
}

// Handle verifies the submitted code and creates the authenticated session.
func (h *VerifyHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpUtil.HandleError(w, errors.NewBadRequestError(errors.ErrMethodNotAllowed))
		return
	}

	pendingUserCookie, err := r.Cookie(cookies.CookieName("pending_user_id"))
	if err != nil || pendingUserCookie.Value == "" {
		httpUtil.HandleError(w, errors.NewAuthError(errors.ErrInvalidCredentials))
		return
	}

	var request dto.VerifyCodeRequest
	if err := httpUtil.DecodeJSONBody(w, r, &request, httpUtil.MaxAuthBodySize); err != nil {
		httpUtil.HandleError(w, err)
		return
	}

	tokens, csrfToken, err := h.userVerificationService.Verify(r.Context(), pendingUserCookie.Value, request.Code)
	if err != nil {
		httpUtil.HandleError(w, err)
		return
	}

	cookies.SetAuthCookie(w, tokens.AccessToken, h.isProduction)
	cookies.SetRefreshCookie(w, tokens.RefreshToken, h.isProduction)
	cookies.SetCSRFCookie(w, csrfToken, h.isProduction)
	cookies.ClearPendingUserCookie(w, h.isProduction)

	httpUtil.SendJSONResponse(w, http.StatusOK, map[string]string{
		"message":    "Account verified successfully",
		"csrf_token": csrfToken,
	})
}
