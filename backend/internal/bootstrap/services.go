// Package bootstrap provides high-level factory functions to initialize and wire the application's infrastructure components.
// This file specifically orchestrates the creation of domain services by injecting required repositories, validators, and cross-cutting services like CSRF.
package bootstrap

import (
	"fmt"

	repository_mysql "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/secondary/repository/mysql"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/services/service_auth"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/services/service_reviews"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/services/service_products"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/security/security_auth"
	"github.com/jmoiron/sqlx"
)

// SetupReviewService initializes the review-related services.
// It creates a single SQL repository and shares it between the retrieval (Get) and creation (Add) services. It also injects a ReviewValidator to ensure business rules are met before persistence.

// Parameters:
//   - db: an active *sqlx.DB connection pool.
//
// Returns:
//   - input.ReviewGetService: service for fetching reviews.
//   - input.ReviewAddService: service for adding new reviews.
func SetupReviewService(db *sqlx.DB) (input.ReviewGetService, input.ReviewAddService, error) {
	reviewRepo, err := repository_mysql.NewSqlReviewRepository(db)
	if err != nil {
		return nil, nil, fmt.Errorf("setup review repository: %w", err)
	}

	reviewValidator := &service_reviews.ReviewValidator{}
	return service_reviews.NewReviewGetService(reviewRepo), service_reviews.NewReviewAddService(reviewRepo, reviewValidator), nil
}

// SetupProductsService initializes the product catalog service.

// Parameters:
//   - db: an active *sqlx.DB connection pool.
//
// Returns:
//   - input.ProductsGetService: the service implementing the product retrieval port.
func SetupProductsService(db *sqlx.DB) (input.ProductsGetService, error) {
	productsRepo, err := repository_mysql.NewSqlProductsRepository(db)
	if err != nil {
		return nil, fmt.Errorf("setup products repository: %w", err)
	}

	return service_products.NewProductsGetService(productsRepo), nil
}

// SetupUserService initializes the authentication and registration services.
// This function wires the user repository with the necessary validators for usernames and passwords. It also injects the TokenService and CSRFService to handle security token generation during the authentication lifecycle.

// Parameters:
//   - userRepo: the user persistence adapter (output port).
//   - tokenService: the service for JWT generation and validation (input port).
//   - csrfService: the service for managing CSRF tokens (input port).
//   - hasher: the service for hashing passwords (input port).
//
// Returns:
//   - input.UserServiceLogin: the service handling user authentication.
//   - input.UserServiceRegister: the service handling new user creation.
func SetupUserService(userRepo output.UserRepository, tokenService output.TokenService, csrfService output.CSRFService) (input.UserServiceLogin, input.UserServiceRegister) {
	// Initialize specific domain validators.
	userNameValidator := &service_auth.UserNameValidator{}
	passwordValidator := &service_auth.PasswordValidator{}
	emailValidator := &service_auth.EmailValidator{}
	passwordHasher := &security_auth.BcryptHasher{}

	return service_auth.NewUserLoginService(userRepo, userNameValidator, passwordValidator, tokenService, csrfService), service_auth.NewUserRegisterService(userRepo, userNameValidator, passwordValidator, emailValidator, tokenService, csrfService, passwordHasher)
}
