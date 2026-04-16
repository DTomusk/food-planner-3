package auth

import (
	"context"
	"database/sql"
	refreshtokens "foodplanner/internal/auth/refresh_tokens"
	"foodplanner/internal/db"
	"foodplanner/internal/events"
	"foodplanner/internal/testutil"
	"foodplanner/internal/user"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newTestAuthService(t *testing.T, tx *sql.Tx) (*AuthService, *events.InMemoryEventBus) {
	t.Helper()

	txRunner := testutil.NewTestTxRunner(tx)
	userService := user.NewUserService(tx, user.NewUserRepo())
	jwtService := NewJWTService("testsecret", 15)
	refreshTokenService := refreshtokens.NewRefreshTokenService(txRunner, refreshtokens.NewRefreshTokenRepo(), "refresh-secret", 7)
	eventBus := events.NewInMemoryEventBus(1, 32, txRunner)
	t.Cleanup(func() {
		_ = eventBus.Close(context.Background())
	})

	authService := NewAuthService(tx, userService, jwtService, refreshTokenService, eventBus)
	return authService, eventBus
}

func TestSignUp_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		authService, _ := newTestAuthService(t, tx)

		email := "blah@test.com"
		password := "securepassword"
		username := "testuser"
		ipAddress := "127.0.0.1"

		// Act
		user, token, refreshToken, err := authService.SignUp(email, password, username, ipAddress, context.Background())

		// Assert
		require.NoError(t, err)
		require.Equal(t, email, user.Email)
		require.NotEmpty(t, user.ID)
		require.NotEqual(t, "securepassword", user.PasswordHash)
		require.Equal(t, username, user.Username)
		require.NotEmpty(t, token)
		require.NotNil(t, refreshToken)
		require.Equal(t, user.ID, refreshToken.UserID)
		require.Equal(t, ipAddress, refreshToken.IPAddress)
		require.False(t, refreshToken.IsRevoked)
		require.Greater(t, refreshToken.ExpiresAt, time.Now().Unix())
	})
}

func TestSignUp_PublishesSignupEvent(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		authService, eventBus := newTestAuthService(t, tx)
		received := make(chan events.UserSignedUpEvent, 1)
		_, err := eventBus.Subscribe(events.UserSignedUpType, events.HandlerFunc(func(ctx context.Context, tx db.DBTX, event events.Event) error {
			signupEvent, ok := event.(events.UserSignedUpEvent)
			if !ok {
				return nil
			}
			received <- signupEvent
			return nil
		}))
		require.NoError(t, err)

		email := "signup-audit@test.com"
		password := "securepassword"
		username := "audited-user"
		ipAddress := "127.0.0.1"

		// Act
		signedUpUser, _, _, err := authService.SignUp(email, password, username, ipAddress, context.Background())
		require.NoError(t, err)

		select {
		case published := <-received:
			require.Equal(t, signedUpUser.ID, published.UserID)
			require.Equal(t, username, published.Username)
			require.Equal(t, email, published.Email)
			require.Equal(t, ipAddress, published.IPAddress)
		case <-time.After(1 * time.Second):
			t.Fatal("expected signup event to be published")
		}
	})
}

func TestSignup_InvalidEmail(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		authService, _ := newTestAuthService(t, tx)
		invalidEmail := "invalid-email"
		password := "securepassword"
		username := "testuser"
		ipAddress := "127.0.0.1"

		// Act
		_, _, _, err := authService.SignUp(invalidEmail, password, username, ipAddress, context.Background())

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrInvalidEmail, err)
	})
}

func TestSignup_ShortPassword(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		authService, _ := newTestAuthService(t, tx)
		email := "test@fun.com"
		invalidPassword := "123"
		ipAddress := "127.0.0.1"

		// Act
		_, _, _, err := authService.SignUp(email, invalidPassword, "testuser", ipAddress, context.Background())

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrInvalidPassword, err)
	})
}

func TestSignup_LongPassword(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		authService, _ := newTestAuthService(t, tx)
		email := "blah@baz.com"
		// Note: password length limit is hardcoded 64 characters
		invalidPassword := "12345678901234567890123456789012345678901234567890123456789012345"
		ipAddress := "127.0.0.1"

		// Act
		_, _, _, err := authService.SignUp(email, invalidPassword, "testuser", ipAddress, context.Background())

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrInvalidPassword, err)
	})
}

func TestSignUp_DuplicateEmail(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		authService, _ := newTestAuthService(t, tx)

		email := "blah@test.com"
		password := "securepassword"
		username := "testuser"
		ipAddress := "127.0.0.1"

		_, _, _, err := authService.SignUp(email, password, username, ipAddress, context.Background())
		require.NoError(t, err)

		// Act
		_, _, _, err = authService.SignUp(email, password, username, ipAddress, context.Background())
		// Assert
		require.Error(t, err)
		require.Equal(t, ErrEmailAlreadyInUse, err)
	})
}

func TestSignIn_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		authService, _ := newTestAuthService(t, tx)
		ipAddress := "127.0.0.1"
		userAgent := "test-agent/1.0"

		email := "test@example.com"
		password := "securepassword"
		username := "testuser"

		createdUser, _, _, err := authService.SignUp(email, password, username, ipAddress, context.Background())
		require.NoError(t, err)

		// Act
		user, token, refreshToken, err := authService.SignIn(email, password, ipAddress, userAgent, context.Background())

		// Assert
		require.NoError(t, err)
		require.Equal(t, createdUser.ID, user.ID)
		require.Equal(t, createdUser.Email, user.Email)
		require.NotEmpty(t, token)
		require.NotNil(t, refreshToken)
		require.Equal(t, user.ID, refreshToken.UserID)
		require.Equal(t, ipAddress, refreshToken.IPAddress)
		require.False(t, refreshToken.IsRevoked)
		require.Greater(t, refreshToken.ExpiresAt, time.Now().Unix())
	})
}

func TestSignIn_PublishesSigninEvent(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		authService, eventBus := newTestAuthService(t, tx)
		received := make(chan events.UserSignedInEvent, 1)
		_, err := eventBus.Subscribe(events.UserSignedInType, events.HandlerFunc(func(ctx context.Context, tx db.DBTX, event events.Event) error {
			signinEvent, ok := event.(events.UserSignedInEvent)
			if !ok {
				return nil
			}
			received <- signinEvent
			return nil
		}))
		require.NoError(t, err)

		email := "signin-audit@test.com"
		password := "securepassword"
		username := "signin-audited-user"
		ipAddress := "127.0.0.1"
		userAgent := "test-agent/1.0"

		signedUpUser, _, _, err := authService.SignUp(email, password, username, ipAddress, context.Background())
		require.NoError(t, err)

		// Act
		_, _, _, err = authService.SignIn(email, password, ipAddress, userAgent, context.Background())
		require.NoError(t, err)

		select {
		case published := <-received:
			require.Equal(t, signedUpUser.ID, published.UserID)
			require.Equal(t, username, published.Username)
			require.Equal(t, email, published.Email)
			require.Equal(t, ipAddress, published.IPAddress)
			require.Equal(t, userAgent, published.UserAgent)
		case <-time.After(1 * time.Second):
			t.Fatal("expected signin event to be published")
		}
	})
}

func TestSignIn_NoUser(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		authService, eventBus := newTestAuthService(t, tx)
		received := make(chan events.UserSigninFailedEvent, 1)
		_, err := eventBus.Subscribe(events.UserSigninFailedType, events.HandlerFunc(func(ctx context.Context, tx db.DBTX, event events.Event) error {
			signinFailedEvent, ok := event.(events.UserSigninFailedEvent)
			if !ok {
				return nil
			}
			received <- signinFailedEvent
			return nil
		}))
		require.NoError(t, err)

		email := "test@example.com"
		password := "wrongpassword"
		ipAddress := "127.0.0.1"
		userAgent := "test-agent/1.0"

		// Act
		_, _, _, err = authService.SignIn(email, password, ipAddress, userAgent, context.Background())

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrInvalidCredentials, err)

		select {
		case published := <-received:
			require.Nil(t, published.UserID)
			require.Equal(t, email, published.Email)
			require.Equal(t, ipAddress, published.IPAddress)
			require.Equal(t, userAgent, published.UserAgent)
			require.Equal(t, "user_not_found", published.FailureReason)
		case <-time.After(1 * time.Second):
			t.Fatal("expected signin failure event to be published")
		}
	})
}

func TestSignIn_WrongPassword(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		authService, eventBus := newTestAuthService(t, tx)
		received := make(chan events.UserSigninFailedEvent, 1)
		_, err := eventBus.Subscribe(events.UserSigninFailedType, events.HandlerFunc(func(ctx context.Context, tx db.DBTX, event events.Event) error {
			signinFailedEvent, ok := event.(events.UserSigninFailedEvent)
			if !ok {
				return nil
			}
			received <- signinFailedEvent
			return nil
		}))
		require.NoError(t, err)
		email := "example@test.com"
		correctPassword := "correctpassword"
		username := "testuser"
		ipAddress := "127.0.0.1"
		userAgent := "test-agent/1.0"

		createdUser, _, _, err := authService.SignUp(email, correctPassword, username, ipAddress, context.Background())
		require.NoError(t, err)

		// Act
		wrongPassword := "wrongpassword"
		_, _, _, err = authService.SignIn(email, wrongPassword, ipAddress, userAgent, context.Background())

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrInvalidCredentials, err)

		select {
		case published := <-received:
			require.NotNil(t, published.UserID)
			require.Equal(t, createdUser.ID, *published.UserID)
			require.Equal(t, email, published.Email)
			require.Equal(t, ipAddress, published.IPAddress)
			require.Equal(t, userAgent, published.UserAgent)
			require.Equal(t, "invalid_password", published.FailureReason)
		case <-time.After(1 * time.Second):
			t.Fatal("expected signin failure event to be published")
		}
	})
}

func TestRefresh_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		authService, _ := newTestAuthService(t, tx)

		email := "refresh-success@example.com"
		password := "securepassword"
		username := "testuser"
		signInIPAddress := "127.0.0.1"
		refreshIPAddress := "10.0.0.1"

		createdUser, _, initialRefreshToken, err := authService.SignUp(email, password, username, signInIPAddress, context.Background())
		require.NoError(t, err)

		// Act
		user, jwt, refreshedToken, err := authService.Refresh(context.Background(), initialRefreshToken.Token, refreshIPAddress)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, createdUser.ID, user.ID)
		require.Equal(t, createdUser.Email, user.Email)
		require.NotEmpty(t, jwt)
		require.NotNil(t, refreshedToken)
		require.Equal(t, user.ID, refreshedToken.UserID)
		require.Equal(t, refreshIPAddress, refreshedToken.IPAddress)
		require.Equal(t, initialRefreshToken.FamilyID, refreshedToken.FamilyID)
		require.NotEqual(t, initialRefreshToken.ID, refreshedToken.ID)
		require.False(t, refreshedToken.IsRevoked)
		require.Greater(t, refreshedToken.ExpiresAt, time.Now().Unix())
	})
}

func TestRefresh_InvalidToken(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		authService, _ := newTestAuthService(t, tx)

		// Act
		user, jwt, refreshedToken, err := authService.Refresh(context.Background(), "missing-refresh-token", "127.0.0.1")

		// Assert
		require.ErrorIs(t, err, refreshtokens.ErrInvalidRefreshToken)
		require.Nil(t, user)
		require.Empty(t, jwt)
		require.Nil(t, refreshedToken)
	})
}

func TestRefresh_ReusedTokenFails(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		authService, _ := newTestAuthService(t, tx)

		email := "refresh-reuse@example.com"
		password := "securepassword"
		username := "testuser"
		ipAddress := "127.0.0.1"

		_, _, initialRefreshToken, err := authService.SignUp(email, password, username, ipAddress, context.Background())
		require.NoError(t, err)

		_, _, _, err = authService.Refresh(context.Background(), initialRefreshToken.Token, "10.0.0.1")
		require.NoError(t, err)

		// Act
		user, jwt, refreshedToken, err := authService.Refresh(context.Background(), initialRefreshToken.Token, "10.0.0.2")

		// Assert
		require.ErrorIs(t, err, refreshtokens.ErrInvalidRefreshToken)
		require.Nil(t, user)
		require.Empty(t, jwt)
		require.Nil(t, refreshedToken)
	})
}

func TestRevoke_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		authService, _ := newTestAuthService(t, tx)

		email := "revoke-success@example.com"
		password := "securepassword"
		username := "testuser"
		ipAddress := "127.0.0.1"

		_, _, initialRefreshToken, err := authService.SignUp(email, password, username, ipAddress, context.Background())
		require.NoError(t, err)

		// Act
		err = authService.Revoke(context.Background(), initialRefreshToken.Token)

		// Assert
		require.NoError(t, err)

		user, jwt, refreshedToken, err := authService.Refresh(context.Background(), initialRefreshToken.Token, "10.0.0.1")
		require.ErrorIs(t, err, refreshtokens.ErrInvalidRefreshToken)
		require.Nil(t, user)
		require.Empty(t, jwt)
		require.Nil(t, refreshedToken)
	})
}

func TestRevoke_MissingToken_NoError(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		authService, _ := newTestAuthService(t, tx)

		// Act
		err := authService.Revoke(context.Background(), "missing-refresh-token")

		// Assert
		require.NoError(t, err)
	})
}
