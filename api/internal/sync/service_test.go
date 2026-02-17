package sync

import (
	"database/sql"
	"foodplanner/internal/ingredient"
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
		ingredients, err := ingredientService.GetAllIngredients(ctx)
		require.NoError(t, err)
		require.Len(t, ingredients, 1)

		// Copy assertions from loader test
		testIngredient := *ingredients[0]
		require.Equal(t, testIngredient.FileKey, "test_ingredient")
		require.Equal(t, testIngredient.Name, "Test Ingredient")
		require.Equal(t, testIngredient.PreferredUnit, unit.Quantum)
	})
}
