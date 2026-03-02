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
	require.Len(t, ingredients, 3)

	testIngredient := *ingredients[0]
	require.Equal(t, testIngredient.ID, "test_ingredient")
	require.Equal(t, testIngredient.Name, "Test Ingredient")
	require.Equal(t, testIngredient.PreferredUnit, 1)
	require.Equal(t, testIngredient.Plural, testutil.PtrString("Test Ingredients"))
	require.Nil(t, testIngredient.Counter, "Expected Counter to be nil for test_ingredient")
	require.Nil(t, testIngredient.CounterPlural, "Expected CounterPlural to be nil for test_ingredient")

	testIngredientGrams := *ingredients[1]
	require.Equal(t, testIngredientGrams.ID, "test_ingredient_grams")
	require.Equal(t, testIngredientGrams.Name, "Test Ingredient (grams)")
	require.Equal(t, testIngredientGrams.PreferredUnit, 2)
	require.Nil(t, testIngredientGrams.Plural)
	require.Nil(t, testIngredientGrams.Counter, "Expected Counter to be nil for test_ingredient_grams")
	require.Nil(t, testIngredientGrams.CounterPlural, "Expected CounterPlural to be nil for test_ingredient_grams")

	testIngredientWithCounter := *ingredients[2]
	require.Equal(t, testIngredientWithCounter.ID, "test_ingredient_with_counter")
	require.Equal(t, testIngredientWithCounter.Name, "Test Ingredient with Counter")
	require.Equal(t, testIngredientWithCounter.PreferredUnit, 1)
	require.Nil(t, testIngredientWithCounter.Plural)
	require.Equal(t, testIngredientWithCounter.Counter, testutil.PtrString("piece"))
	require.Equal(t, testIngredientWithCounter.CounterPlural, testutil.PtrString("pieces"))
}
