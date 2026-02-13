package reference

import (
	"foodplanner/internal/logging"
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
	require.Len(t, ingredients, 1)

	testIngredient := *ingredients[0]
	require.Equal(t, testIngredient.ID, "test_ingredient")
	require.Equal(t, testIngredient.Name, "Test Ingredient")
	require.Equal(t, testIngredient.PreferredUnit, 1)
}
