package public_handlers

import "net/http"

// VerifyHandler is the HTTP entry point reserved for verification-code validation.
type VerifyHandler struct{}

// NewVerifyHandler creates the verification endpoint handler.
func NewVerifyHandler() *VerifyHandler {
	return &VerifyHandler{}
}

// Handle keeps the route available while verification logic is not implemented yet.
func (h *VerifyHandler) Handle(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
