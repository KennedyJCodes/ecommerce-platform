// Package input defines service contracts for comment-related business logic, user operations, and input validation.
package input

// CommentAddService handles creation of new comments.
type CommentAddService interface {
	// AddComment creates a new comment entry.
    // Parameters:
    //   - userID:   ID of the user adding the comment.
    //   - content:  Body text of the comment.
    //   - rating:   Numerical rating (1–5).
    // Returns:
    //   - error: non-nil if validation or persistence fails.
	AddComment(userID int, content string, rating int) error
}
