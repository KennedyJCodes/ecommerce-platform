// Package service_comments implements comment-related domain services, orchestrating validation and persistence of user comments.
package service_comments

import (
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
)

// CommentAddService orchestrates the validation and persistence of a new comment.
// It ensures that the comment data meets business rules before saving it.

// Fields:
//   - commentRepository: handles database operations for comments.
//   - commentValidate: enforces validation rules via the input.Validator interface.
type CommentAddService struct {
	commentRepository output.CommentRepository
    commentValidate input.Validator
}

// NewCommentAddService constructs a CommentAddService with the given dependencies.

// Parameters:
//   - commentRepository: implementation of output.CommentRepository for data access.
//   - commentValidate: implementation of input.Validator for comment data validation.

// Returns:
//   - input.CommentAddService: service to add new comments.
func NewCommentAddService(commentRepository output.CommentRepository, commentValidate input.Validator) input.CommentAddService {
    return &CommentAddService{
        commentRepository: commentRepository,
        commentValidate: commentValidate,
    }
}

// AddComment validates the comment data and saves it to the repository.

// Steps:
//  1. Build CommentValidationData containing content and rating.
//  2. Validate the data; return error if validation fails.
//  3. Call SaveComment on the repository; wrap errors in InternalError.

// Parameters:
//   - userID: ID of the user adding the comment.
//   - content: the text content of the comment.
//   - rating: numerical rating score for the comment.

// Returns:
//   - error: nil on success, or a validation/InternalError on failure.
func (s *CommentAddService) AddComment(userID int, content string, rating int) error {
    // Step 1: Prepare validation payload
    validationData := CommentValidationData{
        Content: content,
        Rating: rating,

    }

    // Step 2: Validate input
    err := s.commentValidate.Validate(validationData) 
    if err != nil{
        return err
    }

    // Step 3: Persist the comment
    err = s.commentRepository.SaveComment(userID, content, rating)
    if err != nil {
        return errors.NewInternalError("Error Saving Comment").WithError(err)
    } 
    return nil
}