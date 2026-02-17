package recipe

import (
	"context"
	"database/sql"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/testutil"
	"foodplanner/internal/testutil/seeds"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateRecipe(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		txRunner := testutil.NewTestTxRunner(tx)

		ingredientID := uuid.New()
		testIngredient := ingredient.Ingredient{
			ID:            ingredientID,
			FileKey:       "test_ingredient",
			Name:          "Test Ingredient",
			PreferredUnit: 1,
		}
		err := seeds.InsertIngredient(context.Background(), tx, &testIngredient)
		require.NoError(t, err, "Failed to seed test ingredient")

		s := NewService(txRunner, NewRepo(), ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100))
		ingredientRequest := CreateIngredientUsageRequest{
			IngredientID: ingredientID.String(),
			Quantity:     200,
			Unit:         1,
		}
		request := CreateRecipeRequest{
			Name:        "Vanilla Ice Cream",
			Ingredients: []CreateIngredientUsageRequest{ingredientRequest},
		}
		recipe, err := s.CreateRecipe(context.Background(), request)
		require.NoError(t, err)
		require.Equal(t, "Vanilla Ice Cream", recipe.Name)
	})
}
