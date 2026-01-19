package recipe

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInstantiateRecipe(t *testing.T) {
	name := "Pancakes"
	recipe, err := NewRecipe(name)
	require.NoError(t, err)
	require.Equal(t, name, recipe.Name)
}

func TestEmptyRecipeName(t *testing.T) {
	_, err := NewRecipe("")
	require.Error(t, err)
	require.Equal(t, ErrEmptyName, err)
}
