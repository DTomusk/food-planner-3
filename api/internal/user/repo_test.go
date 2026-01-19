package users

import (
	"context"
	"database/sql"
	"food-planner/internal/testutil"
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
