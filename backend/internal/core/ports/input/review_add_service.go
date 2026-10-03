// Package input defines service contracts for comment-related business logic, user operations, and input validation.
package input

import "github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/dto/review"

// ReviewAddService handles creation of new reviews.
type ReviewAddService interface {
	// AddReview creates a new review entry.
    // Parameters:
    //   - userID:   ID of the user adding the review.
    //   - content:  Body text of the review.
    //   - rating:   Numerical rating (1–5).
    // Returns:
    //   - error: non-nil if validation or persistence fails.
	AddReview(userID int, request dto.ReviewRequest) error
}
