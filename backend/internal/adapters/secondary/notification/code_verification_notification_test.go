package notification

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestSendCodeVerification(t *testing.T) {
	notification := NewCodeVerificationNotification("test-api-key", "sender@example.com")
	var receivedRequest *http.Request
	notification.httpClient.Transport = roundTripFunc(func(request *http.Request) (*http.Response, error) {
		receivedRequest = request
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"id":"email-id"}`)),
			Header:     make(http.Header),
			Request:    request,
		}, nil
	})

	if err := notification.SendCodeVerification("recipient@example.com", "012345"); err != nil {
		t.Fatalf("SendCodeVerification() error = %v", err)
	}
	if receivedRequest == nil {
		t.Fatal("expected an HTTP request")
	}
	if receivedRequest.Method != http.MethodPost {
		t.Errorf("request method = %s, want %s", receivedRequest.Method, http.MethodPost)
	}
	if receivedRequest.Header.Get("Authorization") != "Bearer test-api-key" {
		t.Errorf("unexpected authorization header: %q", receivedRequest.Header.Get("Authorization"))
	}

	body, err := io.ReadAll(receivedRequest.Body)
	if err != nil {
		t.Fatalf("reading request body: %v", err)
	}
	var payload CodeVerificationPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("decoding request body: %v", err)
	}
	if payload.From != "sender@example.com" || len(payload.To) != 1 || payload.To[0] != "recipient@example.com" {
		t.Errorf("unexpected email addresses: %+v", payload)
	}
	if !strings.Contains(payload.HTML, "012345") {
		t.Errorf("payload HTML does not contain the verification code: %s", payload.HTML)
	}
}

func TestSendCodeVerificationReturnsErrorForFailedResponse(t *testing.T) {
	notification := NewCodeVerificationNotification("test-api-key", "sender@example.com")
	notification.httpClient.Transport = roundTripFunc(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(strings.NewReader(`{"error":"invalid recipient"}`)),
			Header:     make(http.Header),
			Request:    request,
		}, nil
	})

	if err := notification.SendCodeVerification("recipient@example.com", "012345"); err == nil {
		t.Fatal("SendCodeVerification() error = nil, want error")
	}
}
