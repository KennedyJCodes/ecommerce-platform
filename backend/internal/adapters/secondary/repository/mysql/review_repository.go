// Package repository provides SQL-based implementations of output ports for data persistence.
// This file contains SqlReviewRepository, which implements ReviewRepository using MySQL.
package repository_mysql

import (
	"context"

	"github.com/KennedyJCodes/ecommerce-platform/internal/core/domain/models/review"
	"github.com/KennedyJCodes/ecommerce-platform/internal/core/ports/output"
	appErrors "github.com/KennedyJCodes/ecommerce-platform/pkg/errors"
	"github.com/jmoiron/sqlx"
)

type SqlReviewRepository struct {
	db *sqlx.DB
}

func NewSqlReviewRepository(db *sqlx.DB) (output.ReviewRepository, error) {
	if db == nil {
		return nil, appErrors.NewInternalError(appErrors.ErrDatabaseConnection)
	}
	return &SqlReviewRepository{db: db}, nil
}

func (r *SqlReviewRepository) GetReviews(ctx context.Context) ([]models_review.Review, error) {
	var reviews []models_review.Review
	const query = `
	SELECT c.id, c.date, c.content, c.user_id, u.username AS username, c.rating
	FROM comments c
	JOIN user_registration u ON c.user_id = u.user_id
	ORDER BY c.date DESC`

	if err := r.db.SelectContext(ctx, &reviews, query); err != nil {
		return nil, appErrors.NewInternalError(appErrors.ErrDatabaseQuery).WithError(err)
	}
	return reviews, nil
}

func (r *SqlReviewRepository) SaveReview(ctx context.Context, userID int, review models_review.Review) error {
	const query = `INSERT INTO comments (user_id, content, rating, date) VALUES (?, ?, ?, NOW())`
	if _, err := r.db.ExecContext(ctx, query, userID, review.Content, review.Rating); err != nil {
		return appErrors.NewInternalError(appErrors.ErrDatabaseInsert).WithError(err)
	}
	return nil
}
