package refreshtokens

import (
	"context"
	"database/sql"
	"foodplanner/internal/testutil"
	"foodplanner/internal/testutil/seeds"
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
	service := NewRefreshTokenService(nil, secret, 7)

	// Act
	refreshToken, err := service.createRefreshToken(userID, ipAddress)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, refreshToken)
	require.Equal(t, userID, refreshToken.UserID)
	require.Equal(t, ipAddress, refreshToken.IPAddress)
	require.False(t, refreshToken.IsRevoked)
	require.Greater(t, refreshToken.ExpiresAt, time.Now().Unix())
}

func TestNewSession(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		ctx := context.Background()
		ipAddress := "127.0.0.1"
		secret := "my-secret-key"
		service := NewRefreshTokenService(NewRefreshTokenRepo(), secret, 7)

		testuser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		// Act
		refreshToken, err := service.NewSession(ctx, testuser.ID, ipAddress)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, refreshToken)
		require.Equal(t, testuser.ID, refreshToken.UserID)
		require.Equal(t, ipAddress, refreshToken.IPAddress)
		require.False(t, refreshToken.IsRevoked)
		require.Greater(t, refreshToken.ExpiresAt, time.Now().Unix())
	})
}
