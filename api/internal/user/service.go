package user

import (
	"context"
	"foodplanner/internal/db"

	"github.com/google/uuid"
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

func (s *UserService) CreateUser(email, passwordHash, username string, ctx context.Context) (*User, error) {
	existingUser, err := s.repo.getUserByEmail(email, ctx, s.db)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailInUse
	}
	newUser, err := NewUser(email, passwordHash, username)
	if err != nil {
		return nil, err
	}
	createdUser, err := s.repo.createUser(newUser, ctx, s.db)
	if err != nil {
		return nil, err
	}
	if createdUser == nil {
		return nil, nil
	}
	return createdUser, nil
}

func (s *UserService) GetUserByEmail(email string, ctx context.Context) (*User, error) {
	return s.repo.getUserByEmail(email, ctx, s.db)
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repo.getUserByID(ctx, s.db, id)
}
