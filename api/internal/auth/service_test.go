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
		user, err := authService.SignUp(email, password, context.Background())

		// Assert
		require.NoError(t, err)
		require.Equal(t, email, user.Email)
		require.NotEmpty(t, user.ID)
		require.NotEqual(t, "securepassword", user.PasswordHash)
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
		_, err := authService.SignUp(invalidEmail, password, context.Background())

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
		_, err := authService.SignUp(email, invalidPassword, context.Background())

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
		_, err := authService.SignUp(email, invalidPassword, context.Background())

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

		_, err := authService.SignUp(email, password, context.Background())
		require.NoError(t, err)

		// Act
		_, err = authService.SignUp(email, password, context.Background())

		// Assert
		require.Error(t, err)
		require.Equal(t, ErrEmailAlreadyInUse, err)
	})
}
