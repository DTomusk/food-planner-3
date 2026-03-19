package refreshtokens

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID uuid.UUID
	// Actual token, do not persist in db
	Token string
	// Hash of the token, persist in db
	TokenHash    string
	IsRevoked    bool
	ExpiresAt    int64
	UserID       uuid.UUID
	IPAddress    string
	FamilyID     uuid.UUID
	ReplacedByID *uuid.UUID
}

func newRefreshToken(token, tokenHash string, userID uuid.UUID, ipAddress string, expiresInDays int) *RefreshToken {
	expiresAt := time.Now().Add(time.Duration(expiresInDays) * 24 * time.Hour).Unix()
	familyID := uuid.New()
	tokenID := uuid.New()

	return &RefreshToken{
		ID:           tokenID,
		Token:        token,
		TokenHash:    tokenHash,
		UserID:       userID,
		IPAddress:    ipAddress,
		ExpiresAt:    expiresAt,
		FamilyID:     familyID,
		ReplacedByID: nil,
		IsRevoked:    false,
	}
}

func newChildRefreshToken(parentToken *RefreshToken, token, tokenHash string, ipAddress string, expiresInDays int) *RefreshToken {
	expiresAt := time.Now().Add(time.Duration(expiresInDays) * 24 * time.Hour).Unix()
	tokenID := uuid.New()

	return &RefreshToken{
		ID:           tokenID,
		Token:        token,
		TokenHash:    tokenHash,
		UserID:       parentToken.UserID,
		IPAddress:    ipAddress,
		ExpiresAt:    expiresAt,
		FamilyID:     parentToken.FamilyID,
		ReplacedByID: nil,
		IsRevoked:    false,
	}
}
