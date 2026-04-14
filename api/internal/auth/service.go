package auth

import (
	"context"
	"foodplanner/internal/audit"
	refreshtokens "foodplanner/internal/auth/refresh_tokens"
	"foodplanner/internal/db"
	"foodplanner/internal/logging"
	"foodplanner/internal/user"

	"github.com/google/uuid"
)

type AuthService struct {
	db                  db.DBTX
	userService         *user.UserService
	jwtService          *JWTService
	refreshTokenService *refreshtokens.RefreshTokenService
	auditService        *audit.AuditService
}

func NewAuthService(
	db db.DBTX,
	userService *user.UserService,
	jwtService *JWTService,
	refreshTokenService *refreshtokens.RefreshTokenService,
	auditService *audit.AuditService,
) *AuthService {
	return &AuthService{
		db:                  db,
		userService:         userService,
		jwtService:          jwtService,
		refreshTokenService: refreshTokenService,
		auditService:        auditService,
	}
}

func (s *AuthService) SignUp(email, password, username, ipAddress string, ctx context.Context) (*user.User, string, *refreshtokens.RefreshToken, error) {
	logger := logging.FromContext(ctx)

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

	correlationID := uuid.New()
	entry, err := audit.NewUserSignupEvent(correlationID, user.ID, user.Username, ipAddress)
	if err != nil {
		logger.Warn("Failed to create signup audit event", "userID", user.ID, "err", err)
		return user, token, refresh_token, nil
	}

	if err := s.auditService.Log(ctx, entry); err != nil {
		logger.Warn("Failed to persist signup audit event", "userID", user.ID, "correlationID", correlationID, "err", err)
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

	jwt, err := s.jwtService.GenerateToken(user.ID.String())
	if err != nil {
		return nil, "", nil, err
	}
	return user, jwt, newToken, nil
}

func (s *AuthService) Revoke(ctx context.Context, refreshToken string) error {
	return s.refreshTokenService.RevokeSession(ctx, refreshToken)
}
