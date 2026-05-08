// Package middleware provides HTTP middleware utilities.
// This file contains a rate limiting middleware that leverages a custom rate
// limiter to restrict the number of requests per IP address.
package middleware

import (
	"net/http"

	ratelimiter "github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/security/rate_limiter"
)

// RateLimitMiddleware returns a Middleware that enforces rate limiting based on client IP.

// It accepts an IPExtractor to parse the client's IP address from a request and a
// RateLimiterHandler that defines the rate limiting behavior (such as requests per second and burst limits).

// The middleware extracts the client's IP from the request, checks with the rate limiter if the request is allowed, and if not, responds with an HTTP 429 (Too Many Requests) error.

// Otherwise, it forwards the request to the next handler in the chain.
func RateLimitMiddleware(ipExtractor ratelimiter.IPExtractor, limiter ratelimiter.RateLimiterHandler) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := ipExtractor.Extract(r.RemoteAddr)

			if !limiter.Allow(clientIP) {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
