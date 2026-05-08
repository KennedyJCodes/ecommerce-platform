// Package service_products implements the domain logic for watch catalog retrieval.
package service_products

import (
	"log"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
)

// ProductsGetService orchestrates the retrieval of product information.
// It serves as a read-only gateway that decouples the delivery layer from the persistence details, ensuring that catalog queries follow domain standards.
type ProductsGetService struct {
	productsRepository output.ProductsRepository
}

// NewProductsGetService constructs a ProductsGetService with its required repository dependency, satisfying the input.ProductsGetService interface.
func NewProductsGetService(productsRepository output.ProductsRepository) input.ProductsGetService {
	return &ProductsGetService{
		productsRepository: productsRepository,
	}
}

// GetProductsByBrand filters and retrieves watches belonging to a specific brand.
// Parameters:
//   - brand: the manufacturer's name to filter by.
//
// Returns a slice of products or an error if the repository query fails.
func (p *ProductsGetService) GetProductsByBrand(brand string) ([]models.Product, error) {
	products, err := p.productsRepository.GetProductsByBrand(brand)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return products, nil
}

// GetProductByID fetches the details of a single watch based on its unique ID.
// Parameters:
//   - id: the unique numerical identifier of the product.
//
// Returns the product model or an empty model and an error if not found.
func (p *ProductsGetService) GetProductByID(id int) (models.Product, error) {
	product, err := p.productsRepository.GetProductByID(id)
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}

// GetProducts retrieves the entire collection of watches available in the store.
// Returns all product entries or an error in case of persistence failure.
func (p *ProductsGetService) GetProducts() ([]models.Product, error) {
	products, err := p.productsRepository.GetProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}
