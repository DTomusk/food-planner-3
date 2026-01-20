package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword_Success(t *testing.T) {
	password := "SecurePass123!"
	hash, err := hashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	require.NotEqual(t, password, hash)
}

func TestComparePasswordHash_Success(t *testing.T) {
	password := "SecurePass123!"
	hash, err := hashPassword(password)
	require.NoError(t, err)
	err = comparePasswordHash(password, hash)
	require.NoError(t, err)
}

func TestComparePasswordHash_IncorrectPassword(t *testing.T) {
	password := "SecurePass123!"
	wrongPassword := "WrongPass456!"
	hash, err := hashPassword(password)
	require.NoError(t, err)
	err = comparePasswordHash(wrongPassword, hash)
	require.Error(t, err)
}
