// Package output defines output port interfaces for the application.
package output

import "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models"

// ProductsRepository defines the persistence contract for product data management.
// Implementations are responsible for querying the database or external storage.
type ProductsRepository interface {
	// GetProducts fetches all products stored in the database.
	// Returns:
	//   - []models.Product: slice containing all product records.
	//   - error: non-nil if the database query fails.
	GetProducts() ([]models.Product, error)

	// GetProductByID retrieves a single product by its unique database identifier.
	// Parameters:
	//   - id: The primary key ID of the product.
	// Returns:
	//   - models.Product: the matching product record.
	//   - error: non-nil if the product is not found or a database error occurs.
	GetProductByID(id int) (models.Product, error)

	// GetProductsByBrand filters and returns products associated with a specific brand name.
	// Parameters:
	//   - brand: The exact string name of the brand to filter.
	// Returns:
	//   - []models.Product: slice of products matching the criteria.
	//   - error: non-nil if the filtered query fails.
	GetProductsByBrand(brand string) ([]models.Product, error)
}
