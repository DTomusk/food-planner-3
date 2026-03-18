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
	service := NewRefreshTokenService(nil, nil, secret, 7)

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
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()
		ipAddress := "127.0.0.1"
		secret := "my-secret-key"
		service := NewRefreshTokenService(txRunner, NewRefreshTokenRepo(), secret, 7)

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

func TestRefreshSession_RotatesRefreshTokenAndRevokesFamily(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()
		oldIPAddress := "127.0.0.1"
		newIPAddress := "10.0.0.1"
		secret := "my-secret-key"
		repo := NewRefreshTokenRepo()
		service := NewRefreshTokenService(txRunner, repo, secret, 7)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		existingToken, err := service.NewSession(ctx, testUser.ID, oldIPAddress)
		require.NoError(t, err)
		require.NotNil(t, existingToken)

		siblingToken := newTestRefreshToken(testUser.ID, "sibling-token-hash")
		siblingToken.FamilyID = existingToken.FamilyID
		err = repo.saveNewToken(ctx, tx, siblingToken)
		require.NoError(t, err)

		// Act
		newRefreshToken, err := service.RefreshSession(ctx, existingToken.Token, newIPAddress)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, newRefreshToken)
		require.Equal(t, existingToken.UserID, newRefreshToken.UserID)
		require.Equal(t, newIPAddress, newRefreshToken.IPAddress)
		require.Equal(t, existingToken.FamilyID, newRefreshToken.FamilyID)
		require.NotEqual(t, existingToken.ID, newRefreshToken.ID)

		existingRevokedAt, err := getRevokedAtByTokenID(ctx, tx, existingToken.ID)
		require.NoError(t, err)
		require.True(t, existingRevokedAt.Valid)

		siblingRevokedAt, err := getRevokedAtByTokenID(ctx, tx, siblingToken.ID)
		require.NoError(t, err)
		require.True(t, siblingRevokedAt.Valid)

		replacedByTokenID, err := getReplacedByTokenID(ctx, tx, existingToken.ID)
		require.NoError(t, err)
		require.True(t, replacedByTokenID.Valid)
		require.Equal(t, newRefreshToken.ID.String(), replacedByTokenID.String)

		storedNewToken, err := repo.getByHashedToken(ctx, tx, newRefreshToken.TokenHash)
		require.NoError(t, err)
		require.NotNil(t, storedNewToken)
		require.False(t, storedNewToken.IsRevoked)
	})
}

func TestRefreshSession_ReturnsInvalidRefreshToken_WhenTokenNotFound(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()
		ipAddress := "127.0.0.1"
		secret := "my-secret-key"
		service := NewRefreshTokenService(txRunner, NewRefreshTokenRepo(), secret, 7)

		// Act
		refreshToken, err := service.RefreshSession(ctx, "missing-refresh-token", ipAddress)

		// Assert
		require.ErrorIs(t, err, ErrInvalidRefreshToken)
		require.Nil(t, refreshToken)
	})
}

func TestRefreshSession_ReturnsInvalidRefreshToken_WhenTokenAlreadyRevoked(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()
		ipAddress := "127.0.0.1"
		secret := "my-secret-key"
		repo := NewRefreshTokenRepo()
		service := NewRefreshTokenService(txRunner, repo, secret, 7)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		existingToken, err := service.NewSession(ctx, testUser.ID, ipAddress)
		require.NoError(t, err)

		err = repo.revokeFamily(ctx, tx, existingToken.FamilyID, ipAddress)
		require.NoError(t, err)

		// Act
		refreshToken, err := service.RefreshSession(ctx, existingToken.Token, "10.0.0.1")

		// Assert
		require.ErrorIs(t, err, ErrInvalidRefreshToken)
		require.Nil(t, refreshToken)

		tokenCount, err := countTokensByFamilyID(ctx, tx, existingToken.FamilyID)
		require.NoError(t, err)
		require.Equal(t, 1, tokenCount)
	})
}

func TestRefreshSession_ReturnsInvalidRefreshToken_WhenTokenExpired(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()
		ipAddress := "127.0.0.1"
		secret := "my-secret-key"
		repo := NewRefreshTokenRepo()
		service := NewRefreshTokenService(txRunner, repo, secret, 7)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		existingToken, err := service.NewSession(ctx, testUser.ID, ipAddress)
		require.NoError(t, err)

		// Manually expire token in the near future
		// Cant insert expired token due to db constraint
		_, err = tx.ExecContext(ctx, `
			UPDATE refresh_tokens
			SET expires_at = NOW() + INTERVAL '1 second'
			WHERE id = $1
		`, existingToken.ID)
		require.NoError(t, err)

		time.Sleep(2 * time.Second)

		// Act
		refreshToken, err := service.RefreshSession(ctx, existingToken.Token, "10.0.0.1")

		// Assert
		require.ErrorIs(t, err, ErrInvalidRefreshToken)
		require.Nil(t, refreshToken)

		revokedAt, err := getRevokedAtByTokenID(ctx, tx, existingToken.ID)
		require.NoError(t, err)
		require.True(t, revokedAt.Valid)

		tokenCount, err := countTokensByFamilyID(ctx, tx, existingToken.FamilyID)
		require.NoError(t, err)
		require.Equal(t, 1, tokenCount)
	})
}

func countTokensByFamilyID(ctx context.Context, tx *sql.Tx, familyID uuid.UUID) (int, error) {
	var tokenCount int
	err := tx.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM refresh_tokens
		WHERE family_id = $1
	`, familyID).Scan(&tokenCount)
	return tokenCount, err
}
