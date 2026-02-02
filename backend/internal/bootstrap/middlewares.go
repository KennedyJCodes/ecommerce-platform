package bootstrap

import (
	"time"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/adapters/primary/http/middleware"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/adapters/secondary/repository/redis"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/config"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/services/service_csrf"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/rate_limiter"
	"github.com/redis/go-redis/v9"
)

func SetupRateLimiter(appConfig *config.AppConfig) ratelimiter.RateLimiterHandler {
	limiterConfig := appConfig.GetRateLimitConfig()
	return ratelimiter.NewDefaultRateLimiter(limiterConfig.RequestPerSecond, limiterConfig.Burst)
}

func SetupCSRFService(redisClient *redis.Client) input.CSRFService {
	csrfRepo := repository_redis.NewRedisCSRFRepository(redisClient)
	return service_csrf.NewCSRFUseCase(csrfRepo, 24*time.Hour)
}

func SetupCSRFMiddleware(csrfService input.CSRFService) *middleware.CSRFMiddleware {
	if csrfUseCase, ok := csrfService.(*service_csrf.CSRFUseCase); ok {
		return middleware.NewCSRFMiddleware(csrfUseCase)
	}
	return nil
}