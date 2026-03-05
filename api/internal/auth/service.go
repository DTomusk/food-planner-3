package auth

import (
	"context"
	"foodplanner/internal/db"
	"foodplanner/internal/user"
)

type AuthService struct {
	db          db.DBTX
	userService *user.UserService
	jwtService  *JWTService
}

func NewAuthService(db db.DBTX, userService *user.UserService, jwtService *JWTService) *AuthService {
	return &AuthService{
		db:          db,
		userService: userService,
		jwtService:  jwtService,
	}
}

func (s *AuthService) SignUp(email, password, username string, ctx context.Context) (*user.User, string, error) {
	if err := validateEmail(email); err != nil {
		return nil, "", err
	}
	if err := validatePassword(password); err != nil {
		return nil, "", err
	}
	emailUser, err := s.userService.GetUserByEmail(email, ctx)
	if err != nil {
		return nil, "", err
	}
	if emailUser != nil {
		return nil, "", ErrEmailAlreadyInUse
	}
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, "", err
	}
	user, err := s.userService.CreateUser(email, hashedPassword, username, ctx)
	if err != nil {
		return nil, "", err
	}
	token, err := s.jwtService.GenerateToken(user.ID.String())
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (s *AuthService) SignIn(email, password string, ctx context.Context) (*user.User, string, error) {
	user, err := s.userService.GetUserByEmail(email, ctx)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", ErrInvalidCredentials
	}
	err = comparePasswordHash(password, user.PasswordHash)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}
	token, err := s.jwtService.GenerateToken(user.ID.String())
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}
