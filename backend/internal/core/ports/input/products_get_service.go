// Package input defines service contracts for comment-related business logic, user operations, and input validation.
package input

import "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/product"

// ProductsGetService defines the interface for retrieving product information.
// Implementations should handle filtering and fetching products from the data source.
type ProductsGetService interface {
	// GetProducts retrieves all available products in the catalog.
	// Returns:
	//   - []models_product.Product: list of all products found.
	//   - error: non-nil if the retrieval process fails.
	GetProducts() ([]models_product.Product, error)

	// GetProductByID finds a specific product using its unique identifier.
	// Parameters:
	//   - id: The unique numerical ID of the product.
	// Returns:
	//   - models_product.Product: the product details if found.
	//   - error: non-nil if the product does not exist or the query fails.
	GetProductByID(id int) (models_product.Product, error)

	// GetProductsByBrand retrieves a filtered list of products belonging to a specific brand.
	// Parameters:
	//   - brand: The name of the brand to filter by.
	// Returns:
	//   - []models_product.Product: list of products matching the brand.
	//   - error: non-nil if the query fails.
	GetProductsByBrand(brand string) ([]models_product.Product, error)
}
