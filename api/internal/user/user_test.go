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
	user, err := NewUser(email, passwordHash, username)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, email, user.Email)
	require.Equal(t, passwordHash, user.PasswordHash)
	require.Equal(t, username, user.Username)
}

func TestNewUser_UsernameRequired(t *testing.T) {
	// Arrange
	email := "a@b.com"
	passwordHash := "hashedpassword"
	username := ""

	// Act
	user, err := NewUser(email, passwordHash, username)

	// Assert
	require.Error(t, err)
	require.Nil(t, user)
	require.Equal(t, ErrUsernameRequired, err)
}

func TestNewUser_UsernameTooLong(t *testing.T) {
	// Arrange
	email := "g@net.org"
	passwordHash := "hashedpassword"
	username := "thisisaverylongusernamethatisdefinitelymorethanfiftycharacterslongandshouldcauseanerror"

	// Act
	user, err := NewUser(email, passwordHash, username)

	// Assert
	require.Error(t, err)
	require.Nil(t, user)
	require.Equal(t, ErrUsernameTooLong, err)
}
