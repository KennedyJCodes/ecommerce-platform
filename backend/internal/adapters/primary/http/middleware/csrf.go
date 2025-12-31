package middleware

import (
	"fmt"
	"net/http"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/services/service_csrf"
)

type CSRFMiddleware struct {
	service_csrf *service_csrf.CSRFUseCase
}

func NewCSRFMiddleware(service_csrf *service_csrf.CSRFUseCase) *CSRFMiddleware {
	return &CSRFMiddleware{service_csrf: service_csrf}
}

func (m *CSRFMiddleware) ProtectCR(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" && r.Method != "PUT" &&
			r.Method != "DELETE" && r.Method != "PATCH" {
			next.ServeHTTP(w, r)
			return
		}

		userIDValue := r.Context().Value(GetUserIDContextKey())
		userIDInt, ok := userIDValue.(int)
		if !ok {
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}
		userID := fmt.Sprintf("%d", userIDInt)

		cookie, err := r.Cookie("csrf_token")
		if err != nil {
			http.Error(w, "CSRF token not found", http.StatusForbidden)
			return
		}

		tokenFromHeader := r.Header.Get("X-CSRF-Token")
		if tokenFromHeader == "" {
			tokenFromHeader = r.FormValue("csrf_token")
		}

		
		if cookie.Value != tokenFromHeader {
			http.Error(w, "CSRF token does not match", http.StatusForbidden)
			return
		}

		if err := m.service_csrf.ValidateToken(cookie.Value, userID); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
