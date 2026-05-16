// The package app provides startup and lifecycle management of the main application, including configuration loading, service initialization, external resource configuration, and graceful shutdown.
package app

import (
	"fmt"
	"net/http"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/bootstrap"
	bootstrap_database "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/bootstrap/database"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/config"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// The application encapsulates the application lifecycle:
// load configuration, configure common services, initialize
// the database and Redis, and clean up resources upon completion.

// The idea is that the startup order is:
// 1. Load configuration
// 2. Configure common services
// 3. Configure the database
// 4. Configure Redis
// and upon completion or during shutdown, call Close.
type Application struct {
	config      *config.AppConfig // app configuration (e.g. port, DSNs, flags)
	db          *sqlx.DB          // connection to the main database
	redisClient *redis.Client     // Redis client for cache/sessions
	httpServer  *http.Server      // HTTP server responsible for managing incoming requests and lifecycle control
}

// NewApplication creates and returns an empty Application instance.
// It does not perform any initialization that depends on the configuration.
// (That's why LoadConfig and the specific Setup commands run separately).
func NewConfigApplication() *Application {
	return &Application{}
}

// LoadConfig initializes the application configuration and validates
// the required parameters needed for startup.
// Must be run before any Setup that depends on the configuration.
// Side effects:
// - assigns a valid instance to a.config.
// Errors:
// - currently the implementation always returns nil; if you change to a
// validation that can fail, propagate the error appropriately.
func (a *Application) LoadConfig() error {
	a.config = config.NewAppConfig()
	config.ValidateRequiredConfig(a.config.GetConfig())
	return nil
}

// SetupCommonServices registers or configures shared services at the application level (for example: global loggers, dependency providers, common HTTP clients, etc.). It depends on LoadConfig having already been executed.
// It does not return an error in the current deployment, but if the services you register can fail, consider propagating errors to abort the startup.
func (a *Application) SetupCommonServices() error {
	bootstrap.SetupCommonServices(a.config)
	return nil
}

// SetupDatabase establishes the primary connection to the database (MySQL)
// and registers it with the application for later use.
// Returns an error if a connection cannot be established; the caller must decide whether to abort the startup or retry.
// Side effects:
// - assigns the obtained connection to a.db.
func (a *Application) SetupDatabase() error {
	db, err := bootstrap_database.SetupDatabaseMySQL(a.config)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	a.db = db
	return nil
}

// SetupRedis initializes the Redis connection/client used by the application
// (cache, sessions, lightweight queues, etc.). Returns an error if initialization fails.
// Side effects:
// - assigns a.redisClient a ready-to-use client.
func (a *Application) SetupRedis() error {
	redisConfig, err := bootstrap_database.SetupDatabaseRedis(a.config)
	if err != nil {
		return fmt.Errorf("error connecting to Redis: %w", err)
	}
	a.redisClient = config.NewRedisClient(redisConfig)
	return nil
}

// GetPort returns the port configured for the server. It is provided
// as a method to encapsulate the source of the value (internal configuration).
// It does not perform any validations; it assumes that LoadConfig has already populated the value.
func (a *Application) GetPort() string {
	return a.config.GetPort()
}

// IsProduction returns true if the application is running in production mode.
func (a *Application) IsProduction() bool {
	return a.config.IsProduction()
}

// IsSSLEnabled returns whether SSL/TLS should be used.
func (a *Application) IsSSLEnabled() bool {
	return a.config.IsSSLEnabled()
}

// GetSSLCertFile returns the path to the SSL certificate file.
func (a *Application) GetSSLCertFile() string {
	return a.config.GetSSLCertFile()
}

// GetSSLKeyFile returns the path to the SSL private key file.
func (a *Application) GetSSLKeyFile() string {
	return a.config.GetSSLKeyFile()
}

// Close releases external resources opened by the application, such as the
// database connection and the Redis client. It is idempotent in that it checks for nil before closing each resource.
// You must call Close in the shutdown flow to prevent resource leaks.
func (a *Application) Close() {
	if a.db != nil {
		a.db.Close()
	}
	if a.redisClient != nil {
		a.redisClient.Close()
	}
}
