package router

import (
	"net/http"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/primary/http/middleware"
	"github.com/gorilla/mux"
)

// RegisterPublicRoutes registers routes that do not require access-token authentication.
func (c *HandlerConfig) registerPublicRoutes(router *mux.Router, rateLimitMW middleware.Middleware) {
	public := router.PathPrefix("").Subrouter()
	public.Use(mux.MiddlewareFunc(middleware.CORSMiddleware(middleware.PublicCORSConfig())))
	public.Use(mux.MiddlewareFunc(rateLimitMW))

	public.Handle("/", http.HandlerFunc(c.MainPage.Handle)).Methods("GET", "OPTIONS")
	public.Handle("/comments", http.HandlerFunc(c.ReviewsGet.Handle)).Methods("GET", "OPTIONS")
	public.Handle("/products", http.HandlerFunc(c.Products.Handle)).Methods("GET", "OPTIONS")
	public.Handle("/product-id/{id}", http.HandlerFunc(c.Products.HandleGetByID)).Methods("GET", "OPTIONS")
	public.Handle("/products-brand/{brand}", http.HandlerFunc(c.Products.HandleGetByBrand)).Methods("GET", "OPTIONS")
	public.Handle("/refresh", http.HandlerFunc(c.Refresh.Handle)).Methods("POST", "OPTIONS")
}
