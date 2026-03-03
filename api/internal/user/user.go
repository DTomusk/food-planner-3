package user

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Username     string
}

func NewUser(email, passwordHash, username string) *User {
	return &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		Username:     username,
	}
}
