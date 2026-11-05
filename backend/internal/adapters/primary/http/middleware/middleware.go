// Package middleware contains HTTP middleware utilities designed following a hexagonal
// architecture. It provides middleware composition and management for the HTTP router.
package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Middleware defines a function that wraps an http.Handler.
type Middleware func(http.Handler) http.Handler

// Chain combines multiple middlewares. The first middleware is the outermost.
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// MiddlewareManager manages global and route-specific functional middleware.
type MiddlewareManager struct {
	globalMiddlewares []Middleware
}

// NewMiddlewareManager creates an empty MiddlewareManager.
func NewMiddlewareManager() *MiddlewareManager {
	return &MiddlewareManager{globalMiddlewares: []Middleware{}}
}

// AddGlobal registers middleware to be applied globally to all routes.
func (m *MiddlewareManager) AddGlobal(middleware Middleware) {
	m.globalMiddlewares = append(m.globalMiddlewares, middleware)
}

// Apply wraps a handler with global and route-specific middleware.
func (m *MiddlewareManager) Apply(handler http.Handler, middlewares ...Middleware) http.Handler {
	allMiddlewares := append(m.globalMiddlewares, middlewares...)
	return Chain(allMiddlewares...)(handler)
}

// ApplyToRouter applies all registered global middleware to a mux router.
func (m *MiddlewareManager) ApplyToRouter(router *mux.Router) {
	if len(m.globalMiddlewares) > 0 {
		router.Use(func(next http.Handler) http.Handler {
			return Chain(m.globalMiddlewares...)(next)
		})
	}
}
