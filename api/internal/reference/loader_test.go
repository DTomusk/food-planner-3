package reference

import (
	"foodplanner/internal/logging"
	"foodplanner/internal/testutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadTestData(t *testing.T) {
	// Arrange
	filePath := "../../reference/ingredients_test.yaml"
	loader := NewLoader(filePath)
	logger := logging.FromContext(t.Context())

	// Act
	ingredients, err := loader.LoadIngredientData(logger)

	// Assert
	require.NoError(t, err)

	// Note: the length is of the parent ingredients (tree roots), not counting others
	require.Len(t, ingredients, 4)

	testIngredient := *ingredients[0]
	require.Equal(t, testIngredient.FileKey, "test_ingredient")
	require.Equal(t, testIngredient.Name, "Test Ingredient")
	require.Equal(t, testIngredient.PreferredUnit, 1)
	require.Equal(t, testIngredient.Plural, testutil.PtrString("Test Ingredients"))
	require.Nil(t, testIngredient.Counter, "Expected Counter to be nil for test_ingredient")
	require.Nil(t, testIngredient.CounterPlural, "Expected CounterPlural to be nil for test_ingredient")
	require.Nil(t, testIngredient.Children, "Expected Children to be nil for test_ingredient")

	testIngredientGrams := *ingredients[1]
	require.Equal(t, testIngredientGrams.FileKey, "test_ingredient_grams")
	require.Equal(t, testIngredientGrams.Name, "Test Ingredient (grams)")
	require.Equal(t, testIngredientGrams.PreferredUnit, 2)
	require.Nil(t, testIngredientGrams.Plural)
	require.Nil(t, testIngredientGrams.Counter, "Expected Counter to be nil for test_ingredient_grams")
	require.Nil(t, testIngredientGrams.CounterPlural, "Expected CounterPlural to be nil for test_ingredient_grams")
	require.Nil(t, testIngredientGrams.Children, "Expected Children to be nil for test_ingredient_grams")

	testIngredientWithCounter := *ingredients[2]
	require.Equal(t, testIngredientWithCounter.FileKey, "test_ingredient_with_counter")
	require.Equal(t, testIngredientWithCounter.Name, "Test Ingredient with Counter")
	require.Equal(t, testIngredientWithCounter.PreferredUnit, 1)
	require.Nil(t, testIngredientWithCounter.Plural)
	require.Equal(t, testIngredientWithCounter.Counter, testutil.PtrString("piece"))
	require.Equal(t, testIngredientWithCounter.CounterPlural, testutil.PtrString("pieces"))
	require.Nil(t, testIngredientWithCounter.Children, "Expected Children to be nil for test_ingredient_with_counter")

	testIngredientParent := *ingredients[3]
	require.Equal(t, testIngredientParent.FileKey, "test_ingredient_parent")
	require.Equal(t, testIngredientParent.Name, "Test Ingredient Parent")
	require.Equal(t, testIngredientParent.PreferredUnit, 1)
	require.Equal(t, testIngredientParent.Plural, testutil.PtrString("Test Ingredient Parents"))
	require.Nil(t, testIngredientParent.Counter)
	require.Nil(t, testIngredientParent.CounterPlural)
	require.Len(t, testIngredientParent.Children, 2)

	testIngredientChild := *testIngredientParent.Children[0]
	require.Equal(t, testIngredientChild.FileKey, "test_ingredient_child")
	require.Equal(t, testIngredientChild.Name, "Test Ingredient Child")
	require.Equal(t, testIngredientChild.PreferredUnit, 1)
	require.Equal(t, testIngredientChild.Plural, testutil.PtrString("Test Ingredient Children"))
	require.Len(t, testIngredientChild.Children, 1)

	testIngredientGrandchild := *testIngredientChild.Children[0]
	require.Equal(t, testIngredientGrandchild.FileKey, "test_ingredient_grandchild")
	require.Equal(t, testIngredientGrandchild.Name, "Test Ingredient Grandchild")
	require.Equal(t, testIngredientGrandchild.PreferredUnit, 1)
	require.Equal(t, testIngredientGrandchild.Plural, testutil.PtrString("Test Ingredient Grandchildren"))
	require.Nil(t, testIngredientGrandchild.Children, "Expected Children to be nil for test_ingredient_grandchild")

	testIngredientChild2 := *testIngredientParent.Children[1]
	require.Equal(t, testIngredientChild2.FileKey, "test_ingredient_child_2")
	require.Equal(t, testIngredientChild2.Name, "Test Ingredient Child 2")
	require.Equal(t, testIngredientChild2.PreferredUnit, 1)
	require.Equal(t, testIngredientChild2.Plural, testutil.PtrString("Test Ingredient Children 2"))
	require.Nil(t, testIngredientChild2.Children, "Expected Children to be nil for test_ingredient_child_2")
}

// This test verifies that the loader can successfully load the actual reference file without errors and that it contains data.
func TestLoadingReferenceFile(t *testing.T) {
	// Arrange
	filePath := "../../reference/ingredients.yaml"
	loader := NewLoader(filePath)
	logger := logging.FromContext(t.Context())

	// Act
	ingredients, err := loader.LoadIngredientData(logger)

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, ingredients)
}
