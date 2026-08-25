// Package models defines core domain entities for the sale‑watches application.
// This file declares the Review type, which is used to process and validate user feedback submitted via the application's public API.
package models

// Review represents a formal user evaluation of a product or service.
// Unlike internal comment structures, this model acts as a Data Transfer Object (DTO) specifically designed to capture the core sentiment (content and rating provided by the customer during the submission process.
type Review struct {
	// Content: the textual body of the review containing the user's opinion.
	Content  string `json:"content"`

	// Rating: the numeric score assigned by the user (typically on a scale of 1–5).
	// This value is used to calculate product satisfaction metrics.
	Rating   int    `json:"rating"`	
}