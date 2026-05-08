// Package output defines persistence contracts for comments and users.
package output

import "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models"

// CommentRepository persists and retrieves comments.
type CommentRepository interface {
	// GetComments fetches all stored comments.
	// Returns:
	//   - []models.Comment: slice of comments.
	//   - error: non-nil if retrieval fails.
	GetComments() ([]models.Comment, error)

	// SaveComment stores a new comment with associated user ID and rating.
	// Parameters:
	//   - userID:  ID of the author.
	//   - content: Comment text.
	//   - rating:  Numerical rating (1–5).
	// Returns:
	//   - error: non-nil if persistence fails.
	SaveComment(userID int, content string, rating int) error
}
