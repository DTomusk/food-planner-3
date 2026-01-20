package user

import (
	"context"
	"database/sql"
	"foodplanner/internal/testutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetByEmail_DoesntThrow(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		repo := NewUserRepo()
		_, err := repo.GetUserByEmail("test@mail.com", context.Background(), tx)
		require.NoError(t, err)
	})
}

func TestGetByID_Throws(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		repo := NewUserRepo()
		_, err := repo.GetUserByID("non-existent-id", context.Background(), tx)
		require.Error(t, err)
	})
}

func TestCreate_ReturnsUser(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		repo := NewUserRepo()
		user := NewUser("blah@test.com", "securepassword")
		repoUser, err := repo.CreateUser(user, context.Background(), tx)
		require.NoError(t, err)
		require.Equal(t, user.ID, repoUser.ID)
		require.Equal(t, user.Email, repoUser.Email)
		require.Equal(t, user.PasswordHash, repoUser.PasswordHash)
	})
}
