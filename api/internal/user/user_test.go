package user

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	// Arrange
	email := "a@b.com"
	passwordHash := "hashedpassword"
	username := "testuser"

	// Act
	user := NewUser(email, passwordHash, username)

	// Assert
	require.NotNil(t, user)
	require.Equal(t, email, user.Email)
	require.Equal(t, passwordHash, user.PasswordHash)
	require.Equal(t, username, user.Username)
}
