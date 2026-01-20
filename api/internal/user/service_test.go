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
		user, err := s.CreateUser("test@example.com", "hashedpassword", context.Background())
		require.NoError(t, err)
		require.Equal(t, "test@example.com", user.Email)
		require.NotEmpty(t, user.ID)
		require.Equal(t, "hashedpassword", user.PasswordHash)
	})
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		s := NewUserService(tx, NewUserRepo())
		_, err := s.CreateUser("test@example.com", "hashedpassword", context.Background())
		require.NoError(t, err)
		_, err = s.CreateUser("test@example.com", "doesntmatter", context.Background())
		require.Error(t, err)
		require.Equal(t, ErrEmailInUse, err)
	})
}
