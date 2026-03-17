package resolver

import (
	"context"
	"database/sql"
	"foodplanner/internal/auth"
	refreshtokens "foodplanner/internal/auth/refresh_tokens"
	"foodplanner/internal/gql/graph/model"
	"foodplanner/internal/testutil"
	"foodplanner/internal/user"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthResolver_SignUp(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := auth.NewJWTService("testsecret", 15)
		refreshTokenService := refreshtokens.NewRefreshTokenService(tx, refreshtokens.NewRefreshTokenRepo(), "refresh-secret", 7)
		authService := auth.NewAuthService(tx, userService, jwtService, refreshTokenService)
		r := &Resolver{
			AuthService: authService,
		}
		mutationResolver := &mutationResolver{r}
		input := model.SignUpInput{
			Email:    "test@example.com",
			Password: "securepassword",
			Username: "testuser",
		}

		// Act
		ctx := context.Background()
		authPayload, err := mutationResolver.Signup(ctx, input)

		// Assert
		require.NoError(t, err)
		require.Equal(t, input.Email, authPayload.User.Email)
		require.Equal(t, input.Username, authPayload.User.Username)
		require.NotEmpty(t, authPayload.User.ID)
		require.NotEmpty(t, authPayload.Jwt)
	})
}

func TestAuthResolver_SignIn(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := auth.NewJWTService("testsecret", 15)
		refreshTokenService := refreshtokens.NewRefreshTokenService(tx, refreshtokens.NewRefreshTokenRepo(), "refresh-secret", 7)
		authService := auth.NewAuthService(tx, userService, jwtService, refreshTokenService)
		r := &Resolver{
			AuthService: authService,
		}
		mutationResolver := &mutationResolver{r}

		// First, create a user to sign in
		ctx := context.Background()
		ipAddress := "127.0.0.1"
		createdUser, _, _, err := authService.SignUp("test@example.com", "securepassword", "testuser", ipAddress, ctx)
		require.NoError(t, err)

		input := model.SignInInput{
			Email:    "test@example.com",
			Password: "securepassword",
		}

		// Act
		authPayload, err := mutationResolver.Signin(ctx, input)

		// Assert
		require.NoError(t, err)
		require.Equal(t, createdUser.Email, authPayload.User.Email)
		require.Equal(t, createdUser.Username, authPayload.User.Username)
		require.Equal(t, createdUser.ID.String(), authPayload.User.ID)
		require.NotEmpty(t, authPayload.Jwt)
	})
}
