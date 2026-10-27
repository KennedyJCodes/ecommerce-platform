// Package service_reviews implements review-related domain services, orchestrating validation and retrieval of user reviews.
package service_reviews

import (
	"context"
	"github.com/KennedyJCodes/ecommerce-platform/internal/core/domain/models/review"
	"github.com/KennedyJCodes/ecommerce-platform/internal/core/ports/input"
	"github.com/KennedyJCodes/ecommerce-platform/internal/core/ports/output"
	"github.com/KennedyJCodes/ecommerce-platform/pkg/errors"
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
func (s *ReviewGetService) AllReviews(ctx context.Context) ([]models_review.Review, error) {
	// Call repository to get reviews
	comments, err := s.reviewRepository.GetReviews(ctx)
	if err != nil {
		return nil, errors.NewInternalError(errors.ErrDatabaseQuery).WithError(err)
	}
	return comments, nil
}
