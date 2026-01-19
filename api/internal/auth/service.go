package auth

import (
	"context"
	"food-planner/internal/db"
	"food-planner/internal/user"
)

type AuthService struct {
	db          db.DBTX
	userService *user.UserService
	jwtSerivce  *JWTService
}

func NewAuthService(db db.DBTX, userService *user.UserService, jwtService *JWTService) *AuthService {
	return &AuthService{
		db:          db,
		userService: userService,
		jwtSerivce:  jwtService,
	}
}

func (s *AuthService) SignUp(email, password string, ctx context.Context) (*user.User, error) {
	if err := validateEmail(email); err != nil {
		return nil, err
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}
	emailUser, err := s.userService.GetUserByEmail(email, ctx)
	if err != nil {
		return nil, err
	}
	if emailUser != nil {
		return nil, ErrEmailAlreadyInUse
	}
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	user, err := s.userService.CreateUser(email, hashedPassword, ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) SignIn(email, password string, ctx context.Context) (*user.User, string, error) {
	user, err := s.userService.GetUserByEmail(email, ctx)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", ErrInvalidCredentials
	}
	err = comparePasswordHash(user.PasswordHash, password)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}
	token, err := s.jwtSerivce.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}
