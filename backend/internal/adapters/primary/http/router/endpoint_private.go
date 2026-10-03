package router

import (
	"net/http"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/primary/http/middleware"
	"github.com/gorilla/mux"
)

// RegisterPrivateRoutes registers routes that require authentication and CSRF protection.
func (c *HandlerConfig) registerPrivateRoutes(router *mux.Router, rateLimitMW, authMW, csrfMW middleware.Middleware) {
	private := router.PathPrefix("").Subrouter()
	private.Use(mux.MiddlewareFunc(middleware.CORSMiddleware(middleware.PrivateCORSConfig())))
	private.Use(mux.MiddlewareFunc(rateLimitMW))
	private.Use(mux.MiddlewareFunc(authMW))
	private.Use(mux.MiddlewareFunc(csrfMW))
	private.Use(mux.MiddlewareFunc(middleware.NoCacheMiddleware()))

	private.Handle("/comments/newComments", http.HandlerFunc(c.ReviewsAdd.Handle)).Methods("POST", "OPTIONS")
	private.Handle("/logout", http.HandlerFunc(c.Logout.Handle)).Methods("POST", "OPTIONS")
}
