package refreshtokens

import (
	"context"
	"database/sql"
	"foodplanner/internal/db"

	"github.com/google/uuid"
)

type refreshTokenRepo struct{}

func NewRefreshTokenRepo() *refreshTokenRepo {
	return &refreshTokenRepo{}
}

type refreshTokenRecord struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	IPAddress string
	TokenHash string
	ExpiresAt int64
	RevokedAt sql.NullInt64
	FamilyID  uuid.UUID
}

func (r *refreshTokenRepo) saveNewToken(ctx context.Context, db db.DBTX, token *RefreshToken) error {
	query := `
	INSERT INTO refresh_tokens 
	(id, user_id, ip_address, token_hash, expires_at, family_id)
	VALUES ($1, $2, $3, $4, to_timestamp($5), $6)
	`
	_, err := db.ExecContext(ctx, query,
		token.ID,
		token.UserID,
		token.IPAddress,
		token.TokenHash,
		token.ExpiresAt,
		token.FamilyID,
	)
	return err
}

func (r *refreshTokenRepo) getByHashedToken(ctx context.Context, db db.DBTX, hashedToken string) (*RefreshToken, error) {
	query := `
	SELECT id, user_id, ip_address, token_hash, expires_at, revoked_at, family_id
	FROM refresh_tokens
	WHERE token_hash = $1
	`
	var tokenRecord refreshTokenRecord
	err := db.QueryRowContext(ctx, query, hashedToken).Scan(
		&tokenRecord.ID,
		&tokenRecord.UserID,
		&tokenRecord.IPAddress,
		&tokenRecord.TokenHash,
		&tokenRecord.ExpiresAt,
		&tokenRecord.RevokedAt,
		&tokenRecord.FamilyID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	token := RefreshToken{
		ID:        tokenRecord.ID,
		UserID:    tokenRecord.UserID,
		IPAddress: tokenRecord.IPAddress,
		TokenHash: tokenRecord.TokenHash,
		ExpiresAt: tokenRecord.ExpiresAt,
		IsRevoked: tokenRecord.RevokedAt.Valid,
		FamilyID:  tokenRecord.FamilyID,
	}
	return &token, nil
}
