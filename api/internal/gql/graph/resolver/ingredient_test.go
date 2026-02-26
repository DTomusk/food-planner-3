package resolver

import (
	"context"
	"database/sql"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/testutil"
	"foodplanner/internal/testutil/seeds"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIngredientResolver_GetAllIngredients(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		txRunner := testutil.NewTestTxRunner(tx)

		ctx := context.Background()
		testIngredient, err := seeds.SeedTestIngredient(ctx, tx)
		require.NoError(t, err, "Failed to seed test ingredient")

		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)

		r := &Resolver{
			IngredientsService: ingredientService,
		}
		queryResolver := &queryResolver{r}

		ingredients, err := queryResolver.Ingredients(ctx)
		require.NoError(t, err, "GetAllIngredients failed")
		require.Len(t, ingredients, 1)

		ingredientModel := ingredients[0]
		require.Equal(t, testIngredient.ID.String(), ingredientModel.ID)
		require.Equal(t, testIngredient.Name, ingredientModel.Name)
		//require.Equal(t, testIngredient.PreferredUnit, ingredientModel.PreferredUnit)
	})
}

func TestIngredientResolver_GetAllIngredients_NoIngredients(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		txRunner := testutil.NewTestTxRunner(tx)

		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)

		r := &Resolver{
			IngredientsService: ingredientService,
		}
		queryResolver := &queryResolver{r}

		ctx := context.Background()
		ingredients, err := queryResolver.Ingredients(ctx)
		require.NoError(t, err, "GetAllIngredients failed")
		require.Len(t, ingredients, 0)
	})
}
