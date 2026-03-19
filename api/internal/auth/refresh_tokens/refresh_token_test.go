package refreshtokens

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewRefreshToken(t *testing.T) {
	// Arrange
	userID := uuid.New()
	ipAddress := "127.0.0.1"
	token := "sample-refresh-token"
	tokenHash := "hashed-refresh-token"
	expiresInDays := 7

	// Act
	refreshToken := newRefreshToken(token, tokenHash, userID, ipAddress, expiresInDays)

	// Assert
	require.Equal(t, token, refreshToken.Token)
	require.Equal(t, tokenHash, refreshToken.TokenHash)
	require.Equal(t, userID, refreshToken.UserID)
	require.Equal(t, ipAddress, refreshToken.IPAddress)
	require.False(t, refreshToken.IsRevoked)
	require.NotZero(t, refreshToken.ExpiresAt)
	require.NotZero(t, refreshToken.ID)
	require.NotZero(t, refreshToken.FamilyID)
	require.Nil(t, refreshToken.ReplacedByID)
	require.Greater(t, refreshToken.ExpiresAt, time.Now().Unix())
}
