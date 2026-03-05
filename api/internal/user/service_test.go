package user

import (
	"context"
	"database/sql"
	"foodplanner/internal/testutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser_Success(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		s := NewUserService(tx, NewUserRepo())
		user, err := s.CreateUser("test@example.com", "hashedpassword", "username", context.Background())
		require.NoError(t, err)
		require.Equal(t, "test@example.com", user.Email)
		require.NotEmpty(t, user.ID)
		require.Equal(t, "hashedpassword", user.PasswordHash)
		require.Equal(t, "username", user.Username)
	})
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		s := NewUserService(tx, NewUserRepo())
		_, err := s.CreateUser("test@example.com", "hashedpassword", "username", context.Background())
		require.NoError(t, err)
		_, err = s.CreateUser("test@example.com", "doesntmatter", "differentusername", context.Background())
		require.Error(t, err)
		require.Equal(t, ErrEmailInUse, err)
	})
}

func TestCreateUser_DuplicateUsername(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		s := NewUserService(tx, NewUserRepo())
		_, err := s.CreateUser("test1@example.com", "hashedpassword", "username", context.Background())
		require.NoError(t, err)
		_, err = s.CreateUser("test2@example.com", "doesntmatter", "username", context.Background())
		// Returns a sql error
		require.Error(t, err)
	})
}
