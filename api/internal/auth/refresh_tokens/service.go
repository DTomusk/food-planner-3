package refreshtokens

import (
	"crypto/rand"

	"github.com/google/uuid"
)

type RefreshTokenService struct {
	hasher *hasher
}

func NewRefreshTokenService(secret string) *RefreshTokenService {
	return &RefreshTokenService{
		hasher: NewHasher(secret),
	}
}

func (s *RefreshTokenService) createRefreshToken(userID uuid.UUID, ipAddress string, expiresInDays int) (*RefreshToken, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	token := string(b)
	hashedToken := s.hasher.hash(token)
	return newRefreshToken(token, hashedToken, userID, ipAddress, expiresInDays), nil
}
