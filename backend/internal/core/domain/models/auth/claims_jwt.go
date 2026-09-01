// Package models defines core domain entities for the sale‑watches application.

// This file declares the Claims type used in JSON Web Tokens for user authentication.
package models_auth

import "github.com/golang-jwt/jwt/v5"

// Claims represents the JWT payload for authenticated users.

// It embeds jwt.RegisteredClaims—which includes standard fields like ExpiresAt (exp), Issuer (iss), Subject (sub), NotBefore (nbf), IssuedAt (iat), Audience (aud), and ID (jti)—and adds a custom UserName claim for identifying the user. This structure conforms to RFC 7519 and integrates seamlessly with the golang‑jwt library.
type Claims struct {
	UserID int `json:"user_id"` // Custom claim for user id
	UserName string `json:"username"` // Custom claim for the user's username
	Type string `json:"type"` // Custom claim for token type (e.g., "access" or "refresh")
	jwt.RegisteredClaims // Standard JWT claims
}
