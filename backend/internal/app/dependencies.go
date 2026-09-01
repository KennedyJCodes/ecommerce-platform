// Package app provides the main application entry point and dependency orchestration.
// This file focuses on assembling the core dependencies and wiring them into the primary HTTP adapter (the router).
package app

import (
	"fmt"
	"net/http"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/primary/http/middleware"
	primaryRouter "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/primary/http/router"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/bootstrap"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/http/cookies"
	ratelimiter "github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/security/rate_limiter"
)

// Dependencies holds all the initialized services, ports, and handlers required  to run the application.
// By grouping these in a single struct, the application ensures that all required components are ready before starting the HTTP server.
type Dependencies struct {
	UserServiceLogin    input.UserServiceLogin
	UserServiceRegister input.UserServiceRegister
	ReviewGetService   input.ReviewGetService
	ReviewAddService   input.ReviewAddService
	RateHandler         ratelimiter.RateLimiterHandler
	StaticFileAdapter   output.StaticFilePort
	ProductsGetService  input.ProductsGetService
	CSRFMiddleware      *middleware.CSRFMiddleware
	TokenService        output.TokenService
	CSRFService         output.CSRFService
	BlacklistRepo       output.TokenBlacklistPort
}

// BuildDependencies orchestrates the initialization of all internal services and repositories.
// It uses the bootstrap package to set up:
//  1. Repositories (User, Comments, Products).
//  2. Core Services (Auth, CSRF, Comments).
//  3. Infrastructure Adapters (Rate Limiter, Static Files).
//
// Returns:
//   - *Dependencies: a pointer to the fully populated Dependencies struct.
func (a *Application) BuildDependencies() (*Dependencies, error) {
	// Configure cookie prefix to avoid collisions on shared domains.
	cookies.SetCookiePrefix(a.config.GetCookiePrefix())

	// Initialize core repositories and services using shared database connections.
	userRepo, err := bootstrap.SetupUserRepository(a.db)
	if err != nil {
		return nil, fmt.Errorf("build dependencies: %w", err)
	}

	tokenService := bootstrap.SetupTokenService(a.config)
	csrfService := bootstrap.SetupCSRFService(a.redisClient)
	blacklistRepo := bootstrap.SetupTokenBlacklistRepository(a.redisClient)

	// Inject repositories and services into their respective application logic layers.
	userServiceLogin, userServiceRegister := bootstrap.SetupUserService(userRepo, tokenService, csrfService)
	reviewGetService, reviewAddService, err := bootstrap.SetupReviewService(a.db)
	if err != nil {
		return nil, fmt.Errorf("build dependencies: %w", err)
	}

	productsGetService, err := bootstrap.SetupProductsService(a.db)
	if err != nil {
		return nil, fmt.Errorf("build dependencies: %w", err)
	}

	return &Dependencies{
		UserServiceLogin:    userServiceLogin,
		UserServiceRegister: userServiceRegister,
		ReviewGetService:   reviewGetService,
		ReviewAddService:   reviewAddService,
		RateHandler:         bootstrap.SetupRateLimiter(a.config),
		StaticFileAdapter:   bootstrap.SetupStaticFileAdapter(a.config),
		ProductsGetService:  productsGetService,
		TokenService:        tokenService,
		CSRFMiddleware:      bootstrap.SetupCSRFMiddleware(csrfService, a.config.IsProduction()),
		CSRFService:         csrfService,
		BlacklistRepo:       blacklistRepo,
	}, nil
}

// BuildRouter constructs the final HTTP handler for the application.
// This method serves as the Composition Root, where dependencies are resolved and passed into the primary HTTP adapter (NewRouter).
// Returns:
//   - http.Handler: a fully configured router ready to serve requests.
func (a *Application) BuildRouter() (http.Handler, error) {
	// Step 1: Resolve all dependencies.
	deps, err := a.BuildDependencies()
	if err != nil {
		return nil, err
	}

	// Step 2: Inject dependencies into the HTTP router factory.
	router := primaryRouter.NewRouter(primaryRouter.RouterDependencies{
		UserServiceLogin:    deps.UserServiceLogin,
		UserServiceRegister: deps.UserServiceRegister,
		ReviewGetService:   deps.ReviewGetService,
		ReviewAddService:   deps.ReviewAddService,
		RateHandler:         deps.RateHandler,
		StaticFileService:   deps.StaticFileAdapter,
		ProductsGetService:  deps.ProductsGetService,
		CSRFMiddleware:      deps.CSRFMiddleware,
		CSRFService:         deps.CSRFService,
		IsProduction:        a.config.IsProduction(),
		BlacklistRepo:       deps.BlacklistRepo,
		TokenService:        deps.TokenService,
	})

	return router, nil
}
