package recipe

import (
	"context"
	"database/sql"
	"food-planner/internal/testutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateRecipe(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		s := NewService(tx, NewRepo())
		recipe, err := s.CreateRecipe("Vanilla Ice Cream", context.Background())
		require.NoError(t, err)
		require.Equal(t, "Vanilla Ice Cream", recipe.Name)
	})
}
