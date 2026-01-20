package user

import (
	"context"
	"foodplanner/internal/db"
)

type UserService struct {
	db   db.DBTX
	repo *userRepo
}

func NewUserService(db db.DBTX, repo *userRepo) *UserService {
	return &UserService{
		db:   db,
		repo: repo,
	}
}

func (s *UserService) CreateUser(email, passwordHash string, ctx context.Context) (*User, error) {
	existingUser, err := s.repo.GetUserByEmail(email, ctx, s.db)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailInUse
	}
	newUser := NewUser(email, passwordHash)
	createdUser, err := s.repo.CreateUser(newUser, ctx, s.db)
	if err != nil {
		return nil, err
	}
	if createdUser == nil {
		return nil, nil
	}
	return createdUser, nil
}

func (s *UserService) GetUserByEmail(email string, ctx context.Context) (*User, error) {
	return s.repo.GetUserByEmail(email, ctx, s.db)
}
