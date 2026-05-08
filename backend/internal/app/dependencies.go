// Package app provides the main application entry point and dependency orchestration.
// This file focuses on assembling the core dependencies and wiring them into the primary HTTP adapter (the router).
package app

import (
	"net/http"

	primaryHttp "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/primary/http"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/primary/http/middleware"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/bootstrap"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	ratelimiter "github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/security/rate_limiter"
)

// Dependencies holds all the initialized services, ports, and handlers required  to run the application.
// By grouping these in a single struct, the application ensures that all required components are ready before starting the HTTP server.
type Dependencies struct {
	UserServiceLogin    input.UserServiceLogin
	UserServiceRegister input.UserServiceRegister
	CommentGetService   input.CommentGetService
	CommentAddService   input.CommentAddService
	RateHandler         ratelimiter.RateLimiterHandler
	StaticFileAdapter   output.StaticFilePort
	ProductsGetService  input.ProductsGetService
	CSRFMiddleware      *middleware.CSRFMiddleware
	CSRFService         input.CSRFService
}

// BuildDependencies orchestrates the initialization of all internal services and repositories.
// It uses the bootstrap package to set up:
//  1. Repositories (User, Comments, Products).
//  2. Core Services (Auth, CSRF, Comments).
//  3. Infrastructure Adapters (Rate Limiter, Static Files).
//
// Returns:
//   - *Dependencies: a pointer to the fully populated Dependencies struct.
func (a *Application) BuildDependencies() *Dependencies {
	// Initialize core repositories and services using shared database connections.
	userRepo := bootstrap.SetupUserRepository(a.db)
	csrfService := bootstrap.SetupCSRFService(a.redisClient)

	// Inject repositories and services into their respective application logic layers.
	userServiceLogin, userServiceRegister := bootstrap.SetupUserService(userRepo, csrfService)
	commentGetService, commentAddService := bootstrap.SetupCommentService(a.db)

	return &Dependencies{
		UserServiceLogin:    userServiceLogin,
		UserServiceRegister: userServiceRegister,
		CommentGetService:   commentGetService,
		CommentAddService:   commentAddService,
		RateHandler:         bootstrap.SetupRateLimiter(a.config),
		StaticFileAdapter:   bootstrap.SetupStaticFileAdapter(a.config),
		ProductsGetService:  bootstrap.SetupProductsService(a.db),
		CSRFMiddleware:      bootstrap.SetupCSRFMiddleware(csrfService),
		CSRFService:         csrfService,
	}
}

// BuildRouter constructs the final HTTP handler for the application.
// This method serves as the Composition Root, where dependencies are resolved and passed into the primary HTTP adapter (NewRouter).
// Returns:
//   - http.Handler: a fully configured router ready to serve requests.
func (a *Application) BuildRouter() http.Handler {
	// Step 1: Resolve all dependencies.
	deps := a.BuildDependencies()

	// Step 2: Inject dependencies into the HTTP router factory.
	return primaryHttp.NewRouter(
		deps.UserServiceLogin,
		deps.UserServiceRegister,
		deps.CommentGetService,
		deps.CommentAddService,
		deps.RateHandler,
		deps.StaticFileAdapter,
		deps.ProductsGetService,
		deps.CSRFMiddleware,
		deps.CSRFService,
	)
}
