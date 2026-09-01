// Package models defines core domain entities for the sale-watches application.
// This file declares the TokenPair type used for transporting authentication tokens
// between service and handler layers.
package models_auth

// TokenPair holds the access and refresh tokens issued during user authentication.
// AccessToken is a short-lived JWT (15 minutes) used for authorizing API requests.
// RefreshToken is a long-lived JWT (7 days) used exclusively to obtain a new
// token pair without requiring the user to re-enter credentials.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
