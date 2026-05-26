package ratelimiter

import (
	"testing"
	"time"
)

func TestCleanupInactiveLimitersRemovesExpiredEntries(t *testing.T) {
	manager := NewRateLimiterManager()
	ipAddress := "203.0.113.10"

	manager.GetRateLimiterForIP(ipAddress)
	record, ok := manager.ipLimiterCache.Load(ipAddress)
	if !ok {
		t.Fatalf("expected limiter record for IP %s", ipAddress)
	}

	record.(*LimiterEntry).touch(time.Now().Add(-1 * time.Hour))
	manager.CleanupInactiveLimiters(30 * time.Minute)

	if _, ok := manager.ipLimiterCache.Load(ipAddress); ok {
		t.Fatalf("expected expired limiter record for IP %s to be removed", ipAddress)
	}
}

func TestCleanupInactiveLimitersKeepsActiveEntries(t *testing.T) {
	manager := NewRateLimiterManager()
	ipAddress := "203.0.113.11"

	manager.GetRateLimiterForIP(ipAddress)
	manager.CleanupInactiveLimiters(30 * time.Minute)

	if _, ok := manager.ipLimiterCache.Load(ipAddress); !ok {
		t.Fatalf("expected active limiter record for IP %s to be kept", ipAddress)
	}
}
