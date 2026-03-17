package refreshtokens

import (
	"context"
	"foodplanner/internal/db"
)

type refreshTokenRepo struct{}

func NewRefreshTokenRepo() *refreshTokenRepo {
	return &refreshTokenRepo{}
}

func (r *refreshTokenRepo) SaveNewToken(ctx context.Context, db db.DBTX, token *RefreshToken) error {
	// Implementation for saving the new token to the database
	return nil
}
