package recipe

import (
	"context"
	"database/sql"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/testutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateRecipe(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		txRunner := testutil.NewTestTxRunner(tx)
		s := NewService(txRunner, NewRepo(), ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100))
		request := CreateRecipeRequest{Name: "Vanilla Ice Cream"}
		recipe, err := s.CreateRecipe(context.Background(), request)
		require.NoError(t, err)
		require.Equal(t, "Vanilla Ice Cream", recipe.Name)
	})
}
