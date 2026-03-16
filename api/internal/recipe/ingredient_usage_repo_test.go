package recipe

import (
	"context"
	"database/sql"
	"testing"

	"foodplanner/internal/ingredient"
	"foodplanner/internal/testutil"
	"foodplanner/internal/testutil/seeds"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestGetIngredientUsagesForRecipeVersion(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		r := NewIngredientUsageRepo()

		ingredientID := uuid.New()
		err := seeds.InsertIngredient(ctx, tx, &ingredient.Ingredient{
			ID:            ingredientID,
			FileKey:       "test_ingredient",
			Name:          "Test Ingredient",
			PreferredUnit: 1,
		})
		require.NoError(t, err)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		recipeID := uuid.New()
		versionID := uuid.New()

		err = seeds.InsertRecipeContainer(ctx, tx, recipeID, testUser.ID)
		require.NoError(t, err)
		err = seeds.InsertRecipeVersion(ctx, tx, versionID, recipeID, "Recipe With Ingredients", 30, 60, 8, 1)
		require.NoError(t, err)
		err = seeds.SetRecipeContainerCurrentVersion(ctx, tx, recipeID, versionID)
		require.NoError(t, err)

		ingredientUsage, err := NewIngredientUsage(CreateIngredientUsageRequest{
			IngredientID: ingredientID.String(),
			Quantity:     200,
			Unit:         1,
		})
		require.NoError(t, err)

		err = r.insertIngredientUsages(ctx, tx, []*IngredientUsage{ingredientUsage}, versionID)
		require.NoError(t, err, "Failed to insert ingredient usages")

		usages, err := r.getIngredientUsagesForRecipeVersion(ctx, tx, versionID)

		require.NoError(t, err)
		require.Len(t, usages, 1, "Expected to find 1 ingredient usage")
		require.Equal(t, ingredientID, usages[0].IngredientID)
		require.Equal(t, 200.0, usages[0].Quantity)
		require.Equal(t, ingredientUsage.Unit, usages[0].Unit)
	})
}

func TestInsertIngredientUsages(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		ctx := context.Background()
		r := NewIngredientUsageRepo()

		ingredientID1 := uuid.New()
		ingredientID2 := uuid.New()

		err := seeds.InsertIngredient(ctx, tx, &ingredient.Ingredient{
			ID:            ingredientID1,
			FileKey:       "ingredient_1",
			Name:          "Ingredient 1",
			PreferredUnit: 1,
		})
		require.NoError(t, err)

		err = seeds.InsertIngredient(ctx, tx, &ingredient.Ingredient{
			ID:            ingredientID2,
			FileKey:       "ingredient_2",
			Name:          "Ingredient 2",
			PreferredUnit: 2,
		})
		require.NoError(t, err)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		recipeID := uuid.New()
		versionID := uuid.New()

		err = seeds.InsertRecipeContainer(ctx, tx, recipeID, testUser.ID)
		require.NoError(t, err)
		err = seeds.InsertRecipeVersion(ctx, tx, versionID, recipeID, "Recipe With Multiple Ingredients", 30, 60, 8, 1)
		require.NoError(t, err)
		err = seeds.SetRecipeContainerCurrentVersion(ctx, tx, recipeID, versionID)
		require.NoError(t, err)

		usage1, err := NewIngredientUsage(CreateIngredientUsageRequest{
			IngredientID: ingredientID1.String(),
			Quantity:     200,
			Unit:         1,
		})
		require.NoError(t, err)

		usage2, err := NewIngredientUsage(CreateIngredientUsageRequest{
			IngredientID: ingredientID2.String(),
			Quantity:     300,
			Unit:         2,
		})
		require.NoError(t, err)

		// Act
		err = r.insertIngredientUsages(ctx, tx, []*IngredientUsage{usage1, usage2}, versionID)
		require.NoError(t, err, "Failed to insert ingredient usages")

		retrievedUsages, err := r.getIngredientUsagesForRecipeVersion(ctx, tx, versionID)

		// Assert
		require.NoError(t, err)
		require.Len(t, retrievedUsages, 2, "Expected to find 2 ingredient usages")

		ingredientIDs := map[uuid.UUID]bool{}
		for _, usage := range retrievedUsages {
			ingredientIDs[usage.IngredientID] = true
		}
		require.True(t, ingredientIDs[ingredientID1])
		require.True(t, ingredientIDs[ingredientID2])
	})
}
