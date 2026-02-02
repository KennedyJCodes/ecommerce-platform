// Package main provides the entry point for the Watch Store API server.

// It is responsible for loading configuration, initializing all application components, and starting the HTTP server to handle incoming API requests.

// The application follows a clean architecture pattern, separating concerns into adapters, domain services, ports, and infrastructure components.
package main

import (
	"log"
	"net/http"

	primaryHttp "github.com/David-Alejandro-Jimenez/sale-watches/internal/adapters/primary/http"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/bootstrap"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/bootstrap/database"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/unrolled/secure"
)

// main is the application entry point.
// It performs the following steps:
// 1. Loads and validates application configuration.
// 2. Initializes global security services (e.g., JWT).
// 3. Establishes a database connection.
// 4. Creates domain services and their dependencies (repositories, validators).
// 5. Configures the HTTP router with endpoints and middleware.
// 6. Starts listening on the configured port.

// If any of these steps fails, main will log the error and exit the application.
func main() {
	// Step 1: Load and validate configuration
	appConfig := config.NewAppConfig()
	config.ValidateRequiredConfig(appConfig.GetConfig())

	// Step 2: Initialize global services (e.g., JWT auth)
	bootstrap.SetupCommonServices(appConfig)

	// Step 3: Database setup
	db, err := bootstrap_database.SetupDatabaseMySQL(appConfig)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	redisConfig, err := bootstrap_database.SetupDatabaseRedis(appConfig)
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	redisClient := config.NewRedisClient(redisConfig)
	defer redisClient.Close()

	// Step 4: Dependency injection for domain services
	userRepo := bootstrap.SetupUserRepository(db)
	csrfService := bootstrap.SetupCSRFService(redisClient)
	userServiceLogin, userServiceRegister := bootstrap.SetupUserService(userRepo, csrfService)
	commentGetService, commentAddService := bootstrap.SetupCommentService(db)
	rateHandler := bootstrap.SetupRateLimiter(appConfig)
	staticFileAdapter := bootstrap.SetupStaticFileAdapter(appConfig)
	productsGetService := bootstrap.SetupProductsService(db)
	csrfMiddleware := bootstrap.SetupCSRFMiddleware(csrfService)

	// Step 5: Configure HTTP router with handlers and middleware
	router := primaryHttp.NewRouter(
		userServiceLogin,
		userServiceRegister,
		commentGetService,
		commentAddService,
		rateHandler,
		staticFileAdapter,
		productsGetService,
		csrfMiddleware,
		csrfService,
	)

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

	// Endpoints sensibles que no deben ser cacheados
	sensitivePaths := map[string]bool{
		"/login":                true,
		"/register":             true,
		"/comments/newComments": true,
	}

	securedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	// Step 6: Start HTTP server
	port := appConfig.GetPort()
	log.Printf("Server started at http://localhost:%s", port)
	log.Printf("Serving static files from: %s", staticFileAdapter.GetStaticDir())
	log.Fatal(http.ListenAndServe(":"+port, securedHandler))
}
