package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	appErrors "github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
)

// These tests document the shared JSON body contract used by auth and comment handlers.
type jsonBodyTestPayload struct {
	Name string `json:"name"`
}

func TestDecodeJSONBodyAcceptsJSONWithCharset(t *testing.T) {
	req := newJSONBodyRequest(`{"name":"David"}`)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	var payload jsonBodyTestPayload
	err := DecodeJSONBody(httptest.NewRecorder(), req, &payload, 1024)

	if err != nil {
		t.Fatalf("expected valid JSON body, got error: %v", err)
	}
	if payload.Name != "David" {
		t.Fatalf("expected decoded name David, got %q", payload.Name)
	}
}

func TestDecodeJSONBodyRejectsUnknownFields(t *testing.T) {
	req := newJSONBodyRequest(`{"name":"David","role":"admin"}`)

	var payload jsonBodyTestPayload
	err := DecodeJSONBody(httptest.NewRecorder(), req, &payload, 1024)

	assertAppErrorCode(t, err, http.StatusBadRequest)
}

func TestDecodeJSONBodyRejectsOversizedBody(t *testing.T) {
	req := newJSONBodyRequest(`{"name":"David"}`)

	var payload jsonBodyTestPayload
	err := DecodeJSONBody(httptest.NewRecorder(), req, &payload, 8)

	assertAppErrorCode(t, err, http.StatusRequestEntityTooLarge)
}

func TestDecodeJSONBodyRejectsMultipleJSONValues(t *testing.T) {
	req := newJSONBodyRequest(`{"name":"David"}{"name":"Alejandro"}`)

	var payload jsonBodyTestPayload
	err := DecodeJSONBody(httptest.NewRecorder(), req, &payload, 1024)

	assertAppErrorCode(t, err, http.StatusBadRequest)
}

func TestDecodeJSONBodyRejectsUnsupportedMediaType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"David"}`))
	req.Header.Set("Content-Type", "text/plain")

	var payload jsonBodyTestPayload
	err := DecodeJSONBody(httptest.NewRecorder(), req, &payload, 1024)

	assertAppErrorCode(t, err, http.StatusUnsupportedMediaType)
}

func newJSONBodyRequest(body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func assertAppErrorCode(t *testing.T, err error, code int) {
	t.Helper()

	appErr, ok := err.(*appErrors.AppError)
	if !ok {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != code {
		t.Fatalf("expected status code %d, got %d", code, appErr.Code)
	}
}
