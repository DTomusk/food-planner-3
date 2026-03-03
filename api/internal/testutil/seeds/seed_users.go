package seeds

import (
	"context"
	"foodplanner/internal/db"
	"foodplanner/internal/user"

	"github.com/google/uuid"
)

func InsertUser(ctx context.Context, db db.DBTX, user *user.User) error {
	query := `INSERT INTO users (id, email, password_hash, username) VALUES ($1, $2, $3, $4)`
	_, err := db.ExecContext(ctx, query, user.ID, user.Email, user.PasswordHash, user.Username)
	return err
}

func SeedTestUser(ctx context.Context, db db.DBTX) (*user.User, error) {
	userID := uuid.New()
	testUser := user.User{
		ID:           userID,
		Email:        "test@example.com",
		PasswordHash: "test-password-hash",
		Username:     "testuser",
	}
	err := InsertUser(context.Background(), db, &testUser)
	if err != nil {
		return nil, err
	}
	return &testUser, nil
}
