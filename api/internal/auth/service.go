package auth

import (
	"context"
	refreshtokens "foodplanner/internal/auth/refresh_tokens"
	"foodplanner/internal/db"
	"foodplanner/internal/user"
)

type AuthService struct {
	db                  db.DBTX
	userService         *user.UserService
	jwtService          *JWTService
	refreshTokenService *refreshtokens.RefreshTokenService
}

func NewAuthService(db db.DBTX, userService *user.UserService, jwtService *JWTService, refreshTokenService *refreshtokens.RefreshTokenService) *AuthService {
	return &AuthService{
		db:                  db,
		userService:         userService,
		jwtService:          jwtService,
		refreshTokenService: refreshTokenService,
	}
}

func (s *AuthService) SignUp(email, password, username, ipAddress string, ctx context.Context) (*user.User, string, *refreshtokens.RefreshToken, error) {
	if err := validateEmail(email); err != nil {
		return nil, "", nil, err
	}
	if err := validatePassword(password); err != nil {
		return nil, "", nil, err
	}
	emailUser, err := s.userService.GetUserByEmail(email, ctx)
	if err != nil {
		return nil, "", nil, err
	}
	if emailUser != nil {
		return nil, "", nil, ErrEmailAlreadyInUse
	}
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, "", nil, err
	}
	user, err := s.userService.CreateUser(email, hashedPassword, username, ctx)
	if err != nil {
		return nil, "", nil, err
	}
	token, err := s.jwtService.GenerateToken(user.ID.String())
	if err != nil {
		return nil, "", nil, err
	}

	refresh_token, err := s.refreshTokenService.NewSession(ctx, user.ID, ipAddress)
	if err != nil {
		return nil, "", nil, err
	}

	return user, token, refresh_token, nil
}

func (s *AuthService) SignIn(email, password, ipAddress string, ctx context.Context) (*user.User, string, *refreshtokens.RefreshToken, error) {
	user, err := s.userService.GetUserByEmail(email, ctx)
	if err != nil {
		return nil, "", nil, err
	}
	if user == nil {
		return nil, "", nil, ErrInvalidCredentials
	}
	err = comparePasswordHash(password, user.PasswordHash)
	if err != nil {
		return nil, "", nil, ErrInvalidCredentials
	}
	token, err := s.jwtService.GenerateToken(user.ID.String())
	if err != nil {
		return nil, "", nil, err
	}
	refresh_token, err := s.refreshTokenService.NewSession(ctx, user.ID, ipAddress)
	if err != nil {
		return nil, "", nil, err
	}
	return user, token, refresh_token, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken, ipAddress string) (*user.User, string, *refreshtokens.RefreshToken, error) {
	// Hash refresh token
	// Generate new refresh token that's a child of the previous token, and revoke old one in DB (rotate)
	// Fetch refresh token from DB with hash
	// If not found, revoked or expired, return error
	newToken, err := s.refreshTokenService.RefreshSession(ctx, refreshToken, ipAddress)
	if err != nil {
		return nil, "", nil, err
	}

	user, err := s.userService.GetUserByID(ctx, newToken.UserID)
	if err != nil {
		return nil, "", nil, err
	}
	if user == nil {
		return nil, "", nil, ErrInvalidCredentials
	}

	// Fetch user from DB
	// If user not found, return error
	// Generate new JWT token
	jwt, err := s.jwtService.GenerateToken(user.ID.String())
	if err != nil {
		return nil, "", nil, err
	}
	return user, jwt, newToken, nil
}
