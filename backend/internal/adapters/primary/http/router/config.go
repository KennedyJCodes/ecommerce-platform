package router

import (
	"time"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/primary/http/handlers/private"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/primary/http/handlers/public"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/primary/http/middleware"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/security/rate_limiter"
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
	Login       *public_handlers.LoginHandler
	Register    *public_handlers.RegisterHandler
	Refresh     *public_handlers.RefreshHandler
	ReviewsGet *public_handlers.ReviewsGetHandler
	ReviewsAdd *private_handlers.ReviewsAddHandler
	Logout      *private_handlers.LogoutHandler
	MainPage    *public_handlers.MainPageHandler
	StaticFile  *public_handlers.StaticFileHandler
	Products    *public_handlers.ProductsHandler
}

func buildHandlers(
	userServiceLogin input.UserServiceLogin,
	userServiceRegister input.UserServiceRegister,
	reviewGetService input.ReviewGetService,
	reviewAddService input.ReviewAddService,
	staticFileService output.StaticFilePort,
	productsGetService input.ProductsGetService,
	csrfService output.CSRFService,
	tokenService output.TokenService,
	blacklistRepo output.TokenBlacklistPort,
	isProduction bool,
) HandlerConfig {

	mainPageHandler := public_handlers.NewMainPageHandler()
	mainPageHandler.SetStaticDir(staticFileService.GetStaticDir())

	return HandlerConfig{
		Login:       public_handlers.NewLoginHandler(userServiceLogin, csrfService, isProduction),
		Register:    public_handlers.NewRegisterHandler(userServiceRegister, csrfService, isProduction),
		Refresh:     public_handlers.NewRefreshHandler(tokenService, blacklistRepo, isProduction),
		ReviewsGet: public_handlers.NewReviewsGetHandler(reviewGetService),
		ReviewsAdd: private_handlers.NewReviewsAddHandler(reviewAddService),
		Logout:      private_handlers.NewLogoutHandler(tokenService, blacklistRepo, isProduction),
		MainPage:    mainPageHandler,
		StaticFile:  public_handlers.NewStaticFileHandler(staticFileService),
		Products:    public_handlers.NewProductsHandler(productsGetService),
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