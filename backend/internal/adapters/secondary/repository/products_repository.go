package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/output"
	Custom_errors "github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
	"github.com/jmoiron/sqlx"
)

type SqlProductsRepository struct {
	db *sqlx.DB
}

func NewSqlProductsRepository(db *sqlx.DB) output.ProductsRepository {
	if db == nil {
		log.Fatal(Custom_errors.NewInternalError(Custom_errors.ErrDatabaseConnection).Error())
	}

	return &SqlProductsRepository{
		db: db,
	}
}

func (r *SqlProductsRepository) GetProducts() ([]models.Product, error) {
	var product []models.Product
	const sqlQuery = `
	SELECT Product_ID, Product_Name, Product_Description, Product_Price, Stock_Quantity, Brand, Image_URL, Movement_Type
	FROM Products `

	// Execute the query and scan results into comments slice.
	err := r.db.Select(&product, sqlQuery)
	if err != nil {
		// Wrap low-level DB error in a domain-friendly InternalError.
		return nil, Custom_errors.NewInternalError(Custom_errors.ErrDatabaseQuery).WithError(err)
	}

	return product, nil
}

func (r *SqlProductsRepository) GetProductByID(id int) (models.Product, error) {
	var product models.Product
	const sqlQuery = `
	SELECT Product_ID, Product_Name, Product_Description, Product_Price, Stock_Quantity, Brand, Image_URL, Movement_Type
	FROM Products
	WHERE Product_ID = ?`

	if id < 0 {
		return models.Product{}, Custom_errors.NewBadRequestError("Invalid ID format")
	}

	err := r.db.Get(&product, sqlQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Product{}, Custom_errors.NewNotFoundError("product not found")
		}
		return models.Product{}, Custom_errors.NewInternalError(Custom_errors.ErrDatabaseQuery).WithError(err)
	}

	return product, nil
}

func (r *SqlProductsRepository) GetProductsByBrand(brand string) ([]models.Product, error) {
	var productByBrand []models.Product
	const sqlQuery = `
		SELECT Product_ID, Product_Name, Product_Description, Product_Price, Stock_Quantity, Brand, Image_URL, Movement_Type
		FROM Products
		WHERE Brand = ?`
	err := r.db.Select(&productByBrand, sqlQuery, brand)
	if err != nil {
		return []models.Product{}, Custom_errors.NewInternalError(Custom_errors.ErrDatabaseQuery).WithError(err)
	}

	return productByBrand, nil
}
