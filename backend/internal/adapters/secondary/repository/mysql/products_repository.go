// Package repository_mysql provides SQL-based implementations of output ports for data persistence.
// This file contains SqlProductsRepository, which implements the ProductsRepository port using a MySQL database to manage product data.
package repository_mysql

import (
	"database/sql"
	"errors"
	"log"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/output"
	Custom_errors "github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
	"github.com/jmoiron/sqlx"
)

// SqlProductsRepository implements output.ProductsRepository using a SQL database.
// It provides methods to query the product catalog using sqlx for mapping database rows to domain models.
type SqlProductsRepository struct {
	// db: *sqlx.DB instance for executing read operations on the products table.
	db *sqlx.DB
}

// NewSqlProductsRepository creates and returns a new instance of SqlProductsRepository.
// It ensures the database connection is valid; otherwise, it triggers a fatal log.
func NewSqlProductsRepository(db *sqlx.DB) output.ProductsRepository {
	if db == nil {
		log.Fatal(Custom_errors.NewInternalError(Custom_errors.ErrDatabaseConnection).Error())
	}

	return &SqlProductsRepository{
		db: db,
	}
}

// GetProducts retrieves the complete list of products from the database.
// Returns:
//   - []models.Product: a slice containing all products in the catalog.
//   - error: an InternalError if the query execution or scanning fails.
func (r *SqlProductsRepository) GetProducts() ([]models.Product, error) {
	var product []models.Product
	const sqlQuery = `
	SELECT Product_ID, Product_Name, Product_Description, Product_Price, Stock_Quantity, Brand, Image_URL, Movement_Type
	FROM Products `

	// Execute the query and map the entire result set to the product slice.
	err := r.db.Select(&product, sqlQuery)
	if err != nil {
		// Wrap low-level DB error in a domain-friendly InternalError.
		return nil, Custom_errors.NewInternalError(Custom_errors.ErrDatabaseQuery).WithError(err)
	}

	return product, nil
}

// GetProductByID retrieves a single product by its unique identifier.
// It performs a parameterized query to prevent SQL injection and checks for the existence of the record.
// Parameters:
//   - id: the integer ID of the product to retrieve.
// Returns:
//   - models.Product: the found product model.
//   - error: a NotFoundError if the ID does not exist, or a BadRequestError for invalid input formats.
func (r *SqlProductsRepository) GetProductByID(id int) (models.Product, error) {
	var product models.Product
	const sqlQuery = `
	SELECT Product_ID, Product_Name, Product_Description, Product_Price, Stock_Quantity, Brand, Image_URL, Movement_Type
	FROM Products
	WHERE Product_ID = ?`

	// Validate input ID before querying.
	if id < 0 {
		return models.Product{}, Custom_errors.NewBadRequestError("Invalid ID format")
	}

	// Use Get for a single row result.
	err := r.db.Get(&product, sqlQuery, id)
	if err != nil {
		// Specific check for no rows to return a domain-friendly 404 error.
		if errors.Is(err, sql.ErrNoRows) {
			return models.Product{}, Custom_errors.NewNotFoundError("product not found")
		}
		return models.Product{}, Custom_errors.NewInternalError(Custom_errors.ErrDatabaseQuery).WithError(err)
	}

	return product, nil
}

// GetProductsByBrand retrieves all products that belong to a specific brand.
// Parameters:
//   - brand: string representing the brand name (e.g., "Rolex", "Casio").
// Returns:
//   - []models.Product: a slice of products matching the criteria.
func (r *SqlProductsRepository) GetProductsByBrand(brand string) ([]models.Product, error) {
	var productByBrand []models.Product
	const sqlQuery = `
		SELECT Product_ID, Product_Name, Product_Description, Product_Price, Stock_Quantity, Brand, Image_URL, Movement_Type
		FROM Products
		WHERE Brand = ?`

	// Execute parameterized query to filter by brand.	
	err := r.db.Select(&productByBrand, sqlQuery, brand)
	if err != nil {
		return []models.Product{}, Custom_errors.NewInternalError(Custom_errors.ErrDatabaseQuery).WithError(err)
	}

	return productByBrand, nil
}
