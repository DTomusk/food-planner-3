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

func TestGetRecipe_DoesNotErrorWhenNotFound(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		r := NewRecipeRepo()

		// Act
		recipeContainer, err := r.getRecipeByID(context.Background(), tx, uuid.MustParse("04061e4e-6d4c-41d1-abcf-8b214927e1ed"))

		// Assert
		require.NoError(t, err)
		require.Nil(t, recipeContainer)
	})
}

func TestGetRecipesByUserID(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		ctx := context.Background()
		r := NewRecipeRepo()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		otherUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed other test user")

		expectedRecipeNames := []string{"Test Recipe 0", "Test Recipe 1"}
		for _, recipeName := range expectedRecipeNames {
			recipeID := uuid.New()
			versionID := uuid.New()

			err = seeds.InsertRecipeContainer(ctx, tx, recipeID, testUser.ID)
			require.NoError(t, err)

			err = seeds.InsertRecipeVersion(ctx, tx, versionID, recipeID, recipeName, 30, 60, 8, 1)
			require.NoError(t, err)

			err = seeds.SetRecipeContainerCurrentVersion(ctx, tx, recipeID, versionID)
			require.NoError(t, err)
		}

		otherRecipeID := uuid.New()
		otherVersionID := uuid.New()

		err = seeds.InsertRecipeContainer(ctx, tx, otherRecipeID, otherUser.ID)
		require.NoError(t, err)

		err = seeds.InsertRecipeVersion(ctx, tx, otherVersionID, otherRecipeID, "Other User Recipe", 30, 60, 8, 1)
		require.NoError(t, err)

		err = seeds.SetRecipeContainerCurrentVersion(ctx, tx, otherRecipeID, otherVersionID)
		require.NoError(t, err)

		// Act
		recipes, err := r.getRecipesByUserID(ctx, tx, testUser.ID)

		// Assert
		require.NoError(t, err)
		require.Len(t, recipes, 2, "Expected to find 2 recipes for test user")

		actualNames := make(map[string]struct{}, len(recipes))
		for _, recipe := range recipes {
			require.Equal(t, testUser.ID, recipe.UserID, "Expected all recipes to belong to test user")
			require.NotNil(t, recipe.CurrentVersion, "Expected recipe to include current version")
			require.Equal(t, recipe.ID, recipe.CurrentVersion.RecipeID, "Expected current version to point to recipe container")
			actualNames[recipe.CurrentVersion.Name] = struct{}{}
		}

		for _, expectedName := range expectedRecipeNames {
			_, found := actualNames[expectedName]
			require.True(t, found, "Expected recipe name to be present")
		}
	})
}
