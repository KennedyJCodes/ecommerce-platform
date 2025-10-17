package service_products

import (
	"log"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/output"
)

type ProductsGetService struct {
	productsRepository output.ProductsRepository
}
func NewProductsGetService(productsRepository output.ProductsRepository) input.ProductsGetService {
	return &ProductsGetService{
		productsRepository: productsRepository,
	}
}

func (p *ProductsGetService) GetProductsByBrand(brand string) ([]models.Product, error) {
	products, err := p.productsRepository.GetProductsByBrand(brand)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return products, nil
}

func (p *ProductsGetService) GetProductByID(id int) (models.Product, error) {
	product, err := p.productsRepository.GetProductByID(id)
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func (p *ProductsGetService) GetProducts() ([]models.Product, error) {
	products, err := p.productsRepository.GetProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}
