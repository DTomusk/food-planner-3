package auth

import (
	"context"
	"food-planner/internal/db"
)

type AuthService struct {
	db db.DBTX
}

func NewAuthService(db db.DBTX) *AuthService {
	return &AuthService{
		db: db,
	}
}

func (s *AuthService) SignUp(email, password string, ctx context.Context) {
	// Validate email and password
	// Check email not in use
	// Hash password
	// Create user in DB
}

func (s *AuthService) SignIn(email, password string, ctx context.Context) {
	// Retrieve user by email
	// Compare password hash
	// Generate auth token
}
