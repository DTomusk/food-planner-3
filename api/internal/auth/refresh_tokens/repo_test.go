package refreshtokens

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"foodplanner/internal/testutil"
	"foodplanner/internal/testutil/seeds"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSaveNewToken_PersistsRefreshToken(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		ctx := context.Background()
		repo := NewRefreshTokenRepo()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		token := newTestRefreshToken(testUser.ID, "hashed-refresh-token")

		// Act
		err = repo.saveNewToken(ctx, tx, token)
		require.NoError(t, err)

		var (
			storedID              uuid.UUID
			storedUserID          uuid.UUID
			storedIPAddress       string
			storedTokenHash       string
			storedExpiresAt       time.Time
			storedRevokedAt       sql.NullTime
			storedFamilyID        uuid.UUID
			replacedByTokenIsNull bool
		)

		err = tx.QueryRowContext(ctx, `
			SELECT id, user_id, ip_address, token_hash, expires_at, revoked_at, family_id, replaced_by_token_id IS NULL
			FROM refresh_tokens
			WHERE id = $1
		`, token.ID).Scan(
			&storedID,
			&storedUserID,
			&storedIPAddress,
			&storedTokenHash,
			&storedExpiresAt,
			&storedRevokedAt,
			&storedFamilyID,
			&replacedByTokenIsNull,
		)

		// Assert
		require.NoError(t, err)
		require.Equal(t, token.ID, storedID)
		require.Equal(t, token.UserID, storedUserID)
		require.Equal(t, token.IPAddress, storedIPAddress)
		require.Equal(t, token.TokenHash, storedTokenHash)
		require.Equal(t, token.ExpiresAt, storedExpiresAt.Unix())
		require.False(t, storedRevokedAt.Valid)
		require.Equal(t, token.FamilyID, storedFamilyID)
		require.True(t, replacedByTokenIsNull)
	})
}

func TestGetByHashedToken_ReturnsNilWhenNotFound(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		ctx := context.Background()
		repo := NewRefreshTokenRepo()

		// Act
		refreshToken, err := repo.getByHashedToken(ctx, tx, "missing-token-hash")

		// Assert
		require.NoError(t, err)
		require.Nil(t, refreshToken)
	})
}

func TestGetByHashedToken_ReturnsPersistedToken(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		ctx := context.Background()
		repo := NewRefreshTokenRepo()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		token := newTestRefreshToken(testUser.ID, "active-token-hash")
		err = repo.saveNewToken(ctx, tx, token)
		require.NoError(t, err)

		// Act
		storedToken, err := repo.getByHashedToken(ctx, tx, token.TokenHash)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, storedToken)
		require.Equal(t, token.ID, storedToken.ID)
		require.Equal(t, token.UserID, storedToken.UserID)
		require.Equal(t, token.IPAddress, storedToken.IPAddress)
		require.Equal(t, token.TokenHash, storedToken.TokenHash)
		require.Equal(t, token.ExpiresAt, storedToken.ExpiresAt)
		require.False(t, storedToken.IsRevoked)
		require.Equal(t, token.FamilyID, storedToken.FamilyID)
		require.Nil(t, storedToken.ReplacedByID)
	})
}

func TestGetByHashedToken_ReturnsRevokedTokenAsRevoked(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		ctx := context.Background()
		repo := NewRefreshTokenRepo()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		token := newTestRefreshToken(testUser.ID, "revoked-token-hash")
		err = insertRevokedRefreshToken(ctx, tx, token)
		require.NoError(t, err)

		// Act
		storedToken, err := repo.getByHashedToken(ctx, tx, token.TokenHash)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, storedToken)
		require.Equal(t, token.ID, storedToken.ID)
		require.Equal(t, token.ExpiresAt, storedToken.ExpiresAt)
		require.True(t, storedToken.IsRevoked)
	})
}

func newTestRefreshToken(userID uuid.UUID, tokenHash string) *RefreshToken {
	return &RefreshToken{
		ID:           uuid.New(),
		Token:        "refresh-token",
		TokenHash:    tokenHash,
		IsRevoked:    false,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour).Unix(),
		UserID:       userID,
		IPAddress:    "127.0.0.1",
		FamilyID:     uuid.New(),
		ReplacedByID: nil,
	}
}

func insertRevokedRefreshToken(ctx context.Context, tx *sql.Tx, token *RefreshToken) error {
	query := `
	INSERT INTO refresh_tokens
	(id, user_id, ip_address, token_hash, expires_at, revoked_at, family_id)
	VALUES ($1, $2, $3, $4, to_timestamp($5), NOW(), $6)
	`

	_, err := tx.ExecContext(ctx, query,
		token.ID,
		token.UserID,
		token.IPAddress,
		token.TokenHash,
		token.ExpiresAt,
		token.FamilyID,
	)
	return err
}
