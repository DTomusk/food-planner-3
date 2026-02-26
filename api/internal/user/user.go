package user

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
}

func NewUser(email, passwordHash string) *User {
	return &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
	}
}
