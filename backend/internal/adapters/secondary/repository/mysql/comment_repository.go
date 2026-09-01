// Package repository provides SQL-based implementations of output ports for data persistence.
// This file contains SqlCommentRepository, which implements CommentRepository using a MySQL database via sqlx.
package repository_mysql

import (
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/review"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/ecommerce-platform/pkg/errors"
	"github.com/jmoiron/sqlx"
)

// SqlReviewRepository implements output.ReviewRepository using a SQL database.
// It uses sqlx for database interactions and expects a valid *sqlx.DB connection.
//
// Fields:
//   - db: *sqlx.DB instance for executing queries.
type SqlReviewRepository struct {
	db *sqlx.DB
}

// NewSqlReviewRepository creates a new SqlReviewRepository.
// It validates the required database dependency and reports configuration
// errors to the caller instead of exiting the process.

// Parameters:
//   - db: *sqlx.DB connection to the reviews database.

// Returns:
//   - output.ReviewRepository: initialized repository instance.
func NewSqlReviewRepository(db *sqlx.DB) (output.ReviewRepository, error) {
	if db == nil {
		return nil, errors.NewInternalError(errors.ErrDatabaseConnection)
	}

	return &SqlReviewRepository{
		db: db,
	}, nil
}

// GetReviews retrieves all reviews from the database, ordered by date descending.
// It performs a JOIN with the user_registration table to include the commenter's username.

// Returns:
//   - []models_review.Review: slice of Review models containing ID, Date, Content, UserID, UserName, and Rating.
//   - error: non-nil if the query fails, wrapped as an InternalError.
func (r *SqlReviewRepository) GetReviews() ([]models_review.Review, error) {
	var review []models_review.Review
	// Define SQL query to select reviews and join with user table.
	const sqlQuery = `
	SELECT 
		c.id,
		c.date,
		c.content,
		c.user_id,
		u.username AS username,
		c.rating
	FROM comments c
	JOIN user_registration u
		ON c.user_id = u.user_id
		ORDER BY c.date DESC
	`

	// Execute the query and scan results into reviews slice.
	err := r.db.Select(&review, sqlQuery)
	if err != nil {
		// Wrap low-level DB error in a domain-friendly InternalError.
		return nil, errors.NewInternalError(errors.ErrDatabaseQuery).WithError(err)
	}
	return review, nil
}

// SaveReview inserts a new review into the database with the current timestamp.
// It uses parameterized queries to prevent SQL injection.

// Parameters:
//   - userID: ID of the authenticated user adding the review.
//   - content: text content of the review.
//   - rating: numerical rating associated with the review.

// Returns:
//   - error: non-nil if the insert fails, wrapped as an InternalError.
func (r *SqlReviewRepository) SaveReview(userID int, review models_review.Review) error {
	const query = `INSERT INTO comments (user_id, content, rating, date)
	VALUES (?, ?, ?, NOW())`

	// Execute the insert query with provided parameters.
	_, err := r.db.Exec(query, userID, review.Content, review.Rating)
	if err != nil {
		// Return a generic InternalError on failure.
		return errors.NewInternalError("Error querying the database")
	}
	return nil
}
