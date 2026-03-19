package user

import (
	"context"
	"database/sql"
	"foodplanner/internal/testutil"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetByEmail_DoesntThrow(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		repo := NewUserRepo()
		_, err := repo.getUserByEmail("test@mail.com", context.Background(), tx)
		require.NoError(t, err)
	})
}

func TestGetByID_Throws(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		repo := NewUserRepo()
		_, err := repo.getUserByID(context.Background(), tx, uuid.New())
		require.Error(t, err)
	})
}

func TestCreate_ReturnsUser(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		repo := NewUserRepo()
		user, err := NewUser("blah@test.com", "securepassword", "testuser")
		require.NoError(t, err)
		repoUser, err := repo.createUser(user, context.Background(), tx)
		require.NoError(t, err)
		require.Equal(t, user.ID, repoUser.ID)
		require.Equal(t, user.Email, repoUser.Email)
		require.Equal(t, user.PasswordHash, repoUser.PasswordHash)
	})
}
