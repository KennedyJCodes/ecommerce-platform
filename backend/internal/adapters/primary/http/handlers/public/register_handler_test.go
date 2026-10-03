package public_handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	dto "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/dto/auth"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/http/cookies"
)

type registerHandlerServiceMock struct {
	pendingUserID string
	err           error
}

func (m *registerHandlerServiceMock) Register(_ context.Context, _ dto.RegisterAccount) (string, error) {
	return m.pendingUserID, m.err
}

func TestRegisterHandlerHandleSuccess(t *testing.T) {
	cookies.SetCookiePrefix("")
	handler := NewRegisterHandler(&registerHandlerServiceMock{pendingUserID: "pending-1"}, false)
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"username":"alice","email":"alice@example.com","password":"secret"}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var response map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if response["next"] != "/verify" {
		t.Fatalf("next = %q, want /verify", response["next"])
	}
}

func TestRegisterHandlerHandleRejectsInvalidMethod(t *testing.T) {
	handler := NewRegisterHandler(&registerHandlerServiceMock{}, false)
	req := httptest.NewRequest(http.MethodGet, "/register", nil)
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}
