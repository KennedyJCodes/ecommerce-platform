// Package ratelimiter provides IP-based request rate limiting functionality.
// It implements a token bucket algorithm with automatic cleanup of inactive clients.
// Features include:
// - Per-IP rate limiting
// - Configurable request rates and burst sizes
// - Background cleanup of inactive clients
// - Thread-safe operations
package ratelimiter

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models"
	"golang.org/x/time/rate"
)

// RateLimiterHandler defines the core rate limiting interface for request authorization.
type RateLimiterHandler interface {
	// Allow checks if a request from the specified IP address is permitted.
	// Returns true if the request should be allowed, false if rate limited.
	Allow(ipAddress string) bool
}

// RateLimiterManager defines operations for managing rate limiter instances.
type RateLimiterManager interface {
	// GetRateLimiterForIP retrieves or creates a rate limiter for the specified IP.
	GetRateLimiterForIP(ipAddress string) *rate.Limiter

	// CleanupInactiveLimiters removes rate limiters that haven't been used within expirationDuration.
	CleanupInactiveLimiters(expirationDuration time.Duration)

	// SetDefaultLimiterConfig updates the default rate limiting parameters.
	SetDefaultLimiterConfig(config models.LimiterConfig)
}

// DefaultRateLimiter implements RateLimiterHandler using an underlying RateLimiterManager.
// Provides basic rate limiting capabilities with per-IP tracking.
type DefaultRateLimiter struct {
	manager RateLimiterManager
}

// NewDefaultRateLimiter creates a new rate limiter with specified default configuration.
// requestPerSecond: Base allowed request rate (token refresh rate)
// burst: Maximum burst size (bucket capacity)
func NewDefaultRateLimiter(requestPerSecond float64, burst int) RateLimiterHandler {
	manager := NewRateLimiterManager()
	manager.SetDefaultLimiterConfig(models.LimiterConfig{
		RequestPerSecond: requestPerSecond,
		Burst:            burst,
	})

	return &DefaultRateLimiter{
		manager: manager,
	}
}

// NewDefaultRateLimiterWithCleanup creates a rate limiter and starts background cleanup
// for inactive IP entries.
func NewDefaultRateLimiterWithCleanup(config models.LimiterConfig) RateLimiterHandler {
	manager := NewRateLimiterManager()
	manager.SetDefaultLimiterConfig(config)

	cleaner := NewRateLimiterCleaner(manager)
	cleaner.Start(config.ExpirationDuration, config.CleanupInterval)

	return &DefaultRateLimiter{
		manager: manager,
	}
}

// Allow implements rate limiting check for the specified IP address.
// Consumes one token from the IP's rate limiter bucket.
func (d *DefaultRateLimiter) Allow(ipAddress string) bool {
	limiter := d.manager.GetRateLimiterForIP(ipAddress)
	return limiter.Allow()
}

// DefaultRateLimiterManager implements RateLimiterManager with in-memory storage.
// Uses sync.Map for concurrent access safety and automatic cleanup of inactive entries.
type DefaultRateLimiterManager struct {
	defaultConfig  models.LimiterConfig // Default rate limiting parameters
	ipLimiterCache sync.Map             // Concurrent storage for IP limiter records
}

// NewRateLimiterManager creates a new rate limiter manager with default configuration:
// - 1.5 requests/second base rate
// - 5 request burst capacity
func NewRateLimiterManager() *DefaultRateLimiterManager {
	return &DefaultRateLimiterManager{
		defaultConfig: models.LimiterConfig{
			RequestPerSecond: 1.5,
			Burst:            5,
		},
	}
}

// SetDefaultLimiterConfig updates the default rate limiting parameters for new limiters.
// Does not affect existing rate limiter instances.
func (m *DefaultRateLimiterManager) SetDefaultLimiterConfig(config models.LimiterConfig) {
	m.defaultConfig = config
}

// GetRateLimiterForIP retrieves or creates a rate limiter for the specified IP.
// Updates last access time for cleanup tracking. Existing limiters are reused.
func (m *DefaultRateLimiterManager) GetRateLimiterForIP(ipAddress string) *rate.Limiter {
	currentTime := time.Now()
	record, exists := m.ipLimiterCache.Load(ipAddress)
	if exists {
		limiterRecord := record.(*LimiterEntry)
		limiterRecord.touch(currentTime)
		return limiterRecord.limiter
	}

	newLimiter := rate.NewLimiter(rate.Limit(m.defaultConfig.RequestPerSecond), m.defaultConfig.Burst)
	newRecord := &LimiterEntry{
		limiter: newLimiter,
	}
	newRecord.touch(currentTime)

	actualRecord, loaded := m.ipLimiterCache.LoadOrStore(ipAddress, newRecord)
	if loaded {
		limiterRecord := actualRecord.(*LimiterEntry)
		limiterRecord.touch(currentTime)
		return limiterRecord.limiter
	}
	return newLimiter
}

// CleanupInactiveLimiters removes rate limiters that haven't been accessed within expirationDuration.
// Typically run periodically via a background goroutine.
func (m *DefaultRateLimiterManager) CleanupInactiveLimiters(expirationDuration time.Duration) {
	currentTime := time.Now()
	m.ipLimiterCache.Range(func(key, value interface{}) bool {
		limiterRecord := value.(*LimiterEntry)
		if currentTime.Sub(limiterRecord.lastSeen()) > expirationDuration {
			m.ipLimiterCache.Delete(key)
		}
		return true
	})
}

// LimiterEntry represents a rate limiter instance with last access timestamp.
type LimiterEntry struct {
	limiter          *rate.Limiter // Token bucket rate limiter instance
	lastSeenUnixNano atomic.Int64  // Last access time for cleanup tracking
}

func (e *LimiterEntry) touch(t time.Time) {
	e.lastSeenUnixNano.Store(t.UnixNano())
}

func (e *LimiterEntry) lastSeen() time.Time {
	return time.Unix(0, e.lastSeenUnixNano.Load())
}
