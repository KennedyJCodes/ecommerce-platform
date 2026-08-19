// Package middleware provides HTTP middleware utilities.
// This file contains a middleware that applies anti-caching headers to protect
// sensitive responses from being stored in browser or intermediary caches.
package middleware

import (
	"net/http"
)

// NoCacheMiddleware returns a Middleware that sets strict Cache-Control headers
// to prevent browsers and proxies from caching the response. It should be
// applied to any route that returns sensitive or authenticated data.
func NoCacheMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			next.ServeHTTP(w, r)
		})
	}
}
