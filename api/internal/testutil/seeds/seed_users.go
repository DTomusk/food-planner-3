package seeds

import (
	"context"
	"foodplanner/internal/db"
	"foodplanner/internal/user"
	"math/rand"

	"github.com/google/uuid"
)

func InsertUser(ctx context.Context, db db.DBTX, user *user.User) error {
	query := `INSERT INTO users (id, email, password_hash, username) VALUES ($1, $2, $3, $4)`
	_, err := db.ExecContext(ctx, query, user.ID, user.Email, user.PasswordHash, user.Username)
	return err
}

func SeedTestUser(ctx context.Context, db db.DBTX) (*user.User, error) {
	userID := uuid.New()
	email := randSeq(10) + "@example.com"
	passwordHash := "test-password-hash"
	username := randSeq(8)

	testUser := user.User{
		ID:           userID,
		Email:        email,
		PasswordHash: passwordHash,
		Username:     username,
	}
	err := InsertUser(context.Background(), db, &testUser)
	if err != nil {
		return nil, err
	}
	return &testUser, nil
}

// TODO: move somewhere else
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
