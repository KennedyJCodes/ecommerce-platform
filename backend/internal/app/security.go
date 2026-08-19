// Package app provides the main application components and security wrappers.
// This file focuses on global security middleware, applying headers that protect  against common web vulnerabilities like XSS, Clickjacking, and sensitive data caching.
package app

import (
	"log"
	"net/http"

	"github.com/unrolled/secure"
)

// WrapWithSecurityMiddleware applies a robust layer of security headers to the provided router.
// It uses the 'secure' package to set standard headers (CSP, HSTS, etc.) and supplements them with manual headers for Permissions Policy and Cache Control.

// Security Features Applied:
//  - Content Security Policy (CSP): Restricts resource loading to 'self' (same origin).
//  - HSTS (STS): Forces HTTPS for a year (31536000 seconds).
//  - X-Frame-Options: Set to DENY to prevent Clickjacking.
//  - No-Sniff: Prevents the browser from MIME-sniffing away from the declared Content-Type.
//  - Referrer-Policy: Protects privacy by limiting the referrer information sent.
// Parameters:
//   - router: The primary http.Handler (usually the mux router) to be wrapped.
//   - isDevelopment: Whether the app runs in development mode (relaxes some protections).
// Returns:
//   - http.Handler: A new handler protected by the security layer.
func WrapWithSecurityMiddleware(router http.Handler, isDevelopment bool) http.Handler {

	secureMiddleware := secure.New(secure.Options{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		STSSeconds:            31536000,
		STSIncludeSubdomains:  true,
		SSLRedirect:           false, // Handled by proxy/infrastructure in production
		IsDevelopment:         isDevelopment,
		ContentSecurityPolicy: "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data: https:; font-src 'self'; connect-src 'self';",
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Process standard security headers from the 'secure' package.
		err := secureMiddleware.Process(w, r)
		if err != nil {
			log.Println("Error processing security headers:", err)
			return
		}

		// Apply supplementary modern security headers.
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), usb=(), payment=(), accelerometer=(), gyroscope=(), magnetometer=(), clipboard-read=(), clipboard-write=(), fullscreen=(self)")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")

		// Pass control to the next handler in the chain.
		router.ServeHTTP(w, r)
	})
}