package resolver

import (
	"context"
	"database/sql"
	"food-planner/internal/auth"
	"food-planner/internal/gql/graph/model"
	"food-planner/internal/testutil"
	"food-planner/internal/user"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthResolver_SignUp(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := auth.NewJWTService("testsecret", 15)
		authService := auth.NewAuthService(tx, userService, jwtService)
		r := &Resolver{
			AuthService: authService,
		}
		mutationResolver := &mutationResolver{r}
		input := model.SignUpInput{
			Email:    "test@example.com",
			Password: "securepassword",
		}

		// Act
		authPayload, err := mutationResolver.Signup(context.Background(), input)

		// Assert
		require.NoError(t, err)
		require.Equal(t, input.Email, authPayload.User.Email)
		require.NotEmpty(t, authPayload.User.ID)
		require.NotEmpty(t, authPayload.Jwt)
	})
}

func TestAuthResolver_SignIn(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := auth.NewJWTService("testsecret", 15)
		authService := auth.NewAuthService(tx, userService, jwtService)
		r := &Resolver{
			AuthService: authService,
		}
		mutationResolver := &mutationResolver{r}

		// First, create a user to sign in
		createdUser, _, err := authService.SignUp("test@example.com", "securepassword", context.Background())
		require.NoError(t, err)

		input := model.SignInInput{
			Email:    "test@example.com",
			Password: "securepassword",
		}

		// Act
		authPayload, err := mutationResolver.Signin(context.Background(), input)

		// Assert
		require.NoError(t, err)
		require.Equal(t, createdUser.Email, authPayload.User.Email)
		require.Equal(t, createdUser.ID, authPayload.User.ID)
		require.NotEmpty(t, authPayload.Jwt)
	})
}
