// Package input defines service contracts for comment-related business logic, user operations, and input validation.
package input

import "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/review"

// ReviewGetService handles retrieval of reviews.
type ReviewGetService interface {
	// AllReviews returns all reviews ordered by date descending.
	// Returns:
	//   - []models_review.Review: list of reviews including metadata.
	//   - error: non-nil if the query fails.
	AllReviews() ([]models_review.Review, error)
}
