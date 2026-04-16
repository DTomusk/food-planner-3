package auth

import (
	"context"
	refreshtokens "foodplanner/internal/auth/refresh_tokens"
	"foodplanner/internal/db"
	"foodplanner/internal/events"
	"foodplanner/internal/logging"
	"foodplanner/internal/user"
	"log/slog"

	"github.com/google/uuid"
)

type AuthService struct {
	db                  db.DBTX
	userService         *user.UserService
	jwtService          *JWTService
	refreshTokenService *refreshtokens.RefreshTokenService
	eventPublisher      events.Publisher
}

func NewAuthService(
	db db.DBTX,
	userService *user.UserService,
	jwtService *JWTService,
	refreshTokenService *refreshtokens.RefreshTokenService,
	eventPublisher events.Publisher,
) *AuthService {
	return &AuthService{
		db:                  db,
		userService:         userService,
		jwtService:          jwtService,
		refreshTokenService: refreshTokenService,
		eventPublisher:      eventPublisher,
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
	signupEvent := events.NewUserSignedUpEvent(correlationID, user.ID, user.Username, user.Email, ipAddress)
	if err := s.eventPublisher.Publish(ctx, signupEvent); err != nil {
		logger.Warn("Failed to publish signup event", "userID", user.ID, "correlationID", correlationID, "err", err)
	}

	return user, token, refresh_token, nil
}

func (s *AuthService) SignIn(email, password, ipAddress, userAgent string, ctx context.Context) (*user.User, string, *refreshtokens.RefreshToken, error) {
	logger := logging.FromContext(ctx)

	user, err := s.userService.GetUserByEmail(email, ctx)
	if err != nil {
		return nil, "", nil, err
	}
	if user == nil {
		s.publishSigninFailure(ctx, logger, nil, email, ipAddress, userAgent, "user_not_found")
		return nil, "", nil, ErrInvalidCredentials
	}
	err = comparePasswordHash(password, user.PasswordHash)
	if err != nil {
		s.publishSigninFailure(ctx, logger, &user.ID, email, ipAddress, userAgent, "invalid_password")
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

	s.publishSigninSuccess(ctx, logger, user.ID, user.Username, user.Email, ipAddress, userAgent)

	return user, token, refresh_token, nil
}

func (s *AuthService) publishSigninSuccess(ctx context.Context, logger *slog.Logger, userID uuid.UUID, username, email, ipAddress, userAgent string) {
	correlationID := uuid.New()
	signinEvent := events.NewUserSignedInEvent(correlationID, userID, username, email, ipAddress, userAgent)
	if err := s.eventPublisher.Publish(ctx, signinEvent); err != nil {
		logger.Warn("Failed to publish signin event", "userID", userID, "correlationID", correlationID, "err", err)
	}
}

func (s *AuthService) publishSigninFailure(ctx context.Context, logger *slog.Logger, userID *uuid.UUID, email, ipAddress, userAgent, failureReason string) {
	correlationID := uuid.New()
	signinFailureEvent := events.NewUserSigninFailedEvent(correlationID, userID, email, ipAddress, userAgent, failureReason)
	if err := s.eventPublisher.Publish(ctx, signinFailureEvent); err != nil {
		logger.Warn("Failed to publish signin failure event", "userID", userID, "email", email, "correlationID", correlationID, "err", err)
	}
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
