package recipe

import (
	"context"
	"database/sql"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/testutil"
	"foodplanner/internal/testutil/seeds"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetRecipeSourceByRecipeVersionID(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		r := NewRecipeVersionRepo()
		recipeRepo := NewRecipeRepo()
		ingredientID := uuid.New()
		testIngredient := ingredient.Ingredient{
			ID:            ingredientID,
			FileKey:       "test_ingredient",
			Name:          "Test Ingredient",
			PreferredUnit: 1,
		}
		err := seeds.InsertIngredient(context.Background(), tx, &testIngredient)
		require.NoError(t, err)

		testUser, err := seeds.SeedTestUser(context.Background(), tx)
		require.NoError(t, err)

		ingredientUsage, err := NewIngredientUsage(CreateIngredientUsageRequest{
			IngredientID: ingredientID.String(),
			Quantity:     200,
			Unit:         1,
		})
		require.NoError(t, err)

		source, err := NewURLSource("https://example.com/recipe")
		require.NoError(t, err)

		recipeContainer, err := NewRecipe(
			"Recipe With Source",
			testUser.ID,
			[]*IngredientUsage{ingredientUsage},
			30,
			60,
			8,
			source,
		)
		require.NoError(t, err)
		created, err := recipeRepo.createRecipe(context.Background(), tx, recipeContainer)
		require.NoError(t, err)

		// Act
		retrievedSource, err := r.getRecipeSourceByRecipeVersionID(context.Background(), tx, created.CurrentVersion.ID.String())

		// Assert
		require.NoError(t, err)
		require.NotNil(t, retrievedSource, "Expected to find recipe source")
		require.Equal(t, URL, retrievedSource.Type)
		require.Equal(t, "https://example.com/recipe", *retrievedSource.URL)
	})
}
