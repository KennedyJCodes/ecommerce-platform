package output

import "github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/models"

type CSRFRepository interface {
	Save(token *models.CSRFToken) error
	Find(value string) (*models.CSRFToken, error)
	Delete(value string) error
}
