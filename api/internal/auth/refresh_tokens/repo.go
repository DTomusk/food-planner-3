package refreshtokens

import (
	"context"
	"foodplanner/internal/db"
)

type refreshTokenRepo struct{}

func NewRefreshTokenRepo() *refreshTokenRepo {
	return &refreshTokenRepo{}
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
