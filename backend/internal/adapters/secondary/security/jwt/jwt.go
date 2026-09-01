package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models/auth"
	"github.com/golang-jwt/jwt/v5"
)

// jwtService manages operations related to JSON Web Tokens.
// It uses a secret key to sign and verify tokens.
type JWTService struct {
	secretKey []byte
}

// NewJWTService creates a new jwtService with the given secret.
// secretKey: the HMAC secret used to sign and validate tokens.
func NewJWTService(secretKey string) *JWTService {
	return &JWTService{
		secretKey: []byte(secretKey),
	}
}

// GenerateJWT generates a signed JWT for the specified userName.
// The token embeds a unique ID (jti), issued-at timestamp, username,
// and an expiration set to 15 minutes from now.
func (j *JWTService) GenerateToken(userID int, userName string, tokenType models_auth.TokenType) (string, error) {
	jtiBytes := make([]byte, 16)
	if _, err := rand.Read(jtiBytes); err != nil {
		return "", fmt.Errorf("error generating token ID: %w", err)
	}

	if tokenType == models_auth.TokenTypeAccess {
		var claims = models_auth.Claims{
		UserID:   userID,
		UserName: userName,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        hex.EncodeToString(jtiBytes),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			},
		}

		var token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString(j.secretKey)
	} else if tokenType == models_auth.TokenTypeRefresh {
		var claims = models_auth.Claims{
		UserID:   userID,
		UserName: userName,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        hex.EncodeToString(jtiBytes),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			},
		}

		var token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString(j.secretKey)
	}

	return "", fmt.Errorf("invalid token type")
}

// ValidateRefreshToken validates a refresh token string and returns its claims.
// It verifies the signature, checks that the Subject is "refresh", and ensures
// the token has not expired.
func (j *JWTService) ValidateToken(tokenString string) (*models_auth.Claims, error) {
	claims := &models_auth.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
