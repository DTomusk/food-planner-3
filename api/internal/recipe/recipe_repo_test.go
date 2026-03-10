package recipe

import (
	"context"
	"database/sql"
	"strconv"
	"testing"

	"foodplanner/internal/ingredient"
	"foodplanner/internal/testutil"
	"foodplanner/internal/testutil/seeds"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCreateAndGetRecipe(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		r := NewRecipeRepo()
		ingredientID := uuid.New()
		testIngredient := ingredient.Ingredient{
			ID:            ingredientID,
			FileKey:       "test_ingredient",
			Name:          "Test Ingredient",
			PreferredUnit: 1,
		}
		err := seeds.InsertIngredient(context.Background(), tx, &testIngredient)
		require.NoError(t, err, "Failed to seed test ingredient")
		ingredientUsage, err := NewIngredientUsage(CreateIngredientUsageRequest{
			IngredientID: ingredientID.String(),
			Quantity:     200,
			Unit:         1,
		})
		require.NoError(t, err, "Failed to create ingredient usage")

		testUser, err := seeds.SeedTestUser(context.Background(), tx)
		require.NoError(t, err, "Failed to seed test user")

		source := &RecipeSource{
			Type: URL,
			URL:  testutil.PtrString("https://example.com/pancakes"),
		}

		recipeContainer, err := NewRecipe(
			"Chocolate Cake",
			testUser.ID,
			[]*IngredientUsage{ingredientUsage},
			30,
			60,
			8,
			source,
		)
		require.NoError(t, err, "Failed to create recipe")

		// Act
		_, err = r.createRecipe(context.Background(), tx, recipeContainer)
		require.NoError(t, err, "Failed to create recipe in database")

		gotContainer, err := r.getRecipeByID(context.Background(), tx, recipeContainer.ID.String())
		require.NoError(t, err, "Failed to get recipe by ID")

		// Assert
		require.NotNil(t, gotContainer, "Expected to find recipe by ID")
		require.NotNil(t, gotContainer.CurrentVersion, "Expected to find recipe version by ID")
		require.Equal(t, recipeContainer.CurrentVersion.ID, gotContainer.CurrentVersion.ID, "Expected recipe ID to match")
		require.Equal(t, gotContainer.CurrentVersion.Name, "Chocolate Cake", "Expected recipe name to match")
	})
}

func TestGetRecipe_DoesNotErrorWhenNotFound(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		r := NewRecipeRepo()

		// Act
		recipeContainer, err := r.getRecipeByID(context.Background(), tx, "04061e4e-6d4c-41d1-abcf-8b214927e1ed")

		// Assert
		require.NoError(t, err)
		require.Nil(t, recipeContainer)
	})
}

func TestGetRecipesByUserID(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		r := NewRecipeRepo()
		ingredientID := uuid.New()
		testIngredient := ingredient.Ingredient{
			ID:            ingredientID,
			FileKey:       "test_ingredient",
			Name:          "Test Ingredient",
			PreferredUnit: 1,
		}
		err := seeds.InsertIngredient(context.Background(), tx, &testIngredient)
		require.NoError(t, err, "Failed to seed test ingredient")

		testUser, err := seeds.SeedTestUser(context.Background(), tx)
		require.NoError(t, err, "Failed to seed test user")

		otherUser, err := seeds.SeedTestUser(context.Background(), tx)
		require.NoError(t, err, "Failed to seed other test user")

		// Create two recipes for test user
		for i := 0; i < 2; i++ {
			ingredientUsage, err := NewIngredientUsage(CreateIngredientUsageRequest{
				IngredientID: ingredientID.String(),
				Quantity:     200,
				Unit:         1,
			})
			require.NoError(t, err)

			recipeContainer, err := NewRecipe(
				"Test Recipe "+strconv.Itoa(i),
				testUser.ID,
				[]*IngredientUsage{ingredientUsage},
				30,
				60,
				8,
				nil,
			)
			require.NoError(t, err)
			_, err = r.createRecipe(context.Background(), tx, recipeContainer)
			require.NoError(t, err)
		}

		// Create one recipe for other user
		ingredientUsage, err := NewIngredientUsage(CreateIngredientUsageRequest{
			IngredientID: ingredientID.String(),
			Quantity:     200,
			Unit:         1,
		})
		require.NoError(t, err)
		otherRecipe, err := NewRecipe(
			"Other User Recipe",
			otherUser.ID,
			[]*IngredientUsage{ingredientUsage},
			30,
			60,
			8,
			nil,
		)
		require.NoError(t, err)
		_, err = r.createRecipe(context.Background(), tx, otherRecipe)
		require.NoError(t, err)

		// Act
		recipes, err := r.getRecipesByUserID(context.Background(), tx, testUser.ID.String())

		// Assert
		require.NoError(t, err)
		require.Len(t, recipes, 2, "Expected to find 2 recipes for test user")
		for _, recipe := range recipes {
			require.Equal(t, testUser.ID, recipe.UserID, "Expected all recipes to belong to test user")
		}
	})
}
