// Package bootstrap provides high-level factory functions to initialize and wire the application's infrastructure components.
// It serves as an intermediary layer that connects secondary adapters (MySQL, Static Files) with the core ports defined in the domain.
package bootstrap

import (
	repository_mysql "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/secondary/repository/mysql"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/secondary/static"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/config"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/security/security_auth"
	"github.com/jmoiron/sqlx"
)

// SetupStaticFileAdapter initializes the static file service.
// It retrieves the static directory path from the application configuration and
// returns an implementation of the StaticFilePort.

// Parameters:
//   - appConfig: the global application configuration manager.
//
// Returns:
//   - output.StaticFilePort: an adapter capable of serving and validating static assets.
func SetupStaticFileAdapter(appConfig *config.AppConfig) output.StaticFilePort {
	staticDir := appConfig.GetStaticDir()
	return static.NewStaticFileAdapter(staticDir)
}

// SetupUserRepository initializes the user repository with its necessary security dependencies.
// It explicitly injects a BcryptHasher into the SQLUserRepository, ensuring that all user persistence operations follow the defined security standards for password hashing.

// Parameters:
//   - db: an active *sqlx.DB connection pool.
//
// Returns:
//   - output.UserRepository: a repository ready to handle user-related database operations.
func SetupUserRepository(db *sqlx.DB) output.UserRepository {
	// Define the hashing strategy to be used by the repository.
	hasher := security_auth.BcryptHasher{}
	return repository_mysql.NewSQLUserRepository(db, hasher)
}

// SetupTokenService creates a new JWTService instance with the secret key from config.
func SetupTokenService(appConfig *config.AppConfig) *security_auth.JWTService {
	return security_auth.NewJWTService(appConfig.GetJWTSecret())
}

// SetupCommonServices configures global shared services that do not require per-instance state, such as the JWT authentication service.
// It sets the default secret key used for signing and validating authentication tokens across the entire application.

// Parameters:
//   - appConfig: the global application configuration manager.
func SetupCommonServices(appConfig *config.AppConfig) {
	// Initialize the global JWT service with the secret key from config.
	security_auth.SetDefaultJWTService(appConfig.GetJWTSecret())
}
