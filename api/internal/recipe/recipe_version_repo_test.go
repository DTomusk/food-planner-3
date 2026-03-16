package recipe

import (
	"context"
	"database/sql"
	"testing"

	"foodplanner/internal/testutil"
	"foodplanner/internal/testutil/seeds"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestGetRecipeSourceByRecipeVersionID(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		ctx := context.Background()
		r := NewRecipeVersionRepo()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		recipeID := uuid.New()
		versionID := uuid.New()

		err = seeds.InsertRecipeContainer(ctx, tx, recipeID, testUser.ID)
		require.NoError(t, err)
		err = seeds.InsertRecipeVersion(ctx, tx, versionID, recipeID, "Recipe With Source", 30, 60, 8, 1)
		require.NoError(t, err)
		err = seeds.SetRecipeContainerCurrentVersion(ctx, tx, recipeID, versionID)
		require.NoError(t, err)

		url := "https://example.com/recipe"
		err = seeds.InsertRecipeSource(ctx, tx, versionID, int(URL), &url, nil, nil, nil)
		require.NoError(t, err)

		// Act
		retrievedSource, err := r.getRecipeSourceByRecipeVersionID(ctx, tx, versionID)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, retrievedSource, "Expected to find recipe source")
		require.Equal(t, URL, retrievedSource.Type)
		require.NotNil(t, retrievedSource.URL)
		require.Equal(t, url, *retrievedSource.URL)
	})
}
