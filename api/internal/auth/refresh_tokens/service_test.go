package refreshtokens

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateRefreshToken(t *testing.T) {
	// Arrange
	userID := uuid.New()
	ipAddress := "127.0.0.1"
	secret := "my-secret-key"
	service := NewRefreshTokenService(secret)

	// Act
	refreshToken, err := service.createRefreshToken(userID, ipAddress, 7)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, refreshToken)
	require.Equal(t, userID, refreshToken.UserID)
	require.Equal(t, ipAddress, refreshToken.IPAddress)
	require.False(t, refreshToken.IsRevoked)
	require.Greater(t, refreshToken.ExpiresAt, time.Now().Unix())
}
