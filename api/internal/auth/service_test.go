package auth

import (
	"context"
	"database/sql"
	"food-planner/internal/testutil"
	"food-planner/internal/user"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignUp_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		authService := NewAuthService(tx, userService, jwtService)

		email := "blah@test.com"
		password := "securepassword"

		// Act
		user, token, err := authService.SignUp(email, password, context.Background())

		// Assert
		require.NoError(t, err)
		require.Equal(t, email, user.Email)
		require.NotEmpty(t, user.ID)
		require.NotEqual(t, "securepassword", user.PasswordHash)
		require.NotEmpty(t, token)
	})
}

func TestSignup_InvalidEmail(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		authService := NewAuthService(tx, userService, jwtService)
		invalidEmail := "invalid-email"
		password := "securepassword"

		// Act
		_, _, err := authService.SignUp(invalidEmail, password, context.Background())

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrInvalidEmail, err)
	})
}

func TestSignup_ShortPassword(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		authService := NewAuthService(tx, userService, jwtService)
		email := "test@fun.com"
		invalidPassword := "123"

		// Act
		_, _, err := authService.SignUp(email, invalidPassword, context.Background())

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrInvalidPassword, err)
	})
}

func TestSignup_LongPassword(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		authService := NewAuthService(tx, userService, jwtService)
		email := "blah@baz.com"
		// Note: password length limit is hardcoded 64 characters
		invalidPassword := "12345678901234567890123456789012345678901234567890123456789012345"

		// Act
		_, _, err := authService.SignUp(email, invalidPassword, context.Background())

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrInvalidPassword, err)
	})
}

func TestSignUp_DuplicateEmail(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		authService := NewAuthService(tx, userService, jwtService)

		email := "blah@test.com"
		password := "securepassword"

		_, _, err := authService.SignUp(email, password, context.Background())
		require.NoError(t, err)

		// Act
		_, _, err = authService.SignUp(email, password, context.Background())
		// Assert
		require.Error(t, err)
		require.Equal(t, ErrEmailAlreadyInUse, err)
	})
}

func TestSignIn_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		authService := NewAuthService(tx, userService, jwtService)

		email := "test@example.com"
		password := "securepassword"

		createdUser, _, err := authService.SignUp(email, password, context.Background())
		require.NoError(t, err)

		// Act
		user, token, err := authService.SignIn(email, password, context.Background())

		// Assert
		require.NoError(t, err)
		require.Equal(t, createdUser.ID, user.ID)
		require.Equal(t, createdUser.Email, user.Email)
		require.NotEmpty(t, token)
	})
}

func TestSignIn_NoUser(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		authService := NewAuthService(tx, userService, jwtService)

		email := "test@example.com"
		password := "wrongpassword"

		// Act
		_, _, err := authService.SignIn(email, password, context.Background())

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrInvalidCredentials, err)
	})
}

func TestSignIn_WrongPassword(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		authService := NewAuthService(tx, userService, jwtService)
		email := "example@test.com"
		correctPassword := "correctpassword"

		_, _, err := authService.SignUp(email, correctPassword, context.Background())
		require.NoError(t, err)

		// Act
		wrongPassword := "wrongpassword"
		_, _, err = authService.SignIn(email, wrongPassword, context.Background())

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrInvalidCredentials, err)
	})
}
