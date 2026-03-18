package auth

import (
	"context"
	"database/sql"
	refreshtokens "foodplanner/internal/auth/refresh_tokens"
	"foodplanner/internal/testutil"
	"foodplanner/internal/user"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSignUp_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		refreshTokenService := refreshtokens.NewRefreshTokenService(txRunner, refreshtokens.NewRefreshTokenRepo(), "refresh-secret", 7)
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService)

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

func TestSignup_InvalidEmail(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		refreshTokenService := refreshtokens.NewRefreshTokenService(txRunner, refreshtokens.NewRefreshTokenRepo(), "refresh-secret", 7)
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService)
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
		txRunner := testutil.NewTestTxRunner(tx)
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		refreshTokenService := refreshtokens.NewRefreshTokenService(txRunner, refreshtokens.NewRefreshTokenRepo(), "refresh-secret", 7)
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService)
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
		txRunner := testutil.NewTestTxRunner(tx)
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		refreshTokenService := refreshtokens.NewRefreshTokenService(txRunner, refreshtokens.NewRefreshTokenRepo(), "refresh-secret", 7)
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService)
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
		txRunner := testutil.NewTestTxRunner(tx)
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		refreshTokenService := refreshtokens.NewRefreshTokenService(txRunner, refreshtokens.NewRefreshTokenRepo(), "refresh-secret", 7)
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService)

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
		txRunner := testutil.NewTestTxRunner(tx)
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		refreshTokenService := refreshtokens.NewRefreshTokenService(txRunner, refreshtokens.NewRefreshTokenRepo(), "refresh-secret", 7)
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService)
		ipAddress := "127.0.0.1"

		email := "test@example.com"
		password := "securepassword"
		username := "testuser"

		createdUser, _, _, err := authService.SignUp(email, password, username, ipAddress, context.Background())
		require.NoError(t, err)

		// Act
		user, token, refreshToken, err := authService.SignIn(email, password, ipAddress, context.Background())

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

func TestSignIn_NoUser(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		refreshTokenService := refreshtokens.NewRefreshTokenService(txRunner, refreshtokens.NewRefreshTokenRepo(), "refresh-secret", 7)
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService)

		email := "test@example.com"
		password := "wrongpassword"
		ipAddress := "127.0.0.1"

		// Act
		_, _, _, err := authService.SignIn(email, password, ipAddress, context.Background())

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrInvalidCredentials, err)
	})
}

func TestSignIn_WrongPassword(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		refreshTokenService := refreshtokens.NewRefreshTokenService(txRunner, refreshtokens.NewRefreshTokenRepo(), "refresh-secret", 7)
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService)
		email := "example@test.com"
		correctPassword := "correctpassword"
		username := "testuser"
		ipAddress := "127.0.0.1"

		_, _, _, err := authService.SignUp(email, correctPassword, username, ipAddress, context.Background())
		require.NoError(t, err)

		// Act
		wrongPassword := "wrongpassword"
		_, _, _, err = authService.SignIn(email, wrongPassword, ipAddress, context.Background())

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrInvalidCredentials, err)
	})
}
