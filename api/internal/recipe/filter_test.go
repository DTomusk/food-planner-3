package recipe

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterHashForParams(t *testing.T) {
	query := "pasta"
	nilQuery := (*string)(nil)

	t.Run("same inputs produce same hash", func(t *testing.T) {
		h1 := filterHashForParams(RecipeCursorModeNewest, nilQuery, nil)
		h2 := filterHashForParams(RecipeCursorModeNewest, nilQuery, nil)
		require.Equal(t, h1, h2)
	})

	t.Run("different modes produce different hashes", func(t *testing.T) {
		h1 := filterHashForParams(RecipeCursorModeNewest, nilQuery, nil)
		h2 := filterHashForParams(RecipeCursorModeRelevance, &query, nil)
		require.NotEqual(t, h1, h2)
	})

	t.Run("different queries produce different hashes", func(t *testing.T) {
		other := "chicken"
		h1 := filterHashForParams(RecipeCursorModeRelevance, &query, nil)
		h2 := filterHashForParams(RecipeCursorModeRelevance, &other, nil)
		require.NotEqual(t, h1, h2)
	})

	t.Run("nil query does not panic", func(t *testing.T) {
		require.NotPanics(t, func() {
			filterHashForParams(RecipeCursorModeNewest, nil, nil)
		})
	})

	t.Run("hash is 16 hex characters", func(t *testing.T) {
		h := filterHashForParams(RecipeCursorModeNewest, nilQuery, nil)
		require.Len(t, h, 16)
	})
}
