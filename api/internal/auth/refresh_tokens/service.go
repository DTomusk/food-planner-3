package refreshtokens

import (
	"context"
	"crypto/rand"
	"foodplanner/internal/db"
	"foodplanner/internal/logging"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenService struct {
	db            db.DBTX
	repo          *refreshTokenRepo
	hasher        *hasher
	expiresInDays int
}

func NewRefreshTokenService(db db.DBTX, repo *refreshTokenRepo, secret string, expiresInDays int) *RefreshTokenService {
	return &RefreshTokenService{
		db:            db,
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
	err = s.repo.saveNewToken(ctx, s.db, refreshToken)
	if err != nil {
		logger.Error("Failed to persist refresh token", "error", err)
		return nil, err
	}

	return refreshToken, nil
}

func (s *RefreshTokenService) RefreshSession(ctx context.Context, tokenString, ipAddress string) (*RefreshToken, error) {
	logger := logging.FromContext(ctx)
	hashedToken := s.hasher.hash(tokenString)
	existingToken, err := s.repo.getByHashedToken(ctx, s.db, hashedToken)
	if err != nil {
		logger.Error("Failed to fetch refresh token", "error", err)
		return nil, err
	}
	if existingToken == nil || existingToken.IsRevoked || existingToken.ExpiresAt < time.Now().Unix() {
		return nil, ErrInvalidRefreshToken
	}

	return nil, nil
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
