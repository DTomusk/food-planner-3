package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"foodplanner/internal/audit"
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
		auditService := audit.NewAuditService(tx, audit.NewRepo())
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService, auditService)

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

func TestSignUp_PersistsAuditEvent(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		refreshTokenService := refreshtokens.NewRefreshTokenService(txRunner, refreshtokens.NewRefreshTokenRepo(), "refresh-secret", 7)
		auditService := audit.NewAuditService(tx, audit.NewRepo())
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService, auditService)

		email := "signup-audit@test.com"
		password := "securepassword"
		username := "audited-user"
		ipAddress := "127.0.0.1"

		// Act
		signedUpUser, _, _, err := authService.SignUp(email, password, username, ipAddress, context.Background())
		require.NoError(t, err)

		var (
			storedActorID    string
			storedResourceID string
			storedAction     string
			storedResult     string
			storedNewState   []byte
			storedContext    []byte
		)

		err = tx.QueryRowContext(context.Background(), `
			SELECT actor_id::text, resource_id::text, action, result, new_state, context
			FROM audits
			WHERE action = $1 AND resource_id = $2
			ORDER BY created_at DESC
			LIMIT 1
		`, audit.ActionUserSignup, signedUpUser.ID).Scan(
			&storedActorID,
			&storedResourceID,
			&storedAction,
			&storedResult,
			&storedNewState,
			&storedContext,
		)
		require.NoError(t, err)

		require.Equal(t, signedUpUser.ID.String(), storedActorID)
		require.Equal(t, signedUpUser.ID.String(), storedResourceID)
		require.Equal(t, string(audit.ActionUserSignup), storedAction)
		require.Equal(t, string(audit.ResultSuccess), storedResult)

		var newState struct {
			UserID   string `json:"user_id"`
			Username string `json:"username"`
		}
		err = json.Unmarshal(storedNewState, &newState)
		require.NoError(t, err)
		require.Equal(t, signedUpUser.ID.String(), newState.UserID)
		require.Equal(t, username, newState.Username)

		var contextData struct {
			Source    string `json:"source"`
			Operation string `json:"operation"`
			IPAddress string `json:"ip_address"`
		}
		err = json.Unmarshal(storedContext, &contextData)
		require.NoError(t, err)
		require.Equal(t, "graphql", contextData.Source)
		require.Equal(t, "signup", contextData.Operation)
		require.Equal(t, ipAddress, contextData.IPAddress)
	})
}

func TestSignup_InvalidEmail(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		refreshTokenService := refreshtokens.NewRefreshTokenService(txRunner, refreshtokens.NewRefreshTokenRepo(), "refresh-secret", 7)
		auditService := audit.NewAuditService(tx, audit.NewRepo())
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService, auditService)
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
		auditService := audit.NewAuditService(tx, audit.NewRepo())
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService, auditService)
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
		auditService := audit.NewAuditService(tx, audit.NewRepo())
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService, auditService)
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
		auditService := audit.NewAuditService(tx, audit.NewRepo())
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService, auditService)

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
		auditService := audit.NewAuditService(tx, audit.NewRepo())
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService, auditService)
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
		auditService := audit.NewAuditService(tx, audit.NewRepo())
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService, auditService)

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
		auditService := audit.NewAuditService(tx, audit.NewRepo())
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService, auditService)
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

func TestRefresh_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		refreshTokenService := refreshtokens.NewRefreshTokenService(txRunner, refreshtokens.NewRefreshTokenRepo(), "refresh-secret", 7)
		auditService := audit.NewAuditService(tx, audit.NewRepo())
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService, auditService)

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
		txRunner := testutil.NewTestTxRunner(tx)
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		refreshTokenService := refreshtokens.NewRefreshTokenService(txRunner, refreshtokens.NewRefreshTokenRepo(), "refresh-secret", 7)
		auditService := audit.NewAuditService(tx, audit.NewRepo())
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService, auditService)

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
		txRunner := testutil.NewTestTxRunner(tx)
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		refreshTokenService := refreshtokens.NewRefreshTokenService(txRunner, refreshtokens.NewRefreshTokenRepo(), "refresh-secret", 7)
		auditService := audit.NewAuditService(tx, audit.NewRepo())
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService, auditService)

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
		txRunner := testutil.NewTestTxRunner(tx)
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		refreshTokenService := refreshtokens.NewRefreshTokenService(txRunner, refreshtokens.NewRefreshTokenRepo(), "refresh-secret", 7)
		auditService := audit.NewAuditService(tx, audit.NewRepo())
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService, auditService)

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
		txRunner := testutil.NewTestTxRunner(tx)
		userService := user.NewUserService(tx, user.NewUserRepo())
		jwtService := NewJWTService("testsecret", 15)
		refreshTokenService := refreshtokens.NewRefreshTokenService(txRunner, refreshtokens.NewRefreshTokenRepo(), "refresh-secret", 7)
		auditService := audit.NewAuditService(tx, audit.NewRepo())
		authService := NewAuthService(tx, userService, jwtService, refreshTokenService, auditService)

		// Act
		err := authService.Revoke(context.Background(), "missing-refresh-token")

		// Assert
		require.NoError(t, err)
	})
}
