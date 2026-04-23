package ingredient

import (
	"testing"

	"foodplanner/internal/unit"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewIngredient_DefaultsProcessedLevelAndSearchable(t *testing.T) {
	ing, err := NewIngredient(
		"Carrot",
		"carrot",
		int(unit.Gram),
		nil,
		nil,
		nil,
		Vegan,
		false,
		nil,
		nil,
		nil,
	)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, ing.ID)
	require.Equal(t, Raw, ing.ProcessedLevel)
	require.True(t, ing.IsSearchable)
}

func TestNewIngredient_UsesProvidedProcessedLevelAndSearchable(t *testing.T) {
	processedLevel := int(Derived)
	isSearchable := false

	ing, err := NewIngredient(
		"Garlic Powder",
		"garlic_powder",
		int(unit.UnitUnknown),
		nil,
		nil,
		nil,
		Vegan,
		false,
		nil,
		&processedLevel,
		&isSearchable,
	)
	require.NoError(t, err)
	require.Equal(t, Derived, ing.ProcessedLevel)
	require.False(t, ing.IsSearchable)
	require.Equal(t, unit.UnitUnknown, ing.PreferredUnit)
}

func TestNewIngredient_ReturnsErrInvalidName(t *testing.T) {
	ing, err := NewIngredient(
		"",
		"empty_name",
		int(unit.Gram),
		nil,
		nil,
		nil,
		Vegan,
		false,
		nil,
		nil,
		nil,
	)
	require.Nil(t, ing)
	require.ErrorIs(t, err, ErrInvalidName)
}

func TestNewIngredient_ReturnsErrInvalidPreferredUnit_WhenSearchable(t *testing.T) {
	isSearchable := true

	ing, err := NewIngredient(
		"Category Node",
		"category_node",
		int(unit.UnitUnknown),
		nil,
		nil,
		nil,
		Vegan,
		false,
		nil,
		nil,
		&isSearchable,
	)
	require.Nil(t, ing)
	require.ErrorIs(t, err, ErrInvalidPreferredUnit)
}
