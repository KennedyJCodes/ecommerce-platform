// Package service_comments implements comment-related domain services, orchestrating validation and retrieval of user comments.
package service_comments

import (
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
)

// CommentGetService orchestrates the retrieval of user feedback and reviews.
// It acts as a read-only domain service that abstracts the persistence layer from the delivery mechanisms (API/UI).
type CommentGetService struct {
	commentRepository output.CommentRepository
}

// NewCommentGetService constructs and returns a CommentGetService instance 
// using the provided output port for data fetching.
//
// Parameters:
//   - commentRepository: implementation of output.CommentRepository for data access.
func NewCommentGetService(commentRepository output.CommentRepository) input.CommentGetService {
    return &CommentGetService{
        commentRepository: commentRepository,
    }
}

// AllComments fetches the complete list of comments from the repository.
// It ensures that data is retrieved according to domain rules (e.g., chronological order).
func (s *CommentGetService) AllComments() ([]models.Comment, error) {
    // Call repository to get comments
    comments, err := s.commentRepository.GetComments()
    if err != nil {
        // Wrap repository errors in a domain-friendly InternalError
        return nil, errors.NewInternalError("Error while making the query")
    }
    return comments, nil
}