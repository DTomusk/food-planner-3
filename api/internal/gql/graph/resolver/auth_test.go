package resolver

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"foodplanner/internal/audit"
	"foodplanner/internal/auth"
	refreshtokens "foodplanner/internal/auth/refresh_tokens"
	"foodplanner/internal/gql/graph/model"
	"foodplanner/internal/middleware"
	"foodplanner/internal/testutil"
	"foodplanner/internal/testutil/seeds"
	"foodplanner/internal/user"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const testRefreshSecret = "refresh-secret"

func setupAuthMutationResolver(tx *sql.Tx) (*mutationResolver, *auth.AuthService) {
	userService := user.NewUserService(tx, user.NewUserRepo())
	jwtService := auth.NewJWTService("testsecret", 15)
	txRunner := testutil.NewTestTxRunner(tx)
	refreshTokenService := refreshtokens.NewRefreshTokenService(txRunner, refreshtokens.NewRefreshTokenRepo(), testRefreshSecret, 7)
	auditService := audit.NewAuditService(tx, audit.NewRepo())
	authService := auth.NewAuthService(tx, userService, jwtService, refreshTokenService, auditService)
	resolver := &Resolver{
		AuthService: authService,
	}

	return &mutationResolver{resolver}, authService
}

func authContext(w http.ResponseWriter) context.Context {
	ctx := context.WithValue(context.Background(), middleware.IPKey, "127.0.0.1")
	if w != nil {
		ctx = context.WithValue(ctx, middleware.ResponseWriterKey, w)
	}

	return ctx
}

func authContextWithRequest(req *http.Request, w http.ResponseWriter) context.Context {
	ctx := authContext(w)
	if req != nil {
		ctx = context.WithValue(ctx, middleware.RequestKey, req)
	}
	return ctx
}

func assertRefreshTokenCookie(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()

	response := recorder.Result()
	cookies := response.Cookies()
	require.Len(t, cookies, 1)

	cookie := cookies[0]
	require.Equal(t, refreshTokenCookieName, cookie.Name)
	require.NotEmpty(t, cookie.Value)
	require.Equal(t, "/", cookie.Path)
	require.True(t, cookie.HttpOnly)
	require.False(t, cookie.Secure)
	require.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	require.Greater(t, cookie.MaxAge, 0)
}

func assertRefreshTokenCookieCleared(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()

	response := recorder.Result()
	cookies := response.Cookies()
	require.Len(t, cookies, 1)

	cookie := cookies[0]
	require.Equal(t, refreshTokenCookieName, cookie.Name)
	require.Empty(t, cookie.Value)
	require.Equal(t, "/", cookie.Path)
	require.True(t, cookie.HttpOnly)
	require.False(t, cookie.Secure)
	require.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	require.Less(t, cookie.MaxAge, 0)
	require.True(t, cookie.Expires.Before(time.Now()))
}

func hashRefreshTokenForTest(token string) string {
	mac := hmac.New(sha256.New, []byte(testRefreshSecret))
	mac.Write([]byte(token))
	return hex.EncodeToString(mac.Sum(nil))
}

func seedRefreshTokenForResolverTest(t *testing.T, ctx context.Context, tx *sql.Tx, userID, familyID uuid.UUID, token string) {
	t.Helper()

	tokenHash := hashRefreshTokenForTest(token)
	if familyID == uuid.Nil {
		familyID = uuid.New()
	}

	_, err := tx.ExecContext(ctx, `
		INSERT INTO refresh_tokens (id, user_id, ip_address, token_hash, expires_at, family_id)
		VALUES ($1, $2, $3, $4, to_timestamp($5), $6)
	`, uuid.New(), userID, "127.0.0.1", tokenHash, time.Now().Add(24*time.Hour).Unix(), familyID)
	require.NoError(t, err)
}

func assertUnauthenticatedGraphQLError(t *testing.T, err error, expectedMessage string) {
	t.Helper()

	require.Error(t, err)
	gqlErr, ok := err.(*gqlerror.Error)
	require.True(t, ok, "Expected GraphQL error type")
	require.Equal(t, expectedMessage, gqlErr.Message)
	require.Equal(t, "UNAUTHENTICATED", gqlErr.Extensions["code"])
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
		ctx := authContext(nil)
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
		authPayload, err := mutationResolver.Signup(authContext(recorder), input)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, authPayload)
		assertRefreshTokenCookie(t, recorder)
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
		authPayload, err := mutationResolver.Signin(authContext(nil), input)

		// Assert
		require.NoError(t, err)
		require.Equal(t, createdUser.Email, authPayload.User.Email)
		require.Equal(t, createdUser.Username, authPayload.User.Username)
		require.Equal(t, createdUser.ID.String(), authPayload.User.ID)
		require.NotEmpty(t, authPayload.Jwt)
	})
}

func TestAuthResolver_SignIn_SetsRefreshTokenCookie(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		mutationResolver, authService := setupAuthMutationResolver(tx)
		ctx := context.Background()
		ipAddress := "127.0.0.1"
		_, _, _, err := authService.SignUp("signin-cookie@example.com", "securepassword", "signinuser", ipAddress, ctx)
		require.NoError(t, err)

		input := model.SignInInput{
			Email:    "signin-cookie@example.com",
			Password: "securepassword",
		}
		recorder := httptest.NewRecorder()

		// Act
		authPayload, err := mutationResolver.Signin(authContext(recorder), input)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, authPayload)
		assertRefreshTokenCookie(t, recorder)
	})
}

func TestAuthResolver_Refresh_NoRequestInContext(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		mutationResolver, _ := setupAuthMutationResolver(tx)

		authPayload, err := mutationResolver.Refresh(authContext(nil))

		require.Error(t, err)
		require.Nil(t, authPayload)
		require.EqualError(t, err, "no request in context")
	})
}

func TestAuthResolver_Refresh_NoRefreshCookie(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		mutationResolver, _ := setupAuthMutationResolver(tx)

		req := httptest.NewRequest(http.MethodPost, "/query", nil)
		recorder := httptest.NewRecorder()
		authPayload, err := mutationResolver.Refresh(authContextWithRequest(req, recorder))

		require.Nil(t, authPayload)
		assertUnauthenticatedGraphQLError(t, err, "No auth cookie found")
		assertRefreshTokenCookieCleared(t, recorder)
	})
}

func TestAuthResolver_Refresh_EmptyRefreshCookie(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		mutationResolver, _ := setupAuthMutationResolver(tx)

		req := httptest.NewRequest(http.MethodPost, "/query", nil)
		req.AddCookie(&http.Cookie{Name: refreshTokenCookieName, Value: ""})
		recorder := httptest.NewRecorder()

		authPayload, err := mutationResolver.Refresh(authContextWithRequest(req, recorder))

		require.Nil(t, authPayload)
		assertUnauthenticatedGraphQLError(t, err, "Empty refresh token")
		assertRefreshTokenCookieCleared(t, recorder)
	})
}

func TestAuthResolver_Refresh_Success_SetsRefreshTokenCookie(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		mutationResolver, _ := setupAuthMutationResolver(tx)
		ctx := context.Background()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		initialTokenValue := "resolver-refresh-success-token"
		seedRefreshTokenForResolverTest(t, ctx, tx, testUser.ID, uuid.New(), initialTokenValue)

		req := httptest.NewRequest(http.MethodPost, "/query", nil)
		req.AddCookie(&http.Cookie{Name: refreshTokenCookieName, Value: initialTokenValue})
		recorder := httptest.NewRecorder()

		authPayload, err := mutationResolver.Refresh(authContextWithRequest(req, recorder))

		require.NoError(t, err)
		require.NotNil(t, authPayload)
		require.Equal(t, testUser.ID.String(), authPayload.User.ID)
		require.Equal(t, testUser.Email, authPayload.User.Email)
		require.Equal(t, testUser.Username, authPayload.User.Username)
		require.NotEmpty(t, authPayload.Jwt)
		assertRefreshTokenCookie(t, recorder)
	})
}

func TestAuthResolver_Refresh_ReusedRefreshCookieFails(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		mutationResolver, _ := setupAuthMutationResolver(tx)
		ctx := context.Background()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		reusedTokenValue := "resolver-refresh-reused-token"
		familyID := uuid.New()
		seedRefreshTokenForResolverTest(t, ctx, tx, testUser.ID, familyID, reusedTokenValue)

		firstReq := httptest.NewRequest(http.MethodPost, "/query", nil)
		firstReq.AddCookie(&http.Cookie{Name: refreshTokenCookieName, Value: reusedTokenValue})
		firstRecorder := httptest.NewRecorder()

		firstPayload, err := mutationResolver.Refresh(authContextWithRequest(firstReq, firstRecorder))
		require.NoError(t, err)
		require.NotNil(t, firstPayload)

		secondReq := httptest.NewRequest(http.MethodPost, "/query", nil)
		secondReq.AddCookie(&http.Cookie{Name: refreshTokenCookieName, Value: reusedTokenValue})
		secondRecorder := httptest.NewRecorder()

		secondPayload, err := mutationResolver.Refresh(authContextWithRequest(secondReq, secondRecorder))

		require.Nil(t, secondPayload)
		assertUnauthenticatedGraphQLError(t, err, "Invalid refresh token")
		assertRefreshTokenCookieCleared(t, secondRecorder)
	})
}

func TestAuthResolver_Signout_NoRequestInContext(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		mutationResolver, _ := setupAuthMutationResolver(tx)

		signedOut, err := mutationResolver.Signout(authContext(nil))

		require.Error(t, err)
		require.False(t, signedOut)
		require.EqualError(t, err, "no request in context")
	})
}

func TestAuthResolver_Signout_NoRefreshCookie(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		mutationResolver, _ := setupAuthMutationResolver(tx)

		req := httptest.NewRequest(http.MethodPost, "/query", nil)
		recorder := httptest.NewRecorder()

		signedOut, err := mutationResolver.Signout(authContextWithRequest(req, recorder))

		require.NoError(t, err)
		require.True(t, signedOut)
		assertRefreshTokenCookieCleared(t, recorder)
	})
}

func TestAuthResolver_Signout_EmptyRefreshCookie(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		mutationResolver, _ := setupAuthMutationResolver(tx)

		req := httptest.NewRequest(http.MethodPost, "/query", nil)
		req.AddCookie(&http.Cookie{Name: refreshTokenCookieName, Value: ""})
		recorder := httptest.NewRecorder()

		signedOut, err := mutationResolver.Signout(authContextWithRequest(req, recorder))

		require.NoError(t, err)
		require.True(t, signedOut)
		assertRefreshTokenCookieCleared(t, recorder)
	})
}

func TestAuthResolver_Signout_InvalidRefreshCookie_ReturnsTrueAndClearsCookie(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		mutationResolver, _ := setupAuthMutationResolver(tx)

		req := httptest.NewRequest(http.MethodPost, "/query", nil)
		req.AddCookie(&http.Cookie{Name: refreshTokenCookieName, Value: "missing-refresh-token"})
		recorder := httptest.NewRecorder()

		signedOut, err := mutationResolver.Signout(authContextWithRequest(req, recorder))

		require.NoError(t, err)
		require.True(t, signedOut)
		assertRefreshTokenCookieCleared(t, recorder)
	})
}

func TestAuthResolver_Signout_Success_ClearsCookieAndRevokesToken(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		mutationResolver, authService := setupAuthMutationResolver(tx)
		ctx := context.Background()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		tokenValue := "resolver-signout-success-token"
		seedRefreshTokenForResolverTest(t, ctx, tx, testUser.ID, uuid.New(), tokenValue)

		req := httptest.NewRequest(http.MethodPost, "/query", nil)
		req.AddCookie(&http.Cookie{Name: refreshTokenCookieName, Value: tokenValue})
		recorder := httptest.NewRecorder()

		signedOut, err := mutationResolver.Signout(authContextWithRequest(req, recorder))

		require.NoError(t, err)
		require.True(t, signedOut)
		assertRefreshTokenCookieCleared(t, recorder)

		user, jwt, refreshedToken, err := authService.Refresh(ctx, tokenValue, "10.0.0.1")
		require.ErrorIs(t, err, refreshtokens.ErrInvalidRefreshToken)
		require.Nil(t, user)
		require.Empty(t, jwt)
		require.Nil(t, refreshedToken)
	})
}
