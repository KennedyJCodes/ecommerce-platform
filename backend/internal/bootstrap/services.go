// Package bootstrap provides high-level factory functions to initialize and wire the application's infrastructure components.
// This file specifically orchestrates the creation of domain services by injecting required repositories, validators, and cross-cutting services like CSRF.
package bootstrap

import (
	repository_mysql "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/secondary/repository/mysql"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/services/service_auth"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/services/service_comments"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/services/service_products"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/jmoiron/sqlx"
)

// SetupCommentService initializes the comment-related services.
// It creates a single SQL repository and shares it between the retrieval (Get) and creation (Add) services. It also injects a CommentValidator to ensure business rules are met before persistence.

// Parameters:
//   - db: an active *sqlx.DB connection pool.
//
// Returns:
//   - input.CommentGetService: service for fetching comments.
//   - input.CommentAddService: service for adding new comments.
func SetupCommentService(db *sqlx.DB) (input.CommentGetService, input.CommentAddService) {
	commentRepo := repository_mysql.NewSqlCommentRepository(db)
	commentValidator := &service_comments.CommentValidator{}
	return service_comments.NewCommentGetService(commentRepo), service_comments.NewCommentAddService(commentRepo, commentValidator)
}

// SetupProductsService initializes the product catalog service.

// Parameters:
//   - db: an active *sqlx.DB connection pool.
//
// Returns:
//   - input.ProductsGetService: the service implementing the product retrieval port.
func SetupProductsService(db *sqlx.DB) input.ProductsGetService {
	productsRepo := repository_mysql.NewSqlProductsRepository(db)
	return service_products.NewProductsGetService(productsRepo)
}

// SetupUserService initializes the authentication and registration services.
// This function wires the user repository with the necessary validators for usernames and passwords. It also injects the CSRFService to handle security token generation during the authentication lifecycle.

// Parameters:
//   - userRepo: the user persistence adapter (output port).
//   - csrfService: the service for managing CSRF tokens (input port).
//
// Returns:
//   - input.UserServiceLogin: the service handling user authentication.
//   - input.UserServiceRegister: the service handling new user creation.
func SetupUserService(userRepo output.UserRepository, csrfService input.CSRFService) (input.UserServiceLogin, input.UserServiceRegister) {
	// Initialize specific domain validators.
	userNameValidator := &service_auth.UserNameValidator{}
	passwordValidator := &service_auth.PasswordValidator{}

	// Create services with shared dependencies.
	// The 'nil' parameter is reserved for future extensions (e.g., specific cookie setters).
	return service_auth.NewUserLoginService(userRepo, userNameValidator, passwordValidator, csrfService, nil), service_auth.NewUserRegisterService(userRepo, userNameValidator, passwordValidator, csrfService, nil)
}
