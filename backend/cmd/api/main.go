// Package main provides the entry point for the Watch Store API server.

// It is responsible for loading configuration, initializing all application components, and starting the HTTP server to handle incoming API requests.

// The application follows a clean architecture pattern, separating concerns into adapters, domain services, ports, and infrastructure components.
package main

import (
	"fmt"
	"log"
	"net/http"

	primaryHttp "github.com/David-Alejandro-Jimenez/sale-watches/internal/adapters/primary/http"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/adapters/secondary/repository"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/adapters/secondary/static"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/config"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/services/service_auth"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/services/service_comments"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/output"
	ratelimiter "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/rate_limiter"
	securityAuth "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/security_auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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
	initializeCommonServices(appConfig)

	// Step 3: Database setup
	db, err := setupDatabase(appConfig)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Step 4: Dependency injection for domain services
	userRepo := setupUserRepository(db)
	userServiceLogin := setupLoginService(userRepo)
	userServiceRegister := setupRegisterService(userRepo)
	commentGetService, commentAddService := setupCommentService(db)
	rateHandler := setupRateLimiter(appConfig)
	staticFileAdapter := setupStaticFileAdapter(appConfig)

	// Step 5: Configure HTTP router with handlers and middleware
	router := primaryHttp.NewRouter(
		userServiceLogin,
		userServiceRegister,
		commentGetService,
		commentAddService,
		rateHandler,
		staticFileAdapter,
	)

	// Step 6: Start HTTP server
	port := appConfig.GetPort()
	log.Printf("Server started at http://localhost:%s", port)
	log.Printf("Serving static files from: %s", staticFileAdapter.GetStaticDir())
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// initializeCommonServices sets up services that are shared globally across the application.

// Currently, this function initializes the default JWT authentication service using the secret key from configuration.
func initializeCommonServices(appConfig *config.AppConfig) {
	securityAuth.SetDefaultJWTService(appConfig.GetJWTSecret())
}

// setupDatabase establishes a connection to the MySQL database.

// It uses configuration values such as username, password, host, and database name to construct the DSN string and open the connection. It returns a *sqlx.DB instance and an error if the connection fails.
func setupDatabase(appConfig *config.AppConfig) (*sqlx.DB, error) {
	cfg := appConfig.GetConfig()
	user := cfg.GetString("database.user")
	password := cfg.GetString("database.password")
	host := cfg.GetString("database.host")
	port := cfg.GetInt("database.port")
	dbName := cfg.GetString("database.name")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, password, host, port, dbName)
	return sqlx.Connect("mysql", dsn)
}

// setupUserRepository returns an implementation of the UserRepository interface.

// It sets up dependencies for user authentication such as salt generation and password hashing and injects them into the SQL-based repository.
func setupUserRepository(db *sqlx.DB) output.UserRepository {
	hasher := securityAuth.BcryptHasher{}
	return repository.NewSQLUserRepository(db, hasher)
}

// setupLoginService initializes and returns the user login service.

// This service validates credentials and authenticates users.
// It relies on validators for username and password and uses the user repository to query user data.
func setupLoginService(userRepo output.UserRepository) input.UserServiceLogin {
	userNameValidator := &service_auth.UserNameValidator{}
	passwordValidator := &service_auth.PasswordValidator{}
	return service_auth.NewUserLoginService(userRepo, userNameValidator, passwordValidator)
}

// setupRegisterService initializes and returns the user registration service.
// It validates user input and stores new users in the database using the provided repository.
func setupRegisterService(userRepo output.UserRepository) input.UserServiceRegister {
	userNameValidator := &service_auth.UserNameValidator{}
	passwordValidator := &service_auth.PasswordValidator{}
	return service_auth.NewUserRegisterService(userRepo, userNameValidator, passwordValidator)
}

// setupCommentService initializes services for retrieving and creating user comments.
// This binds the comment repository and validation rules into service implementations.
// Parameters:
//   - db: active *sqlx.DB connection

// Returns:
//   - input.CommentGetService: service interface to fetch comments
//   - input.CommentAddService: service interface to add new comments
func setupCommentService(db *sqlx.DB) (input.CommentGetService, input.CommentAddService) {
	commentRepo := repository.NewSqlCommentRepository(db)
	commentValidator := &service_comments.CommentValidator{}
	return service_comments.NewCommentGetService(commentRepo, commentValidator), service_comments.NewCommentAddService(commentRepo, commentValidator)
}

// setupRateLimiter configures and returns a rate limiting handler.
// It uses rate limit settings (requests per second and burst) defined in the application configuration to protect the API against abuse or DoS attacks.
func setupRateLimiter(appConfig *config.AppConfig) ratelimiter.RateLimiterHandler {
	limiterConfig := appConfig.GetRateLimitConfig()
	return ratelimiter.NewDefaultRateLimiter(limiterConfig.RequestPerSecond, limiterConfig.Burst)
}

// setupStaticFileAdapter creates and returns an adapter for serving static files.

// The adapter serves assets such as images, stylesheets, or JavaScript files
// from a directory defined in the configuration.
func setupStaticFileAdapter(appConfig *config.AppConfig) output.StaticFilePort {
	staticDir := appConfig.GetStaticDir()
	return static.NewStaticFileAdapter(staticDir)
}
