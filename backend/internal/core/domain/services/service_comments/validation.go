// Package service_comments provides validation logic for comment data.
package service_comments

import (
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
)

// CommentValidationData represents the payload required to validate a comment.
// Fields:
//   - Content: the text body of the comment; must be non-empty.
//   - Rating: numerical score for the comment; must be between 1 and 5 (inclusive).
type CommentValidationData struct {
	Content string
	Rating  int
}

// CommentValidator enforces business rules for comment creation.
// It implements the input.Validator interface, validating CommentValidationData.

// Validation rules:
//  1. Input must be of type CommentValidationData.
//  2. Content must not be an empty string.
//  3. Rating must be provided (non-zero).
//  4. Rating must be between 1 and 5 (inclusive).

// On validation failure, returns a ValidationError with appropriate message.
type CommentValidator struct{}

// Validate checks the provided input against comment rules.

// Parameters:
//   - input: interface{} expected to be of type CommentValidationData.

// Returns:
//   - error: nil if validation passes; ValidationError otherwise.
func (r *CommentValidator) Validate(input interface{}) error {
	// Assert correct type
	data, ok := input.(CommentValidationData)
	if !ok {
		return errors.NewValidationError("Incorrect validation data")
	}

	// Rule 1: Content must not be empty
	if data.Content == "" {
		return errors.NewValidationError("Comment content cannot be empty")
	}

	// Rule 2: Rating must be provided
	if data.Rating == 0 {
		return errors.NewValidationError("You need to enter the product rating")
	}

	// Rule 3: Rating range must be 1 to 5
	if data.Rating < 1 || data.Rating > 5 {
		return errors.NewValidationError("The rating must be between 1 to 5")
	}

	return nil
}
