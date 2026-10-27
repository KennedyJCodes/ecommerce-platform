// Package service_comments implements comment-related domain services, orchestrating validation and persistence of user comments.
package service_reviews

import (
	"context"
	dto "github.com/KennedyJCodes/ecommerce-platform/internal/core/domain/dto/review"
	models_review "github.com/KennedyJCodes/ecommerce-platform/internal/core/domain/models/review"
	"github.com/KennedyJCodes/ecommerce-platform/internal/core/ports/input"
	"github.com/KennedyJCodes/ecommerce-platform/internal/core/ports/output"
	"github.com/KennedyJCodes/ecommerce-platform/pkg/errors"
)

// ReviewAddService orchestrates the validation and persistence of a new review.
// It ensures that the review data meets business rules before saving it.

// Fields:
//   - reviewRepository: handles database operations for reviews.
//   - reviewValidate: enforces validation rules via the input.Validator interface.
type ReviewAddService struct {
	ReviewRepository output.ReviewRepository
	ReviewValidate   input.Validator
}

// NewReviewAddService constructs a ReviewAddService with the given dependencies.

// Parameters:
//   - reviewRepository: implementation of output.ReviewRepository for data access.
//   - reviewValidate: implementation of input.Validator for review data validation.

// Returns:
//   - input.ReviewAddService: service to add new reviews.
func NewReviewAddService(reviewRepository output.ReviewRepository, reviewValidate input.Validator) input.ReviewAddService {
	return &ReviewAddService{
		ReviewRepository: reviewRepository,
		ReviewValidate:   reviewValidate,
	}
}

// AddReview validates the review data and saves it to the repository.

// Steps:
//  1. Build ReviewValidationData containing content and rating.
//  2. Validate the data; return error if validation fails.
//  3. Call SaveReview on the repository; wrap errors in InternalError.

// Parameters:
//   - userID: ID of the user adding the review.
//   - content: the text content of the review.
//   - rating: numerical rating score for the review.

// Returns:
//   - error: nil on success, or a validation/InternalError on failure.
func (s *ReviewAddService) AddReview(ctx context.Context, userID int, request dto.ReviewRequest) error {
	validationData := ReviewValidationData{
		Content: request.Content,
		Rating:  request.Rating,
	}

	err := s.ReviewValidate.Validate(validationData)
	if err != nil {
		return err
	}

	newReview := models_review.Review{
		Content: request.Content,
		Rating:  request.Rating,
	}

	err = s.ReviewRepository.SaveReview(ctx, userID, newReview)
	if err != nil {
		return errors.NewInternalError("Error Saving Review").WithError(err)
	}
	return nil
}
