// Package http implements HTTP handlers for the ecommerce-platform application.
// This file, which contains the CommentsHandler, is planned to be further enhanced in a future development phase. Note that only this file's functionality will be completed later.
package public_handlers

import (
	"net/http"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
	httpUtil "github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/http"
)

// CommentsHandler handles HTTP requests related to comments.

// It acts as an adapter between HTTP requests and the business logic provided by the CommentService interface defined in the core domain. This handler currently supports retrieving comments.
type ReviewsGetHandler struct {
	ReviewService input.ReviewGetService
}

// NewReviewsGetHandler creates and returns a new instance of ReviewsGetHandler.

// It receives an implementation of the ReviewGetService interface, which contains the business logic for managing reviews.
func NewReviewsGetHandler(reviewService input.ReviewGetService) *ReviewsGetHandler {
	return &ReviewsGetHandler{
		ReviewService: reviewService,
	}
}

// Handle processes incoming HTTP requests to retrieve reviews.

// It calls the GetReviews method of the commentService to fetch reviews. If an error occurs during the retrieval, it sends an HTTP error response with a 500 (Internal Server Error) status using a utility function. If successful, it returns the reviews in JSON format with an HTTP 200 (OK) status.
func (h *ReviewsGetHandler) Handle(w http.ResponseWriter, r *http.Request) {
	comments, err := h.ReviewService.AllReviews()
	if err != nil {
		httpUtil.HandleError(w, errors.NewInternalError("Error getting feedback"))
		return
	}

	httpUtil.SendJSONResponse(w, http.StatusOK, comments)
}
