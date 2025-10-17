// Package service_comments implements comment-related domain services, orchestrating validation and retrieval of user comments.
package service_comments

import (
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
)

// CommentGetService handles the retrieval of comments from the repository and can apply additional business rules or transformations if needed.

// Fields:
//   - commentRepository: provides access to persisted comment data.
//   - commentValidate: validator for input parameters (unused currently, reserved for
//     potential future filters or pagination validations).
type CommentGetService struct {
	commentRepository output.CommentRepository
}

// NewCommentGetService constructs and returns a CommentGetService instance.

// Parameters:
//   - commentRepository: implementation of output.CommentRepository for data fetching.
//   - commentValidate: implementation of input.Validator for any retrieval constraints.

// Returns:
//   - input.CommentGetService: service interface for fetching all comments.
func NewCommentGetService(commentRepository output.CommentRepository) input.CommentGetService {
    return &CommentGetService{
        commentRepository: commentRepository,
    }
}

// AllComments retrieves all comments sorted by date (descending) via the repository.
// It returns an InternalError if the underlying query fails.
//
// Returns:
//   - []models.Comment: slice of Comment models including ID, Date, Content, UserID, UserName, and Rating.
//   - error: non-nil if database retrieval fails.
func (s *CommentGetService) AllComments() ([]models.Comment, error) {
    // Call repository to get comments
    comments, err := s.commentRepository.GetComments()
    if err != nil {
        // Wrap repository errors in a domain-friendly InternalError
        return nil, errors.NewInternalError("Error while making the query")
    }
    return comments, nil
}