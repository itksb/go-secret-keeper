package auth

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// TestJwtTokenProvider_ValidateToken - Test for ValidateToken
func TestJwtTokenProvider_ValidateToken(t *testing.T) {
	tokenKey := []byte("secretKey")
	tokenProvider := NewJwtTokenProvider(tokenKey, func() time.Time {
		return time.Date(2035, 1, 1, 0, 0, 0, 0, time.UTC)
	})

	// Generate a valid token
	validToken, err := tokenProvider.GenerateToken(context.Background(), "12345")
	if err != nil {
		t.Errorf("Error generating valid token: %v", err)
	}

	expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIwNTEzMDg4MDAsInVpZCI6IjEyMzQ1In0.JJeEVUmuPddmZgG0yredFphjdcjrHyCmwyiyRm94UUo"
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

// TestJwtTokenProvider_GenerateToken - Test for GenerateToken
func TestJwtTokenProvider_GenerateToken(t *testing.T) {
	secretKey := []byte("secretKey")
	expectedTokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIwNTEzMDg4MDAsInVpZCI6InF3ZXJ0eSJ9.atqqmr4qUe2djJuxcRN7bSBTadMvnE1j-J4EQ8IaLJ0"
	userID := "qwerty"

	tokenProvider := NewJwtTokenProvider(secretKey, func() time.Time {
		return time.Date(2035, 1, 1, 0, 0, 0, 0, time.UTC)
	})

	genereratedToken, err := tokenProvider.GenerateToken(context.Background(), userID)
	require.NoError(t, err)
	require.Equal(t, expectedTokenString, genereratedToken)
}
