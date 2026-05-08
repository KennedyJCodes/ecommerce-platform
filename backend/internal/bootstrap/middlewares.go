// Package bootstrap provides high-level factory functions to initialize and wire the application's infrastructure components.
// This file specifically handles the setup of security-related services, including rate limiting and CSRF protection mechanisms.
package bootstrap

import (
	"time"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/primary/http/middleware"
	repository_redis "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/secondary/repository/redis"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/config"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/services/service_csrf"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	ratelimiter "github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/security/rate_limiter"
	"github.com/redis/go-redis/v9"
)

// SetupRateLimiter initializes the rate limiting handler based on application configuration.
// It retrieves the threshold and burst settings from the config provider to create a limiter that protects the API from brute-force and DoS attacks.

// Parameters:
//   - appConfig: the global application configuration manager.
//
// Returns:
//   - ratelimiter.RateLimiterHandler: a configured rate limiter ready for middleware injection.
func SetupRateLimiter(appConfig *config.AppConfig) ratelimiter.RateLimiterHandler {
	limiterConfig := appConfig.GetRateLimitConfig()
	return ratelimiter.NewDefaultRateLimiter(limiterConfig.RequestPerSecond, limiterConfig.Burst)
}

// SetupCSRFService initializes the CSRF domain service using Redis for token persistence.
// It wires the Redis repository with the CSRF Use Case, setting a default token expiration of 24 hours.

// Parameters:
//   - redisClient: an active *redis.Client connection.
//
// Returns:
//   - input.CSRFService: the initialized service implementing the CSRF input port.
func SetupCSRFService(redisClient *redis.Client) input.CSRFService {
	csrfRepo := repository_redis.NewRedisCSRFRepository(redisClient)
	// Initialize the Use Case with the repository and a 24-hour TTL for tokens.
	return service_csrf.NewCSRFUseCase(csrfRepo, 24*time.Hour)
}

// SetupCSRFMiddleware creates the HTTP middleware responsible for CSRF validation.
// It performs a type assertion to ensure the provided CSRF service is compatible with the middleware requirements.

// Parameters:
//   - csrfService: the CSRF service instance (input port).
//
// Returns:
//   - *middleware.CSRFMiddleware: the configured middleware, or nil if the type assertion fails.
func SetupCSRFMiddleware(csrfService input.CSRFService) *middleware.CSRFMiddleware {
	// Type assertion to bridge the input port with the specific middleware implementation.
	if csrfUseCase, ok := csrfService.(*service_csrf.CSRFUseCase); ok {
		return middleware.NewCSRFMiddleware(csrfUseCase)
	}
	return nil
}
