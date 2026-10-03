// Package http implements HTTP handlers and the routing configuration for the ecommerce-platform application.
// It provides the RouterConfig type and NewRouter factory function, which wire up handlers, middlewares, rate limiters, and static file serving to produce a fully configured *mux.Router* ready to handle API requests.
package router

import (
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/primary/http/middleware"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/security/rate_limiter"
	"github.com/gorilla/mux"
)

// RouterConfiguration defines the interface for configuring routes in the application.
type RouterConfiguration interface {
	SetupRoutes(router *mux.Router)
}

type RouterDependencies struct {
	UserServiceLogin    input.UserServiceLogin
	UserServiceRegister input.UserServiceRegister
	CommentGetService   input.CommentGetService
	CommentAddService   input.CommentAddService
	RateHandler         ratelimiter.RateLimiterHandler
	StaticFileService   output.StaticFilePort
	ProductsGetService  input.ProductsGetService
	CSRFMiddleware      *middleware.CSRFMiddleware
	CSRFService         output.CSRFService
	IsProduction        bool
	BlacklistRepo       output.TokenBlacklistPort
	TokenService        output.TokenService
}

// NewRouter constructs and returns a *mux.Router configured with all application routes, handlers, and global middleware.
// It performs the following steps:
//  1. Instantiate handler objects for login, registration, comments, etc.
//  2. Set up the MainPageHandler with the static directory path.
//  3. Create and configure a MiddlewareManager, adding global middleware
//     for logging, timing, and CORS.
//  4. Build a RouterConfig with dependencies and call SetupRoutes.

// Parameters:
//   - userServiceLogin: service for authenticating users on login.
//   - userServiceRegister: service for registering new users.
//   - commentGetService: service for fetching existing comments.
//   - commentAddService: service for adding new comments.
//   - rateHandler: rate limiting handler middleware for DoS protection.
//   - staticFileService: adapter for serving static files from disk.
//   - productsGetService: service for retrieving product catalog data.
//   - csrfMiddleware: middleware for Cross-Site Request Forgery protection.
//   - csrfService: service for CSRF token generation and validation.
//   - isProduction: flag indicating whether the app runs in production mode.
//   - blacklistRepo: repository for persisting and checking revoked JWT tokens.

// Returns:
//   - *mux.Router: fully configured router ready to be passed to http.ListenAndServe.
func NewRouter(deps RouterDependencies) *mux.Router {
	// 1. Initialize a new router
	router := mux.NewRouter()

	// 2. Build all HTTP handlers using the reusable helper
	handlers := buildHandlers(
		deps.UserServiceLogin,
		deps.UserServiceRegister,
		deps.CommentGetService,
		deps.CommentAddService,
		deps.StaticFileService,
		deps.ProductsGetService,
		deps.CSRFService,
		deps.TokenService,
		deps.BlacklistRepo,
		deps.IsProduction,
	)

	// 3. Configure global middleware using the reusable helper
	middlewareManager := configureMiddleware(router)

	// 4. Build RouterConfig with dependencies
	c := &RouterConfig{
		Handlers:          handlers,
		IPExtractor:        ratelimiter.NewDefaultIPExtractor("10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "127.0.0.1/32"),
		RateLimiter:        deps.RateHandler,
		MiddlewareManager:  middlewareManager,
		CSRFMiddleware:     deps.CSRFMiddleware,
		TokenService:       deps.TokenService,
		BlacklistRepo:      deps.BlacklistRepo,
	}

	// 5. Register routes on router
	c.SetupRoutes(router)

	return router
}
