package user

import (
	"strings"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Username     string
}

func NewUser(email, passwordHash, username string) (*User, error) {
	trimmedUsername := strings.TrimSpace(username)
	if trimmedUsername == "" {
		return nil, ErrUsernameRequired
	}
	if len(trimmedUsername) > 50 {
		return nil, ErrUsernameTooLong
	}
	return &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		Username:     trimmedUsername,
	}, nil
}
