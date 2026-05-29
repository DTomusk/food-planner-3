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

// TODO: these tests are brittle because they depend on sql order, which isn't guaranteed
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

func TestSyncIngredientData_Idempotent(t *testing.T) {
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
		ingredients_first_fetch, err := ingredientService.GetAllIngredients(ctx, logger)
		require.NoError(t, err)
		require.Len(t, ingredients_first_fetch, 7)

		// Act again
		err = service.SyncIngredientData(ctx)

		// Assert no error and no duplicates
		require.NoError(t, err)

		ingredients_second_fetch, err := ingredientService.GetAllIngredients(ctx, logger)
		require.NoError(t, err)
		require.Len(t, ingredients_second_fetch, 7)

		for i := range ingredients_first_fetch {
			require.Equal(t, ingredients_first_fetch[i].ID, ingredients_second_fetch[i].ID, "Expected ingredient IDs to be the same after second sync")
			require.Equal(t, ingredients_first_fetch[i].FileKey, ingredients_second_fetch[i].FileKey, "Expected ingredient FileKeys to be the same after second sync")
			require.Equal(t, ingredients_first_fetch[i].Name, ingredients_second_fetch[i].Name, "Expected ingredient Names to be the same after second sync")
			require.Equal(t, ingredients_first_fetch[i].PreferredUnit, ingredients_second_fetch[i].PreferredUnit, "Expected ingredient PreferredUnits to be the same after second sync")
			require.Equal(t, ingredients_first_fetch[i].Plural, ingredients_second_fetch[i].Plural, "Expected ingredient Plurals to be the same after second sync")
			require.Equal(t, ingredients_first_fetch[i].Counter, ingredients_second_fetch[i].Counter, "Expected ingredient Counters to be the same after second sync")
			require.Equal(t, ingredients_first_fetch[i].CounterPlural, ingredients_second_fetch[i].CounterPlural, "Expected ingredient CounterPlurals to be the same after second sync")
			require.Equal(t, ingredients_first_fetch[i].TaxonomyParentID, ingredients_second_fetch[i].TaxonomyParentID, "Expected ingredient TaxonomyParentIDs to be the same after second sync")
		}
	})
}

func TestSyncIngredientData_Idempotent_WithNonSearchableParentChild(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := t.Context()
		logger := logging.FromContext(ctx)
		filePath := "../../reference/ingredients_non_searchable_parent_test.yaml"
		loader := reference.NewLoader(filePath)
		txRunner := testutil.NewTestTxRunner(tx)
		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)
		service := NewSyncService(ingredientService, loader)

		err := service.SyncIngredientData(ctx)
		require.NoError(t, err)

		err = service.SyncIngredientData(ctx)
		require.NoError(t, err)

		ingredients, err := ingredientService.GetAllIngredientsUnfiltered(ctx, logger)
		require.NoError(t, err)
		require.Len(t, ingredients, 2)

		var parent, child *ingredient.Ingredient
		for _, synced := range ingredients {
			switch synced.FileKey {
			case "test_parent_non_searchable":
				parent = synced
			case "test_child_non_searchable":
				child = synced
			}
		}

		require.NotNil(t, parent)
		require.NotNil(t, child)
		require.NotNil(t, child.TaxonomyParentID)
		require.Equal(t, parent.ID, *child.TaxonomyParentID)
	})
}
