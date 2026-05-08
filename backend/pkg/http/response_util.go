// Package http provides response handling utilities for HTTP APIs.
// It includes helpers for standardized JSON responses and error handling integration with the application's error types.
package http

import (
	"encoding/json"
	"net/http"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
)

// SendJSONResponse sends a JSON-encoded response with proper headers and status code.
// Sets the Content-Type header to "application/json" and encodes the provided data.
// Note: If JSON encoding fails, this will send a partial response header. Consider using error handling middleware to catch encoding errors.
func SendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// HandleError processes application errors and sends appropriate HTTP responses.
// Recognizes errors of type *errors.AppError to send structured responses with proper status codes and messages. Falls back to 500 Internal Server Error for unexpected error types.
// Usage note: Should typically be used as the final error handler in request chains.
func HandleError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		http.Error(w, appErr.Message, appErr.Code)
	} else {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
