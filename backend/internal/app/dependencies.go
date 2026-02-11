package app

import (
	"net/http"

	primaryHttp "github.com/David-Alejandro-Jimenez/sale-watches/internal/adapters/primary/http"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/adapters/primary/http/middleware"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/bootstrap"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/output"
	ratelimiter "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/rate_limiter"
)

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

func (a *Application) BuildDependencies() *Dependencies {
    userRepo := bootstrap.SetupUserRepository(a.db)
    csrfService := bootstrap.SetupCSRFService(a.redisClient)
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

func (a *Application) BuildRouter() http.Handler {
    deps := a.BuildDependencies()
    
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