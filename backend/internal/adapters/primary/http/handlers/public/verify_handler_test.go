package public_handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	models_auth "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/auth"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/http/cookies"
)

type verifyHandlerServiceMock struct {
	pendingUserID string
	code          string
	tokens        *models_auth.TokenPair
	csrfToken     string
	err           error
}

func (m *verifyHandlerServiceMock) Verify(_ context.Context, pendingUserID, code string) (*models_auth.TokenPair, string, error) {
	m.pendingUserID = pendingUserID
	m.code = code
	return m.tokens, m.csrfToken, m.err
}

func TestVerifyHandlerHandleSuccess(t *testing.T) {
	cookies.SetCookiePrefix("")
	service := &verifyHandlerServiceMock{
		tokens:    &models_auth.TokenPair{AccessToken: "access", RefreshToken: "refresh"},
		csrfToken: "csrf",
	}
	handler := NewVerifyHandler(service, false)

	req := httptest.NewRequest(http.MethodPost, "/verify", strings.NewReader(`{"code":"123456"}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "pending_user_id", Value: "pending-1"})
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if service.pendingUserID != "pending-1" || service.code != "123456" {
		t.Fatalf("service received (%q, %q), want (%q, %q)", service.pendingUserID, service.code, "pending-1", "123456")
	}
	if len(recorder.Result().Cookies()) != 4 {
		t.Fatalf("cookies = %d, want access, refresh, csrf and pending cookie", len(recorder.Result().Cookies()))
	}
}

func TestVerifyHandlerHandleRejectsMissingPendingCookie(t *testing.T) {
	handler := NewVerifyHandler(&verifyHandlerServiceMock{}, false)
	req := httptest.NewRequest(http.MethodPost, "/verify", strings.NewReader(`{"code":"123456"}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}

func TestVerifyHandlerHandlePropagatesServiceError(t *testing.T) {
	handler := NewVerifyHandler(&verifyHandlerServiceMock{err: errors.New("verification failed")}, false)
	req := httptest.NewRequest(http.MethodPost, "/verify", strings.NewReader(`{"code":"123456"}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "pending_user_id", Value: "pending-1"})
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
}
