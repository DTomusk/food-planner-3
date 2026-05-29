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
		require.Len(t, ingredients, 7)

		testIngredient := *ingredients[0]
		require.Equal(t, testIngredient.FileKey, "test_ingredient")
		require.Equal(t, testIngredient.Name, "Test Ingredient")
		require.Equal(t, testIngredient.PreferredUnit, unit.Quantum)
		require.Equal(t, testIngredient.Plural, testutil.PtrString("Test Ingredients"))
		require.Nil(t, testIngredient.Counter, "Expected Counter to be nil for test_ingredient")
		require.Nil(t, testIngredient.CounterPlural, "Expected CounterPlural to be nil for test_ingredient")
		require.Nil(t, testIngredient.TaxonomyParentID, "Expected TaxonomyParentID to be nil for test_ingredient")

		testIngredientGrams := *ingredients[1]
		require.Equal(t, testIngredientGrams.FileKey, "test_ingredient_grams")
		require.Equal(t, testIngredientGrams.Name, "Test Ingredient (grams)")
		require.Equal(t, testIngredientGrams.PreferredUnit, unit.Gram)
		require.Nil(t, testIngredientGrams.Plural)
		require.Nil(t, testIngredientGrams.Counter, "Expected Counter to be nil for test_ingredient_grams")
		require.Nil(t, testIngredientGrams.CounterPlural, "Expected CounterPlural to be nil for test_ingredient_grams")
		require.Nil(t, testIngredientGrams.TaxonomyParentID, "Expected TaxonomyParentID to be nil for test_ingredient_grams")

		testIngredientWithCounter := *ingredients[2]
		require.Equal(t, testIngredientWithCounter.FileKey, "test_ingredient_with_counter")
		require.Equal(t, testIngredientWithCounter.Name, "Test Ingredient with Counter")
		require.Equal(t, testIngredientWithCounter.PreferredUnit, unit.Quantum)
		require.Nil(t, testIngredientWithCounter.Plural)
		require.Equal(t, testIngredientWithCounter.Counter, testutil.PtrString("piece"))
		require.Equal(t, testIngredientWithCounter.CounterPlural, testutil.PtrString("pieces"))
		require.Nil(t, testIngredientWithCounter.TaxonomyParentID, "Expected TaxonomyParentID to be nil for test_ingredient_with_counter")

		testIngredientParent := *ingredients[3]
		require.Equal(t, testIngredientParent.FileKey, "test_ingredient_parent")
		require.Equal(t, testIngredientParent.Name, "Test Ingredient Parent")
		require.Equal(t, testIngredientParent.PreferredUnit, unit.Quantum)
		require.Equal(t, testIngredientParent.Plural, testutil.PtrString("Test Ingredient Parents"))
		require.Nil(t, testIngredientParent.Counter)
		require.Nil(t, testIngredientParent.CounterPlural)
		require.Nil(t, testIngredientParent.TaxonomyParentID, "Expected TaxonomyParentID to be nil for test_ingredient_parent")

		testIngredientChild := *ingredients[4]
		require.Equal(t, testIngredientChild.FileKey, "test_ingredient_child")
		require.Equal(t, testIngredientChild.Name, "Test Ingredient Child")
		require.Equal(t, testIngredientChild.PreferredUnit, unit.Quantum)
		require.Equal(t, testIngredientChild.Plural, testutil.PtrString("Test Ingredient Children"))
		require.Equal(t, *testIngredientChild.TaxonomyParentID, testIngredientParent.ID, "Expected TaxonomyParentID to be set to parent ID for test_ingredient_child")

		testIngredientGrandchild := *ingredients[5]
		require.Equal(t, testIngredientGrandchild.FileKey, "test_ingredient_grandchild")
		require.Equal(t, testIngredientGrandchild.Name, "Test Ingredient Grandchild")
		require.Equal(t, testIngredientGrandchild.PreferredUnit, unit.Quantum)
		require.Equal(t, testIngredientGrandchild.Plural, testutil.PtrString("Test Ingredient Grandchildren"))
		require.Equal(t, *testIngredientGrandchild.TaxonomyParentID, testIngredientChild.ID, "Expected TaxonomyParentID to be set to child ID for test_ingredient_grandchild")

		testIngredientChild2 := *ingredients[6]
		require.Equal(t, testIngredientChild2.FileKey, "test_ingredient_child_2")
		require.Equal(t, testIngredientChild2.Name, "Test Ingredient Child 2")
		require.Equal(t, testIngredientChild2.PreferredUnit, unit.Quantum)
		require.Equal(t, testIngredientChild2.Plural, testutil.PtrString("Test Ingredient Children 2"))
		require.Equal(t, *testIngredientChild2.TaxonomyParentID, testIngredientParent.ID, "Expected TaxonomyParentID to be set to parent ID for test_ingredient_child_2")
	})
}
