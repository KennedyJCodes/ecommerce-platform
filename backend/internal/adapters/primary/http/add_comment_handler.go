// Package http provides HTTP handlers for the Watch Store API's comment endpoints.
// It defines adapters that translate HTTP requests into domain service calls and format service responses (or errors) as JSON HTTP responses.
package http

import (
	"net/http"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/adapters/primary/http/middleware"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
	httpUtil "github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/http"
)

// CommentsAddHandler handles HTTP requests for adding new comments (reviews).
// It delegates business logic to the CommentAddService and writes appropriate JSON responses based on the outcome of the service call.

// Fields:
//   - commentService: domain service used to add new comments.
type CommentsAddHandler struct {
	commentService input.CommentAddService
}

// NewCommentsAddHandler constructs and returns a new CommentsAddHandler.

// Parameters:
//   - commentService: an implementation of input.CommentAddService responsible for persisting new comments.

// Returns:
//   - *CommentsAddHandler: ready-to-use HTTP handler for the "add comment" endpoint.
func NewCommentAddsHandler(commentService input.CommentAddService) *CommentsAddHandler {
	return &CommentsAddHandler{
		commentService: commentService,
	}
}

// Handle processes HTTP requests to add a new comment.
// It expects a POST request with a JSON body matching the models.Review schema.
// The handler extracts the authenticated user ID from the request context, calls the domain service to add the comment, and returns:

//   - 200 OK with a success message on success.
//   - 400 Bad Request if the method is not POST or the JSON body is invalid.
//   - 500 Internal Server Error if context lacks user ID or adding the comment fails.

// Steps:
//  1. Verify HTTP method is POST; otherwise, return 400 with "Method Not Allowed".
//  2. Decode request body into models.Review; on JSON syntax errors, return 400.
//  3. Extract user ID from context using middleware.GetUserIDContextKey(); if missing or wrong type, return 500.
//  4. Invoke commentService.AddComment with user ID, review content, and rating; on error, return 500.
//  5. Send a JSON response with status 200 and message "Comment added".

// Parameters:
//   - w: http.ResponseWriter to write HTTP response headers and body.
//   - r: *http.Request containing HTTP request data and context.
func (h *CommentsAddHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Step 1: Method check
	if r.Method != http.MethodPost {
		httpUtil.HandleError(w, errors.NewBadRequestError(errors.ErrMethodNotAllowed))
		return
	}

	// Step 2: Decode JSON body
	var account models.Review
	if err := httpUtil.DecodeJSONBody(w, r, &account, httpUtil.MaxCommentBodySize); err != nil {
		httpUtil.HandleError(w, err)
		return
	}

	// Step 3: Extract authenticated user ID from context
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.GetUserIDContextKey())
	userIdInt, ok := userIDValue.(int)
	if !ok {
		httpUtil.HandleError(w, errors.NewInternalError(errors.ErrInternalServer))
		return
	}

	// Step 4: Call domain service to add the comment
	err := h.commentService.AddComment(userIdInt, account.Content, account.Rating)
	if err != nil {
		httpUtil.HandleError(w, errors.NewInternalError(errors.ErrInternalServer))
		return
	}

	// Step 5: Send success response
	httpUtil.SendJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Comment added",
	})
}
