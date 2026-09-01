// Package models defines core domain entities for the sale‑watches application.

// This file declares the Comment type, representing user feedback with rating.
package models_review

// Comment represents a user’s feedback on a product or service.

// Fields:
//   - ID:        unique identifier of the comment.
//   - Date:      timestamp when the comment was posted, usually in ISO 8601 format.
//   - UserName:  identifier of the user who posted the comment.
//   - Content:   textual body of the comment.
//   - Rating:    numeric score given by the user (e.g., 1–5).
type Review struct {
	ID int `db:"id"`
	Date string `db:"date"`
	UserID int `db:"user_id"`
	UserName string `db:"user_name"`
	Content string `db:"content"`
	Rating int `db:"rating"`
} 