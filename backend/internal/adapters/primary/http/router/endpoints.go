package router

import (

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/primary/http/middleware"
	"github.com/gorilla/mux"
)

// SetupRoutes registers all application endpoints on the given router and applies route-specific middleware.
func (c *RouterConfig) SetupRoutes(router *mux.Router) {
	c.Handlers.StaticFile.RegisterRoutes(router)

	rateLimitMW := middleware.RateLimitMiddleware(c.IPExtractor, c.RateLimiter)
	authOpts := middleware.DefaultAuthOptions()
	authOpts.TokenService = c.TokenService
	authOpts.BlacklistRepo = c.BlacklistRepo
	authMW := middleware.AuthMiddleware(authOpts)
	csrfMW := c.CSRFMiddleware.ProtectCR

	c.Handlers.registerPublicRoutes(router, rateLimitMW)
	c.Handlers.registerPrivateRoutes(router, rateLimitMW, authMW, csrfMW)
	c.Handlers.registerLoginRoutes(router, rateLimitMW)
}
