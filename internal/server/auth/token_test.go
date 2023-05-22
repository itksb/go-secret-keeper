package auth

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJwtTokenProvider_ValidateToken(t *testing.T) {
	tokenKey := []byte("secretKey")
	tokenProvider := NewJwtTokenProvider(tokenKey, func() time.Time {
		return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	})

	// Generate a valid token
	validToken, err := tokenProvider.GenerateToken(context.Background(), "12345")
	if err != nil {
		t.Errorf("Error generating valid token: %v", err)
	}

	expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDk1NDU2MDAsInVzZXJJRCI6IjEyMzQ1In0.dLEkI9cLkN4Ai7dgjcWsBZxk1VQEEW6lvxvtrlvY6XI"
	require.Equal(t, expectedToken, validToken, "The two tokens should be the same.")

	// Test validation with valid token
	userID, err := tokenProvider.ValidateToken(context.Background(), validToken)
	if err != nil {
		t.Errorf("Valid token validation failed: %v", err)
	}
	if userID != "12345" {
		t.Errorf("Valid token validation returned incorrect userID. Expected: 12345, Got: %s", userID)
	}

}
