package user

import (
	"context"
	"database/sql"
	"foodplanner/internal/db"
)

type userRepo struct{}

const (
	selectUsersBaseQuery   = `SELECT id, email, password_hash, username FROM users`
	selectUserByEmailQuery = selectUsersBaseQuery + ` WHERE email = $1`
	selectUserByIDQuery    = selectUsersBaseQuery + ` WHERE id = $1`
)

func NewUserRepo() *userRepo {
	return &userRepo{}
}

func (r *userRepo) CreateUser(user *User, ctx context.Context, db db.DBTX) (*User, error) {
	var newUser User
	query := `INSERT INTO users (id, email, password_hash, username) VALUES ($1, $2, $3, $4) RETURNING id, email, password_hash, username`
	err := db.QueryRowContext(ctx, query, user.ID, user.Email, user.PasswordHash, user.Username).Scan(&newUser.ID, &newUser.Email, &newUser.PasswordHash, &newUser.Username)
	if err != nil {
		return nil, err
	}
	return &newUser, nil
}

// Do not return error if no rows
func (r *userRepo) GetUserByEmail(email string, ctx context.Context, db db.DBTX) (*User, error) {
	var user User
	err := db.QueryRowContext(ctx, selectUserByEmailQuery, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Username)
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
	err := db.QueryRowContext(ctx, selectUserByIDQuery, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
