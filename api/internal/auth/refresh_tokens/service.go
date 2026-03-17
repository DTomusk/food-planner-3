package refreshtokens

import (
	"context"
	"crypto/rand"
	"foodplanner/internal/db"
	"foodplanner/internal/logging"

	"github.com/google/uuid"
)

type RefreshTokenService struct {
	db            db.DBTX
	repo          *refreshTokenRepo
	hasher        *hasher
	expiresInDays int
}

func NewRefreshTokenService(repo *refreshTokenRepo, secret string, expiresInDays int) *RefreshTokenService {
	return &RefreshTokenService{
		repo:          repo,
		hasher:        NewHasher(secret),
		expiresInDays: expiresInDays,
	}
}

func (s *RefreshTokenService) NewSession(ctx context.Context, userID uuid.UUID, ipAddress string) (*RefreshToken, error) {
	logger := logging.FromContext(ctx)
	refreshToken, err := s.createRefreshToken(userID, ipAddress)
	if err != nil {
		logger.Error("Failed to create refresh token", "error", err)
		return nil, err
	}

	// Persist
	err = s.repo.SaveNewToken(ctx, nil, refreshToken)
	if err != nil {
		logger.Error("Failed to persist refresh token", "error", err)
		return nil, err
	}

	return refreshToken, nil
}

func (s *RefreshTokenService) createRefreshToken(userID uuid.UUID, ipAddress string) (*RefreshToken, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	token := string(b)
	hashedToken := s.hasher.hash(token)
	return newRefreshToken(token, hashedToken, userID, ipAddress, s.expiresInDays), nil
}
