// Package ratelimiter provides background cleanup utilities for inactive rate limiters.
// It implements a maintenance service that periodically purges expired rate limiting records to prevent memory leaks and optimize storage usage.
package ratelimiter

import "time"

// RateLimiterCleaner manages background cleanup of inactive rate limiter entries.
// Maintains a reference to a RateLimiterManager to perform periodic cleanup operations.
// Should be instantiated once per application lifecycle.
type RateLimiterCleaner struct {
	manager RateLimiterManager
}

// NewRateLimiterCleaner creates a new cleanup service instance.
// Requires a RateLimiterManager implementation that provides the CleanupInactiveLimiters method.
// Typical usage:
//
//	cleaner := NewRateLimiterCleaner(rateLimiterManager)
//	cleaner.Start(10*time.Minute, 1*time.Minute)
func NewRateLimiterCleaner(manager RateLimiterManager) *RateLimiterCleaner {
	return &RateLimiterCleaner{manager: manager}
}

// Start begins the background cleanup goroutine with specified intervals.
// Parameters:
//   expirationDuration - Time since last access after which limiters are considered inactive
//   cleanupInterval - Frequency between cleanup cycles

// The cleanup runs indefinitely until application shutdown.
// Example: Start(15*time.Minute, 5*time.Minute) cleans every 5 minutes, removing limiters inactive for 15+ minutes
func (c *RateLimiterCleaner) Start(expirationDuration, cleanupInterval time.Duration) {
	if c == nil || c.manager == nil {
		return
	}
	if expirationDuration <= 0 {
		expirationDuration = 30 * time.Minute
	}
	if cleanupInterval <= 0 {
		cleanupInterval = 5 * time.Minute
	}

	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()
		for range ticker.C {
			c.manager.CleanupInactiveLimiters(expirationDuration)
		}
	}()
}
