// Package input defines service contracts for comment-related business logic, user operations, and input validation.
package input

import "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models"

// CommentGetService handles retrieval of comments.
type CommentGetService interface {
	// AllComments returns all comments ordered by date descending.
	// Returns:
	//   - []models.Comment: list of comments including metadata.
	//   - error: non-nil if the query fails.
	AllComments() ([]models.Comment, error)
}
