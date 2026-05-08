// Package http implements HTTP handlers and adapters for the ecommerce-platform application.
package http

import (
	"net/http"
	"os"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/http/cookies"
)

// csrfCookieSetter implements output.CSRFCookieSetter interface.
// It wraps an http.ResponseWriter to set CSRF token cookies.
type csrfCookieSetter struct {
	w http.ResponseWriter
}

// NewCSRFCookieSetter creates a new CSRFCookieSetter implementation.
func NewCSRFCookieSetter(w http.ResponseWriter) output.CSRFCookieSetter {
	return &csrfCookieSetter{w: w}
}

// SetCSRFCookie sets a CSRF token cookie using the wrapped ResponseWriter.
func (c *csrfCookieSetter) SetCSRFCookie(token string) {
	isProduction := os.Getenv("ENV") == "production"
	cookies.SetCSRFCookie(c.w, token, isProduction)
}
