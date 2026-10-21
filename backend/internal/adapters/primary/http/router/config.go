package router

import (
	"time"

	"github.com/KennedyJCodes/ecommerce-platform/internal/adapters/primary/http/handlers/private"
	"github.com/KennedyJCodes/ecommerce-platform/internal/adapters/primary/http/handlers/public"
	"github.com/KennedyJCodes/ecommerce-platform/internal/adapters/primary/http/middleware"
	"github.com/KennedyJCodes/ecommerce-platform/internal/core/ports/input"
	"github.com/KennedyJCodes/ecommerce-platform/internal/core/ports/output"
	"github.com/KennedyJCodes/ecommerce-platform/pkg/security/rate_limiter"
	"github.com/gorilla/mux"
)

type RouterConfig struct {
	Handlers          HandlerConfig
	IPExtractor       ratelimiter.IPExtractor
	RateLimiter       ratelimiter.RateLimiterHandler
	MiddlewareManager *middleware.MiddlewareManager
	CSRFMiddleware    *middleware.CSRFMiddleware
	TokenService      output.TokenService
	BlacklistRepo     output.TokenBlacklistPort
}

type HandlerConfig struct {
	Login      *public_handlers.LoginHandler
	Register   *public_handlers.RegisterHandler
	Verify     *public_handlers.VerifyHandler
	Refresh    *public_handlers.RefreshHandler
	ReviewsGet *public_handlers.ReviewsGetHandler
	ReviewsAdd *private_handlers.ReviewsAddHandler
	Logout     *private_handlers.LogoutHandler
	MainPage   *public_handlers.MainPageHandler
	StaticFile *public_handlers.StaticFileHandler
	Products   *public_handlers.ProductsHandler
}

// HandlerDependencies contains the dependencies required to construct all HTTP handlers.
// Keeping this separate from RouterDependencies makes buildHandlers explicit about the
// dependencies it needs without passing unrelated router configuration through it.
type HandlerDependencies struct {
	UserServiceLogin        input.UserServiceLogin
	UserServiceRegister     input.UserServiceRegister
	UserVerificationService input.UserVerificationService
	ReviewGetService        input.ReviewGetService
	ReviewAddService        input.ReviewAddService
	StaticFileService       output.StaticFilePort
	ProductsGetService      input.ProductsGetService
	CSRFService             output.CSRFService
	TokenService            output.TokenService
	BlacklistRepo           output.TokenBlacklistPort
	IsProduction            bool
}

func buildHandlers(deps HandlerDependencies) HandlerConfig {

	mainPageHandler := public_handlers.NewMainPageHandler()
	mainPageHandler.SetStaticDir(deps.StaticFileService.GetStaticDir())

	return HandlerConfig{
		Login:      public_handlers.NewLoginHandler(deps.UserServiceLogin, deps.CSRFService, deps.IsProduction),
		Register:   public_handlers.NewRegisterHandler(deps.UserServiceRegister, deps.IsProduction),
		Verify:     public_handlers.NewVerifyHandler(deps.UserVerificationService, deps.IsProduction),
		Refresh:    public_handlers.NewRefreshHandler(deps.TokenService, deps.BlacklistRepo, deps.IsProduction),
		ReviewsGet: public_handlers.NewReviewsGetHandler(deps.ReviewGetService),
		ReviewsAdd: private_handlers.NewReviewsAddHandler(deps.ReviewAddService),
		Logout:     private_handlers.NewLogoutHandler(deps.TokenService, deps.BlacklistRepo, deps.IsProduction),
		MainPage:   mainPageHandler,
		StaticFile: public_handlers.NewStaticFileHandler(deps.StaticFileService),
		Products:   public_handlers.NewProductsHandler(deps.ProductsGetService),
	}
}

func configureMiddleware(router *mux.Router) *middleware.MiddlewareManager {
	manager := middleware.NewMiddlewareManager()

	timingConfig := middleware.DefaultTimingConfig()
	timingConfig.WarningThreshold = 200 * time.Millisecond

	manager.AddGlobal(middleware.LoggingMiddleware)
	manager.AddGlobal(middleware.TimingMiddleware(timingConfig))

	manager.ApplyToRouter(router)

	return manager
}
