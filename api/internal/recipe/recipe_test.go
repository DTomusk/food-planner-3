package recipe

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestInstantiateRecipe(t *testing.T) {
	name := "Pancakes"
	userID := uuid.New()
	recipe, err := NewRecipe(name, userID, nil)
	require.NoError(t, err)
	require.Equal(t, name, recipe.Name)
	require.Equal(t, userID, recipe.UserID)
}

func TestEmptyRecipeName(t *testing.T) {
	userID := uuid.New()
	_, err := NewRecipe("", userID, nil)
	require.Error(t, err)
	require.Equal(t, ErrEmptyName, err)
}
