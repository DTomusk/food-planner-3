package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateToken_Success(t *testing.T) {
	jwtService := NewJWTService("testsecret", 15)
	userID := "user-123"
	token, err := jwtService.GenerateToken(userID)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestValidateToken_Success(t *testing.T) {
	jwtService := NewJWTService("testsecret", 15)
	userID := "user-123"
	token, err := jwtService.GenerateToken(userID)
	require.NoError(t, err)
	claims, err := jwtService.ValidateToken(token)
	require.NoError(t, err)
	require.Equal(t, userID, claims.UserID)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	jwtService := NewJWTService("testsecret", 15)
	invalidToken := "invalid.token.string"
	_, err := jwtService.ValidateToken(invalidToken)
	require.Error(t, err)
	require.Equal(t, ErrInvalidToken, err)
}
