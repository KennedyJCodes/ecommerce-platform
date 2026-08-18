package router

import (
	"net/http"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/primary/http/middleware"
	"github.com/gorilla/mux"
)

// RegisterLoginRoutes registers unauthenticated auth-entry routes.
func (c *HandlerConfig) registerLoginRoutes(router *mux.Router, rateLimitMW middleware.Middleware) {
	loginRouter := router.PathPrefix("").Subrouter()
	loginRouter.Use(mux.MiddlewareFunc(middleware.CORSMiddleware(middleware.PrivateCORSConfig())))
	loginRouter.Use(mux.MiddlewareFunc(rateLimitMW))
	loginRouter.Use(mux.MiddlewareFunc(middleware.NoCacheMiddleware()))

	loginRouter.Handle("/login", http.HandlerFunc(c.Login.Handle)).Methods("POST", "OPTIONS")
	loginRouter.Handle("/register", http.HandlerFunc(c.Register.Handle)).Methods("POST", "OPTIONS")
}
