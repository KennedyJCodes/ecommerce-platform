package app

import (
	"log"
	"net/http"

	"github.com/unrolled/secure"
)

var sensitivePaths = map[string]bool{
	"/login":                true,
	"/register":             true,
	"/comments/newComments": true,
}

// WrapWithSecurityMiddleware envuelve el handler con headers de seguridad
func WrapWithSecurityMiddleware(router http.Handler) http.Handler {
	secureMiddleware := secure.New(secure.Options{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		STSSeconds:            31536000,
		STSIncludeSubdomains:  true,
		SSLRedirect:           false,
		IsDevelopment:         true,
		ContentSecurityPolicy: "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data: https:; font-src 'self'; connect-src 'self';",
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := secureMiddleware.Process(w, r)
		if err != nil {
			log.Println("Error processing security headers:", err)
			return
		}

		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), usb=(), payment=(), accelerometer=(), gyroscope=(), magnetometer=(), clipboard-read=(), clipboard-write=(), fullscreen=(self)")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")

		if sensitivePaths[r.URL.Path] {
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		}

		router.ServeHTTP(w, r)
	})
}