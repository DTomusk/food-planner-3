package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMiddleware_NoAuthHeader(t *testing.T) {
	// Arrange
	jwtService := NewJWTService("testsecret", 15)
	middleware := Middleware(jwtService)
	handlerCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// Act
	middleware(nextHandler).ServeHTTP(rr, req)

	// Assert
	require.True(t, handlerCalled, "next handler should be called when no Authorization header is present")
}

func TestMiddleware_InvalidAuthHeaderFormat(t *testing.T) {
	// Arrange
	jwtService := NewJWTService("testsecret", 15)
	middleware := Middleware(jwtService)
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	rr := httptest.NewRecorder()

	// Act
	middleware(nextHandler).ServeHTTP(rr, req)

	// Assert
	require.Equal(t, http.StatusUnauthorized, rr.Code)
	require.Equal(t, "invalid Authorization header format\n", rr.Body.String())
}

func TestMiddleware_ValidToken(t *testing.T) {
	// Arrange
	jwtService := NewJWTService("testsecret", 15)
	middleware := Middleware(jwtService)
	userID := "user-123"
	token, err := jwtService.GenerateToken(userID)
	require.NoError(t, err)
	handlerCalled := false
	var claimsInContext *Claims
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		claimsInContext, err = ClaimsFromContext(r.Context())
	})
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	// Act
	middleware(nextHandler).ServeHTTP(rr, req)

	// Assert
	require.True(t, handlerCalled, "next handler should be called with valid token")
	require.NoError(t, err)
	require.NotNil(t, claimsInContext)
	require.Equal(t, userID, claimsInContext.UserID)
}
