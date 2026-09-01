// Package output defines persistence contracts for comments and users.
package output

import (
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/review"
)

// ReviewRepository persists and retrieves reviews.
type ReviewRepository interface {
	// GetReviews fetches all stored reviews.
	// Returns:
	//   - []models_review.Review: slice of reviews.
	//   - error: non-nil if retrieval fails.
	GetReviews() ([]models_review.Review, error)

	// SaveReview stores a new review with associated user ID and rating.
	// Parameters:
	//   - userID:  ID of the author.
	//   - review:  The review object to be saved.
	// Returns:
	//   - error: non-nil if persistence fails.
	SaveReview(userID int, review models_review.Review) error
}
