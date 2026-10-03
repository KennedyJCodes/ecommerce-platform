// Package service_reviews implements review-related domain services, orchestrating validation and retrieval of user reviews.
package service_reviews

import (
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/review"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
)

// ReviewGetService orchestrates the retrieval of user feedback and reviews.
// It acts as a read-only domain service that abstracts the persistence layer from the delivery mechanisms (API/UI).
type ReviewGetService struct {
	reviewRepository output.ReviewRepository
}

// NewReviewGetService constructs and returns a ReviewGetService instance
// using the provided output port for data fetching.
//
// Parameters:
//   - reviewRepository: implementation of output.ReviewRepository for data access.
func NewReviewGetService(reviewRepository output.ReviewRepository) input.ReviewGetService {
	return &ReviewGetService{
		reviewRepository: reviewRepository,
	}
}

// AllReviews fetches the complete list of reviews from the repository.
// It ensures that data is retrieved according to domain rules (e.g., chronological order).
func (s *ReviewGetService) AllReviews() ([]models_review.Review, error) {
	// Call repository to get reviews
	comments, err := s.reviewRepository.GetReviews()
	if err != nil {
		// Wrap repository errors in a domain-friendly InternalError
		return nil, errors.NewInternalError("Error while making the query")
	}
	return comments, nil
}
