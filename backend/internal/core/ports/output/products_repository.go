package output

import "github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/models"

type ProductsRepository interface {
	GetProducts() ([]models.Product, error)
	GetProductByID(id int) (models.Product, error)
}
