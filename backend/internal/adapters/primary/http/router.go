// Package http implements HTTP handlers and the routing configuration for the ecommerce-platform application.
// It provides the RouterConfig type and NewRouter factory function, which wire up handlers, middlewares, rate limiters, and static file serving to produce a fully configured *mux.Router* ready to handle API requests.
package http

import (
	"net/http"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/primary/http/middleware"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	ratelimiter "github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/security/rate_limiter"
	"github.com/gorilla/mux"
)

// RouterConfiguration defines the interface for configuring routes in the application.
type RouterConfiguration interface {
	SetupRoutes(router *mux.Router)
}

// RouterConfig holds dependencies required to set up application routes.
// Fields correspond to handlers for each endpoint, middleware manager, and rate limiter components.
//
// Fields:
//   - IPExtractor: extracts client IP from *http.Request* for rate limiting.
//   - RateLimiter: handles request rate limiting based on extracted IP.
//   - LoginHandler: processes user login requests.
//   - RegisterHandler: processes user registration requests.
//   - CommentsGetHandler: handles retrieval of comments.
//   - CommentsAddHandler: handles creation of new comments.
//   - LogoutHandler: handles user logout and token revocation.
//   - MainPageHandler: serves the application's main page.
//   - StaticFileHandler: serves static assets like CSS/JS/images.
//   - MiddlewareManager: orchestrates application of global and route-specific middleware.
//   - BlacklistRepo: repository for checking revoked JWT tokens.
type RouterConfig struct {
	IPExtractor        ratelimiter.IPExtractor
	RateLimiter        ratelimiter.RateLimiterHandler
	LoginHandler       *LoginHandler
	RegisterHandler    *RegisterHandler
	CommentsGetHandler *CommentsGetHandler
	CommentsAddHandler *CommentsAddHandler
	LogoutHandler      *LogoutHandler
	MainPageHandler    *MainPageHandler
	StaticFileHandler  *StaticFileHandler
	MiddlewareManager  *middleware.MiddlewareManager
	ProductsHandler    *ProductsHandler
	CSRFMiddleware     *middleware.CSRFMiddleware
	BlacklistRepo      output.TokenBlacklistPort
}

// SetupRoutes registers all application endpoints on the given router and applies route-specific middleware for authentication and rate limiting.
// Routes include:
//   - Static files (CSS, JS, images)
//   - Public endpoints: GET /, POST /register, POST /login
//   - Protected endpoints: GET /comments, POST /comments/newComments, POST /logout

// Each route is wrapped with authentication and rate limiting via the MiddlewareManager.Apply method.

// Parameters:
//   - router: *mux.Router instance to configure routes on.
func (c *RouterConfig) SetupRoutes(router *mux.Router) {
	c.StaticFileHandler.RegisterRoutes(router)

	rateLimitMW := middleware.RateLimitMiddleware(c.IPExtractor, c.RateLimiter)
	authOpts := middleware.DefaultAuthOptions()
	authOpts.BlacklistRepo = c.BlacklistRepo
	authMW := middleware.AuthMiddleware(authOpts)
	csrfMW := c.CSRFMiddleware.ProtectCR

	public := router.PathPrefix("").Subrouter()
	public.Use(mux.MiddlewareFunc(middleware.CORSMiddleware(middleware.PublicCORSConfig())))
	public.Use(mux.MiddlewareFunc(rateLimitMW))

	public.Handle("/", http.HandlerFunc(c.MainPageHandler.Handle)).Methods("GET", "OPTIONS")
	public.Handle("/comments", http.HandlerFunc(c.CommentsGetHandler.Handle)).Methods("GET", "OPTIONS")
	public.Handle("/products", http.HandlerFunc(c.ProductsHandler.Handle)).Methods("GET", "OPTIONS")
	public.Handle("/product-id/{id}", http.HandlerFunc(c.ProductsHandler.HandleGetByID)).Methods("GET", "OPTIONS")
	public.Handle("/products-brand/{brand}", http.HandlerFunc(c.ProductsHandler.HandleGetByBrand)).Methods("GET", "OPTIONS")

	private := router.PathPrefix("").Subrouter()
	private.Use(mux.MiddlewareFunc(middleware.CORSMiddleware(middleware.PrivateCORSConfig())))
	private.Use(mux.MiddlewareFunc(rateLimitMW))
	private.Use(mux.MiddlewareFunc(authMW))
	private.Use(mux.MiddlewareFunc(csrfMW))

	private.Handle("/comments/newComments", http.HandlerFunc(c.CommentsAddHandler.Handle)).Methods("POST", "OPTIONS")
	private.Handle("/logout", http.HandlerFunc(c.LogoutHandler.Handle)).Methods("POST", "OPTIONS")

	loginRouter := router.PathPrefix("").Subrouter()
	loginRouter.Use(mux.MiddlewareFunc(middleware.CORSMiddleware(middleware.PrivateCORSConfig())))
	loginRouter.Use(mux.MiddlewareFunc(rateLimitMW))
	loginRouter.Handle("/login", http.HandlerFunc(c.LoginHandler.Handle)).Methods("POST", "OPTIONS")
	loginRouter.Handle("/register", http.HandlerFunc(c.RegisterHandler.Handle)).Methods("POST", "OPTIONS")
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
func NewRouter(
	userServiceLogin input.UserServiceLogin,
	userServiceRegister input.UserServiceRegister,
	commentGetService input.CommentGetService,
	commentAddService input.CommentAddService,
	rateHandler ratelimiter.RateLimiterHandler,
	staticFileService output.StaticFilePort,
	productsGetService input.ProductsGetService,
	csrfMiddleware *middleware.CSRFMiddleware,
	csrfService input.CSRFService,
	isProduction bool,
	blacklistRepo output.TokenBlacklistPort,
) *mux.Router {
	// 1. Initialize a new router
	router := mux.NewRouter()

	// 2. Instantiate HTTP handlers with injected domain services
	loginHandler := NewLoginHandler(userServiceLogin, csrfService, isProduction)
	registerHandler := NewRegisterHandler(userServiceRegister, csrfService, isProduction)
	commentsGetHandler := NewCommentsGetHandler(commentGetService)
	commentsAddHandler := NewCommentAddsHandler(commentAddService)
	logoutHandler := NewLogoutHandler(blacklistRepo, isProduction)
	mainPageHandler := NewMainPageHandler()
	staticFileHandler := NewStaticFileHandler(staticFileService)
	productsHandler := NewProductsHandler(productsGetService)

	// 3. Configure main page handler with static directory
	mainPageHandler.SetStaticDir(staticFileService.GetStaticDir())

	// 4. Create and configure MiddlewareManager
	middlewareManager := middleware.NewMiddlewareManager()
	timingConfig := middleware.DefaultTimingConfig()
	timingConfig.WarningThreshold = 200 * 1000 * 1000 // 200 milliseconds

	// Add global middleware: logging, timing, CORS
	middlewareManager.AddGlobal(middleware.LoggingMiddleware)
	middlewareManager.AddGlobal(middleware.TimingMiddleware(timingConfig))
	middlewareManager.ApplyToRouter(router)

	// 5. Build RouterConfig with dependencies
	config := &RouterConfig{
		IPExtractor:        ratelimiter.NewDefaultIPExtractor("10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "127.0.0.1/32"),
		RateLimiter:        rateHandler,
		LoginHandler:       loginHandler,
		RegisterHandler:    registerHandler,
		CommentsGetHandler: commentsGetHandler,
		CommentsAddHandler: commentsAddHandler,
		LogoutHandler:      logoutHandler,
		BlacklistRepo:      blacklistRepo,
		MainPageHandler:    mainPageHandler,
		StaticFileHandler:  staticFileHandler,
		MiddlewareManager:  middlewareManager,
		ProductsHandler:    productsHandler,
		CSRFMiddleware:     csrfMiddleware,
	}

	// 6. Register routes on router
	config.SetupRoutes(router)
	return router
}
