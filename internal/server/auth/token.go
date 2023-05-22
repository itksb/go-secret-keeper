package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"time"
)

var _ contract.ITokenProvider = &JwtTokenProvider{}

// JwtTokenProvider - JWT Token Provider
type JwtTokenProvider struct {
	tokenKey []byte
	nowFunc  func() time.Time
}

// NewJwtTokenProvider - Create new JWT Token Provider
func NewJwtTokenProvider(tokenKey []byte, f func() time.Time) *JwtTokenProvider {
	return &JwtTokenProvider{tokenKey: tokenKey, nowFunc: f}
}

// GenerateToken - Generate new token
func (j *JwtTokenProvider) GenerateToken(
	ctx context.Context,
	userID string,
) (string, error) {
	// Create a new token with the standard claims
	token := jwt.New(jwt.SigningMethodHS256)
	// Set the claims for the token
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userID
	claims["exp"] = j.nowFunc().Add(time.Hour * 24).Unix()
	// Sign the token with a secret key
	// Replace "your-secret-key" with your actual secret key
	tokenString, err := token.SignedString([]byte(j.tokenKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

// ValidateToken - Validate token
func (j *JwtTokenProvider) ValidateToken(
	ctx context.Context,
	tokenString string,
) (string, error) {
	// Parse the token string
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return j.tokenKey, nil
	})

	if err != nil {
		return "", err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["userID"].(string); ok {
			return userID, nil
		}
	}

	return "", fmt.Errorf("invalid token")
}
