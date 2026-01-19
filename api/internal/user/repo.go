package user

import (
	"context"
	"database/sql"
	"food-planner/internal/db"
)

type userRepo struct{}

func NewUserRepo() *userRepo {
	return &userRepo{}
}

func (r *userRepo) CreateUser(user *User, ctx context.Context, db db.DBTX) (*User, error) {
	var newUser User
	query := `INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3) RETURNING id, email, password_hash`
	err := db.QueryRowContext(ctx, query, user.ID, user.Email, user.PasswordHash).Scan(&newUser.ID, &newUser.Email, &newUser.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &newUser, nil
}

// Do not return error if no rows
func (r *userRepo) GetUserByEmail(email string, ctx context.Context, db db.DBTX) (*User, error) {
	var user User
	query := `SELECT id, email, password_hash FROM users WHERE email = $1`
	err := db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Return error if no rows
func (r *userRepo) GetUserByID(id string, ctx context.Context, db db.DBTX) (*User, error) {
	var user User
	query := `SELECT id, email, password_hash FROM users WHERE id = $1`
	err := db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
