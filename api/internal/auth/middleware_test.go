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
	jwtService := NewJWTService("testsecret", 15)
	middleware := Middleware(jwtService)

	handlerCalled := false
	var claimsInContext *Claims

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		var ok bool
		claimsInContext, ok = ClaimsFromContext(r.Context())
		require.False(t, ok, "claims should not be present in context for invalid auth header format")
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	rr := httptest.NewRecorder()

	middleware(nextHandler).ServeHTTP(rr, req)

	require.True(t, handlerCalled, "next handler should always be called")
	require.Nil(t, claimsInContext)
}

func TestMiddleware_InvalidToken(t *testing.T) {
	jwtService := NewJWTService("testsecret", 15)
	middleware := Middleware(jwtService)

	handlerCalled := false
	var claimsInContext *Claims

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		var ok bool
		claimsInContext, ok = ClaimsFromContext(r.Context())
		require.False(t, ok, "claims should not be present in context for invalid token")
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rr := httptest.NewRecorder()

	middleware(nextHandler).ServeHTTP(rr, req)

	require.True(t, handlerCalled)
	require.Nil(t, claimsInContext)
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

		claims, ok := ClaimsFromContext(r.Context())
		require.True(t, ok, "claims should be present in context for valid token")

		claimsInContext = claims
	})
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	// Act
	middleware(nextHandler).ServeHTTP(rr, req)

	// Assert
	require.True(t, handlerCalled, "next handler should be called with valid token")
	require.NotNil(t, claimsInContext)
	require.Equal(t, userID, claimsInContext.UserID)
}
