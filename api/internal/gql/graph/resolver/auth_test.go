package resolver

import (
	"context"
	"database/sql"
	"foodplanner/internal/auth"
	refreshtokens "foodplanner/internal/auth/refresh_tokens"
	"foodplanner/internal/gql/graph/model"
	"foodplanner/internal/middleware"
	"foodplanner/internal/testutil"
	"foodplanner/internal/user"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupAuthMutationResolver(tx *sql.Tx) (*mutationResolver, *auth.AuthService) {
	userService := user.NewUserService(tx, user.NewUserRepo())
	jwtService := auth.NewJWTService("testsecret", 15)
	refreshTokenService := refreshtokens.NewRefreshTokenService(tx, refreshtokens.NewRefreshTokenRepo(), "refresh-secret", 7)
	authService := auth.NewAuthService(tx, userService, jwtService, refreshTokenService)
	resolver := &Resolver{
		AuthService: authService,
	}

	return &mutationResolver{resolver}, authService
}

func signupContext(w http.ResponseWriter) context.Context {
	ctx := context.WithValue(context.Background(), middleware.IPKey, "127.0.0.1")
	if w != nil {
		ctx = context.WithValue(ctx, middleware.ResponseWriterKey, w)
	}

	return ctx
}

func TestAuthResolver_SignUp(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		mutationResolver, _ := setupAuthMutationResolver(tx)
		input := model.SignUpInput{
			Email:    "test@example.com",
			Password: "securepassword",
			Username: "testuser",
		}

		// Act
		ctx := signupContext(nil)
		authPayload, err := mutationResolver.Signup(ctx, input)

		// Assert
		require.NoError(t, err)
		require.Equal(t, input.Email, authPayload.User.Email)
		require.Equal(t, input.Username, authPayload.User.Username)
		require.NotEmpty(t, authPayload.User.ID)
		require.NotEmpty(t, authPayload.Jwt)
	})
}

func TestAuthResolver_SignUp_SetsRefreshTokenCookie(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		mutationResolver, _ := setupAuthMutationResolver(tx)
		input := model.SignUpInput{
			Email:    "cookie-test@example.com",
			Password: "securepassword",
			Username: "cookietestuser",
		}
		recorder := httptest.NewRecorder()

		// Act
		authPayload, err := mutationResolver.Signup(signupContext(recorder), input)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, authPayload)

		response := recorder.Result()
		cookies := response.Cookies()
		require.Len(t, cookies, 1)

		cookie := cookies[0]
		require.Equal(t, "refreshToken", cookie.Name)
		require.NotEmpty(t, cookie.Value)
		require.Equal(t, "/", cookie.Path)
		require.True(t, cookie.HttpOnly)
		require.False(t, cookie.Secure)
		require.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
		require.Greater(t, cookie.MaxAge, 0)
	})
}

func TestAuthResolver_SignIn(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		mutationResolver, authService := setupAuthMutationResolver(tx)

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
