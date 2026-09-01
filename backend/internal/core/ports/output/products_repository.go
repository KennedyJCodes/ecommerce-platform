// Package output defines output port interfaces for the application.
package output

import "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/product"

// ProductsRepository defines the persistence contract for product data management.
// Implementations are responsible for querying the database or external storage.
type ProductsRepository interface {
	// GetProducts fetches all products stored in the database.
	// Returns:
	//   - []models_product.Product: slice containing all product records.
	//   - error: non-nil if the database query fails.
	GetProducts() ([]models_product.Product, error)

	// GetProductByID retrieves a single product by its unique database identifier.
	// Parameters:
	//   - id: The primary key ID of the product.
	// Returns:
	//   - models_product.Product: the matching product record.
	//   - error: non-nil if the product is not found or a database error occurs.
	GetProductByID(id int) (models_product.Product, error)

	// GetProductsByBrand filters and returns products associated with a specific brand name.
	// Parameters:
	//   - brand: The exact string name of the brand to filter.
	// Returns:
	//   - []models_product.Product: slice of products matching the criteria.
	//   - error: non-nil if the filtered query fails.
	GetProductsByBrand(brand string) ([]models_product.Product, error)
}
