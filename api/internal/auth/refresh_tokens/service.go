package refreshtokens

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"foodplanner/internal/db"
	"foodplanner/internal/logging"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenService struct {
	txRunner      db.TxRunner
	repo          *refreshTokenRepo
	hasher        *hasher
	expiresInDays int
}

func NewRefreshTokenService(txRunner db.TxRunner, repo *refreshTokenRepo, secret string, expiresInDays int) *RefreshTokenService {
	return &RefreshTokenService{
		txRunner:      txRunner,
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
	err = s.repo.saveNewToken(ctx, s.txRunner.DB(), refreshToken)
	if err != nil {
		logger.Error("Failed to persist refresh token", "error", err)
		return nil, err
	}

	return refreshToken, nil
}

func (s *RefreshTokenService) RefreshSession(ctx context.Context, tokenString, ipAddress string) (*RefreshToken, error) {
	logger := logging.FromContext(ctx)
	hashedToken := s.hasher.hash(tokenString)
	existingToken, err := s.repo.getByHashedToken(ctx, s.txRunner.DB(), hashedToken)
	if err != nil {
		logger.Error("Failed to fetch refresh token", "error", err)
		return nil, err
	}
	if existingToken == nil {
		logger.Warn("No token found for hash", "ipAddress", ipAddress)
		return nil, ErrInvalidRefreshToken
	}

	var newRefreshToken *RefreshToken

	err = s.txRunner.WithTx(ctx, func(tx *sql.Tx) error {
		// always revoke existing token family
		err := s.repo.revokeFamily(ctx, tx, existingToken.FamilyID, ipAddress)
		if err != nil {
			logger.Error("Failed to revoke token family", "error", err, "familyID", existingToken.FamilyID, "ipAddress", ipAddress)
			return err
		}

		if existingToken.IsRevoked || existingToken.ExpiresAt < time.Now().Unix() {
			logger.Warn("Attempt to use revoked or expired refresh token", "tokenID", existingToken.ID, "ipAddress", ipAddress)
			// revoke family
			return ErrInvalidRefreshToken
		}

		// create child token
		newRefreshToken, err = s.createRefreshToken(existingToken.UserID, ipAddress)
		if err != nil {
			logger.Error("Failed to create new refresh token", "error", err, "ipAddress", ipAddress)
			return err
		}

		newRefreshToken = newChildRefreshToken(existingToken, newRefreshToken.Token, newRefreshToken.TokenHash, ipAddress, s.expiresInDays)

		// persist child token
		err = s.repo.saveNewToken(ctx, tx, newRefreshToken)
		if err != nil {
			logger.Error("Failed to persist new refresh token", "error", err, "ipAddress", ipAddress)
			return err
		}

		// set replaced by
		err = s.repo.setReplacedBy(ctx, tx, existingToken.ID, newRefreshToken.ID)
		if err != nil {
			logger.Error("Failed to set replaced by for refresh token", "error", err, "ipAddress", ipAddress)
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return newRefreshToken, nil
}

func (s *RefreshTokenService) createRefreshToken(userID uuid.UUID, ipAddress string) (*RefreshToken, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	token := base64.RawURLEncoding.EncodeToString(b)
	hashedToken := s.hasher.hash(token)
	return newRefreshToken(token, hashedToken, userID, ipAddress, s.expiresInDays), nil
}
