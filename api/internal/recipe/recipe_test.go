package recipe

import (
	"foodplanner/internal/testutil"
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
	source := &RecipeSource{
		Type: URL,
		URL:  testutil.PtrString("https://example.com/pancakes"),
	}
	recipeContainer, err := NewRecipe(name, userID, []*IngredientUsage{ingredientUsage}, 10, 20, 4, source, nil)
	require.NoError(t, err)
	recipeVersion := recipeContainer.CurrentVersion

	require.Equal(t, recipeContainer.ID, recipeVersion.RecipeID)
	require.Equal(t, name, recipeVersion.Name)
	require.Equal(t, 10, recipeVersion.PrepMins)
	require.Equal(t, 20, recipeVersion.CookMins)
	require.Equal(t, 4, recipeVersion.Portions)
	require.NotNil(t, recipeVersion.Source)
	require.Equal(t, URL, recipeVersion.Source.Type)
	require.Equal(t, testutil.PtrString("https://example.com/pancakes"), recipeVersion.Source.URL)

	require.Equal(t, userID, recipeContainer.UserID)
	require.Equal(t, recipeVersion.ID, recipeContainer.CurrentVersionID)
}

func TestEmptyRecipeName(t *testing.T) {
	userID := uuid.New()
	ingredientUsage := &IngredientUsage{
		IngredientID: uuid.New(),
		Quantity:     200,
		Unit:         1,
	}
	source := &RecipeSource{
		Type: URL,
		URL:  testutil.PtrString("https://example.com/pancakes"),
	}
	_, err := NewRecipe("", userID, []*IngredientUsage{ingredientUsage}, 10, 20, 4, source, nil)
	require.Error(t, err)
	require.Equal(t, ErrEmptyName, err)
}

func TestNoIngredients(t *testing.T) {
	name := "Pancakes"
	userID := uuid.New()
	source := &RecipeSource{
		Type: URL,
		URL:  testutil.PtrString("https://example.com/pancakes"),
	}
	_, err := NewRecipe(name, userID, []*IngredientUsage{}, 10, 20, 4, source, nil)
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
	source := &RecipeSource{
		Type: URL,
		URL:  testutil.PtrString("https://example.com/pancakes"),
	}
	_, err := NewRecipe(name, userID, []*IngredientUsage{ingredientUsage}, -5, 20, 4, source, nil)
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
	source := &RecipeSource{
		Type: URL,
		URL:  testutil.PtrString("https://example.com/pancakes"),
	}
	_, err := NewRecipe(name, userID, []*IngredientUsage{ingredientUsage}, 10, -5, 4, source, nil)
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
	source := &RecipeSource{
		Type: URL,
		URL:  testutil.PtrString("https://example.com/pancakes"),
	}
	_, err := NewRecipe(name, userID, []*IngredientUsage{ingredientUsage}, 10, 20, 0, source, nil)
	require.Error(t, err)
	require.Equal(t, ErrInvalidPortions, err)
}
