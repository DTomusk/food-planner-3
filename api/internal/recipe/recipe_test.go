package recipe

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestInstantiateRecipe(t *testing.T) {
	name := "Pancakes"
	userID := uuid.New()
	ingredientUsage := &IngredientUsage{
		IngredientID: uuid.New(),
		Quantity:     200,
		Unit:         1,
	}
	recipe, err := NewRecipe(name, userID, []*IngredientUsage{ingredientUsage}, 10, 20, 4)
	require.NoError(t, err)
	require.Equal(t, name, recipe.Name)
	require.Equal(t, userID, recipe.UserID)
	require.Equal(t, 10, recipe.PrepMins)
	require.Equal(t, 20, recipe.CookMins)
	require.Equal(t, 4, recipe.Portions)
}

func TestEmptyRecipeName(t *testing.T) {
	userID := uuid.New()
	ingredientUsage := &IngredientUsage{
		IngredientID: uuid.New(),
		Quantity:     200,
		Unit:         1,
	}
	_, err := NewRecipe("", userID, []*IngredientUsage{ingredientUsage}, 10, 20, 4)
	require.Error(t, err)
	require.Equal(t, ErrEmptyName, err)
}

func TestNoIngredients(t *testing.T) {
	name := "Pancakes"
	userID := uuid.New()
	_, err := NewRecipe(name, userID, []*IngredientUsage{}, 10, 20, 4)
	require.Error(t, err)
	require.Equal(t, ErrNoIngredients, err)
}

func TestNegativePrepMins(t *testing.T) {
	name := "Pancakes"
	userID := uuid.New()
	ingredientUsage := &IngredientUsage{
		IngredientID: uuid.New(),
		Quantity:     200,
		Unit:         1,
	}
	_, err := NewRecipe(name, userID, []*IngredientUsage{ingredientUsage}, -5, 20, 4)
	require.Error(t, err)
	require.Equal(t, ErrInvalidPrepMins, err)
}

func TestNegativeCookMins(t *testing.T) {
	name := "Pancakes"
	userID := uuid.New()
	ingredientUsage := &IngredientUsage{
		IngredientID: uuid.New(),
		Quantity:     200,
		Unit:         1,
	}
	_, err := NewRecipe(name, userID, []*IngredientUsage{ingredientUsage}, 10, -5, 4)
	require.Error(t, err)
	require.Equal(t, ErrInvalidCookMins, err)
}

func TestInvalidPortions(t *testing.T) {
	name := "Pancakes"
	userID := uuid.New()
	ingredientUsage := &IngredientUsage{
		IngredientID: uuid.New(),
		Quantity:     200,
		Unit:         1,
	}
	_, err := NewRecipe(name, userID, []*IngredientUsage{ingredientUsage}, 10, 20, 0)
	require.Error(t, err)
	require.Equal(t, ErrInvalidPortions, err)
}
