// Package securityAuth provides JWT creation and validation services for the sale‑watches application’s authentication flow.
package security_auth

import (
	"fmt"
	"time"

	"github.com/David-Alejandro-Jimenez/ecommerce-platform/internal/core/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

// JWTService manages operations related to JSON Web Tokens.
// It uses a secret key to sign and verify tokens.
type JWTService struct {
	secretKey []byte
}

// NewJWTService creates a new JWTService with the given secret.
// secretKey: the HMAC secret used to sign and validate tokens.
func NewJWTService(secretKey string) *JWTService {
	return &JWTService{
		secretKey: []byte(secretKey),
	}
}

// GenerateJWT generates a signed JWT for the specified userName.
// The token embeds the username and an expiration set to one hour from now.
func (j *JWTService) GenerateJWT(userId int, userName string) (string, error) {
	var claims = models.Claims{
		UserID:   userId,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Hour)),
		},
	}

	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// defaultJWTService holds the globally configured JWTService for convenience.
var defaultJWTService *JWTService

// SetDefaultJWTService configures the package‑level JWTService.
// It should be called once at application startup with the secret key.
func SetDefaultJWTService(secretKey string) {
	defaultJWTService = NewJWTService(secretKey)
}

// GenerateJWT signs a token for userName using the default service.
// Returns an error if the service has not been initialized.
func GenerateJWT(userId int, userName string) (string, error) {
	if defaultJWTService == nil {
		return "", fmt.Errorf("JWT service not initialized")
	}
	return defaultJWTService.GenerateJWT(userId, userName)
}

func ParseTokenWithClaims(tokenString string) (*models.Claims, error) {
	if defaultJWTService == nil {
		return nil, fmt.Errorf("JWT service not initialized")
	}

	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return defaultJWTService.secretKey, nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
