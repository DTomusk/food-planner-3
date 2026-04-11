package recipe

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterHash(t *testing.T) {
	query := "pasta"
	nilQuery := (*string)(nil)

	t.Run("same inputs produce same hash", func(t *testing.T) {
		h1 := filterHash(RecipeCursorModeNewest, nilQuery, normalizedRecipeFilter{})
		h2 := filterHash(RecipeCursorModeNewest, nilQuery, normalizedRecipeFilter{})
		require.Equal(t, h1, h2)
	})

	t.Run("different modes produce different hashes", func(t *testing.T) {
		h1 := filterHash(RecipeCursorModeNewest, nilQuery, normalizedRecipeFilter{})
		h2 := filterHash(RecipeCursorModeRelevance, &query, normalizedRecipeFilter{})
		require.NotEqual(t, h1, h2)
	})

	t.Run("different queries produce different hashes", func(t *testing.T) {
		other := "chicken"
		h1 := filterHash(RecipeCursorModeRelevance, &query, normalizedRecipeFilter{})
		h2 := filterHash(RecipeCursorModeRelevance, &other, normalizedRecipeFilter{})
		require.NotEqual(t, h1, h2)
	})

	t.Run("nil query does not panic", func(t *testing.T) {
		require.NotPanics(t, func() {
			filterHash(RecipeCursorModeNewest, nil, normalizedRecipeFilter{})
		})
	})

	t.Run("hash is 16 hex characters", func(t *testing.T) {
		h := filterHash(RecipeCursorModeNewest, nilQuery, normalizedRecipeFilter{})
		require.Len(t, h, 16)
	})

	t.Run("different contains gluten filters produce different hashes", func(t *testing.T) {
		containsGluten := true
		glutenFree := false
		h1 := filterHash(RecipeCursorModeNewest, nilQuery, normalizedRecipeFilter{ContainsGluten: &containsGluten})
		h2 := filterHash(RecipeCursorModeNewest, nilQuery, normalizedRecipeFilter{ContainsGluten: &glutenFree})
		require.NotEqual(t, h1, h2)
	})
}

func TestNormalizedAnimalProductLevelFilter(t *testing.T) {
	t.Run("nil remains nil", func(t *testing.T) {
		require.Nil(t, normalizedAnimalProductLevelFilter(nil))
	})

	t.Run("zero is kept", func(t *testing.T) {
		v := 0
		result := normalizedAnimalProductLevelFilter(&v)
		require.NotNil(t, result)
		require.Equal(t, 0, *result)
	})

	t.Run("one is kept", func(t *testing.T) {
		v := 1
		result := normalizedAnimalProductLevelFilter(&v)
		require.NotNil(t, result)
		require.Equal(t, 1, *result)
	})

	t.Run("unsupported values are treated as nil", func(t *testing.T) {
		for _, v := range []int{-1, 2, 3, 99} {
			value := v
			require.Nil(t, normalizedAnimalProductLevelFilter(&value))
		}
	})
}
