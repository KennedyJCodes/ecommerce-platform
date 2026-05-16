// Package http implements HTTP handlers for the ecommerce-platform application.
// This file contains the RegisterHandler, which is responsible for processing user registration requests. It decodes the registration payload, calls the user registration service to create a new account, sets the authentication cookie, and sends a JSON response indicating the result.
package http

import (
	"encoding/json"
	"net/http"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/services/service_auth"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
	httpUtil "github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/http"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/http/cookies"
)

// RegisterHandler handles HTTP requests for user registration.

// It serves as an adapter between HTTP requests and the core business logic for registering users, utilizing the UserServiceRegister interface.
type RegisterHandler struct {
	userServiceRegister input.UserServiceRegister
	csrfService         input.CSRFService
	isProduction        bool
}

// NewRegisterHandler creates a new instance of RegisterHandler.

// It receives an implementation of the UserServiceRegister interface that encapsulates the business logic for user registration.
func NewRegisterHandler(userServiceRegister input.UserServiceRegister, csrfService input.CSRFService, isProduction bool) *RegisterHandler {
	return &RegisterHandler{
		userServiceRegister: userServiceRegister,
		csrfService:         csrfService,
		isProduction:        isProduction,
	}
}

// Handle processes HTTP registration requests.

// It validates that the request method is POST and decodes the incoming JSON payload into an Account model. After invoking the registration service to create a new user account, it sets an authentication cookie (using secure settings if in production) and sends a JSON response indicating a successful registration.
// In case of errors, it responds with appropriate HTTP error messages.
func (h *RegisterHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Ensure the HTTP method is POST.
	if r.Method != http.MethodPost {
		httpUtil.HandleError(w, errors.NewBadRequestError(errors.ErrMethodNotAllowed))
		return
	}

	// Decode the JSON request body into an Account instance.
	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		httpUtil.HandleError(w, errors.NewBadRequestError(errors.ErrInvalidRequest))
		return
	}

	// Create CSRF cookie setter for this request and inject into service
	csrfCookieSetter := NewCSRFCookieSetter(w, h.isProduction)
	if registerSvc, ok := h.userServiceRegister.(*service_auth.UserRegisterService); ok {
		registerSvc.BaseAuthService.CSRFCookieSetter = csrfCookieSetter
		registerSvc.BaseAuthService.CSRFService = h.csrfService
	}

	// Attempt to register the user and generate an authentication token.
	token, err := h.userServiceRegister.Register(account)
	if err != nil {
		httpUtil.HandleError(w, err)
		return
	}

	cookies.SetAuthCookie(w, token, h.isProduction)

	// Send a JSON response indicating successful registration.
	httpUtil.SendJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Successfully registered user",
	})
}
