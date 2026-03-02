package sync

import (
	"database/sql"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/logging"
	"foodplanner/internal/reference"
	"foodplanner/internal/testutil"
	"foodplanner/internal/unit"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSyncIngredientData(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		ctx := t.Context()
		logger := logging.FromContext(ctx)
		filePath := "../../reference/ingredients_test.yaml"
		loader := reference.NewLoader(filePath)
		txRunner := testutil.NewTestTxRunner(tx)
		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)
		service := NewSyncService(ingredientService, loader)

		// Act
		err := service.SyncIngredientData(ctx)

		// Assert
		require.NoError(t, err)

		// Fetch ingredient via ingredient service to ensure it's populate
		ingredients, err := ingredientService.GetAllIngredients(ctx, logger)
		require.NoError(t, err)
		require.Len(t, ingredients, 3)

		// Copy assertions from loader test
		testIngredient := *ingredients[0]
		require.Equal(t, testIngredient.FileKey, "test_ingredient")
		require.Equal(t, testIngredient.Name, "Test Ingredient")
		require.Equal(t, testIngredient.PreferredUnit, unit.Quantum)
		require.Equal(t, testIngredient.Plural, testutil.PtrString("Test Ingredients"))
		require.Nil(t, testIngredient.Counter, "Expected Counter to be nil for test_ingredient")
		require.Nil(t, testIngredient.CounterPlural, "Expected CounterPlural to be nil for test_ingredient")

		testIngredientGrams := *ingredients[1]
		require.Equal(t, testIngredientGrams.FileKey, "test_ingredient_grams")
		require.Equal(t, testIngredientGrams.Name, "Test Ingredient (grams)")
		require.Equal(t, testIngredientGrams.PreferredUnit, unit.Gram)
		require.Nil(t, testIngredientGrams.Plural)
		require.Nil(t, testIngredientGrams.Counter, "Expected Counter to be nil for test_ingredient_grams")
		require.Nil(t, testIngredientGrams.CounterPlural, "Expected CounterPlural to be nil for test_ingredient_grams")

		testIngredientWithCounter := *ingredients[2]
		require.Equal(t, testIngredientWithCounter.FileKey, "test_ingredient_with_counter")
		require.Equal(t, testIngredientWithCounter.Name, "Test Ingredient with Counter")
		require.Equal(t, testIngredientWithCounter.PreferredUnit, unit.Quantum)
		require.Nil(t, testIngredientWithCounter.Plural)
		require.Equal(t, testIngredientWithCounter.Counter, testutil.PtrString("piece"))
		require.Equal(t, testIngredientWithCounter.CounterPlural, testutil.PtrString("pieces"))
	})
}
