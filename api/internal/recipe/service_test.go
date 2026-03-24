package recipe

import (
	"context"
	"database/sql"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/testutil"
	"foodplanner/internal/testutil/seeds"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateRecipe(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()

		ingredientID := uuid.New()
		testIngredient := ingredient.Ingredient{
			ID:            ingredientID,
			FileKey:       "test_ingredient",
			Name:          "Test Ingredient",
			PreferredUnit: 1,
		}
		err := seeds.InsertIngredient(ctx, tx, &testIngredient)
		require.NoError(t, err, "Failed to seed test ingredient")

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		s := NewService(
			txRunner,
			NewRecipeRepo(),
			NewRecipeVersionRepo(),
			ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
			NewIngredientUsageRepo(),
			nil,
		)
		ingredientRequest := CreateIngredientUsageRequest{
			IngredientID: ingredientID.String(),
			Quantity:     200,
			Unit:         1,
		}

		source := CreateRecipeSourceRequest{
			Type: 1,
			URL:  testutil.PtrString("https://example.com/pancakes"),
		}

		request := CreateRecipeRequest{
			Name:        "Vanilla Ice Cream",
			Ingredients: []CreateIngredientUsageRequest{ingredientRequest},
			UserID:      testUser.ID.String(),
			PrepMins:    15,
			CookMins:    0,
			Portions:    6,
			Source:      source,
		}

		// Act
		recipeContainer, err := s.CreateRecipe(ctx, request)
		recipeVersion := recipeContainer.CurrentVersion

		// Assert
		require.NoError(t, err, "Expected no error when creating recipe")

		require.Equal(t, recipeContainer.UserID, uuid.MustParse(request.UserID), "Expected user ID to match the request")
		require.Equal(t, recipeContainer.ID, recipeVersion.RecipeID, "Expected recipe ID to match the version's recipe ID")
		require.NotNil(t, recipeContainer.CreatedAt, "Expected CreatedAt to be set")

		require.Equal(t, "Vanilla Ice Cream", recipeVersion.Name, "Expected recipe name to match the request")
		require.Equal(t, 15, recipeVersion.PrepMins, "Expected prep minutes to match the request")
		require.Equal(t, 0, recipeVersion.CookMins, "Expected cook minutes to match the request")
		require.Equal(t, 6, recipeVersion.Portions, "Expected portions to match the request")
	})
}

func TestCreateRecipeWithDuplicateIngredients(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()
		ingredientID := uuid.New()
		testIngredient := ingredient.Ingredient{
			ID:            ingredientID,
			FileKey:       "test_ingredient",
			Name:          "Test Ingredient",
			PreferredUnit: 1,
		}
		err := seeds.InsertIngredient(ctx, tx, &testIngredient)
		require.NoError(t, err, "Failed to seed test ingredient")

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		s := NewService(
			txRunner,
			NewRecipeRepo(),
			NewRecipeVersionRepo(),
			ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
			NewIngredientUsageRepo(),
			nil,
		)
		ingredientRequest := CreateIngredientUsageRequest{
			IngredientID: ingredientID.String(),
			Quantity:     200,
			Unit:         1,
		}
		request := CreateRecipeRequest{
			Name:        "Vanilla Ice Cream",
			Ingredients: []CreateIngredientUsageRequest{ingredientRequest, ingredientRequest},
			UserID:      testUser.ID.String(),
			PrepMins:    15,
			CookMins:    0,
			Portions:    6,
		}
		_, err = s.CreateRecipe(ctx, request)
		require.ErrorIs(t, err, ErrDuplicateIngredient)
	})
}

func TestCreateRecipeWithNonexistentIngredient(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()
		ingredientID := uuid.New()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")
		s := NewService(
			txRunner,
			NewRecipeRepo(),
			NewRecipeVersionRepo(),
			ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
			NewIngredientUsageRepo(),
			nil,
		)

		ingredientRequest := CreateIngredientUsageRequest{
			IngredientID: ingredientID.String(),
			Quantity:     200,
			Unit:         1,
		}

		request := CreateRecipeRequest{
			Name:        "Vanilla Ice Cream",
			Ingredients: []CreateIngredientUsageRequest{ingredientRequest},
			UserID:      testUser.ID.String(),
			PrepMins:    15,
			CookMins:    0,
			Portions:    6,
		}
		_, err = s.CreateRecipe(ctx, request)
		require.ErrorIs(t, err, ErrIngredientNotFound)
	})
}

func TestCreateRecipeWithInvalidUnit(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()
		ingredientID := uuid.New()
		testIngredient := ingredient.Ingredient{
			ID:            ingredientID,
			FileKey:       "test_ingredient",
			Name:          "Test Ingredient",
			PreferredUnit: 1,
		}
		err := seeds.InsertIngredient(ctx, tx, &testIngredient)
		require.NoError(t, err, "Failed to seed test ingredient")
		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")
		s := NewService(
			txRunner,
			NewRecipeRepo(),
			NewRecipeVersionRepo(),
			ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
			NewIngredientUsageRepo(),
			nil,
		)
		ingredientRequest := CreateIngredientUsageRequest{
			IngredientID: ingredientID.String(),
			Quantity:     200,
			Unit:         999, // this unit doesn't exist (for now!)
		}
		request := CreateRecipeRequest{
			Name:        "Vanilla Ice Cream",
			Ingredients: []CreateIngredientUsageRequest{ingredientRequest},
			UserID:      testUser.ID.String(),
			PrepMins:    15,
			CookMins:    0,
			Portions:    6,
		}
		_, err = s.CreateRecipe(ctx, request)
		require.ErrorIs(t, err, ErrInvalidUnit)
	})
}

func TestCreateRecipeNotPreferredUnit(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()
		ingredientID := uuid.New()
		testIngredient := ingredient.Ingredient{
			ID:            ingredientID,
			FileKey:       "test_ingredient",
			Name:          "Test Ingredient",
			PreferredUnit: 1,
		}
		err := seeds.InsertIngredient(ctx, tx, &testIngredient)
		require.NoError(t, err, "Failed to seed test ingredient")
		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")
		s := NewService(
			txRunner,
			NewRecipeRepo(),
			NewRecipeVersionRepo(),
			ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
			NewIngredientUsageRepo(),
			nil,
		)
		ingredientRequest := CreateIngredientUsageRequest{
			IngredientID: ingredientID.String(),
			Quantity:     200,
			Unit:         2, // valid unit, but it's not the preferred one
		}
		request := CreateRecipeRequest{
			Name:        "Vanilla Ice Cream",
			Ingredients: []CreateIngredientUsageRequest{ingredientRequest},
			UserID:      testUser.ID.String(),
			PrepMins:    15,
			CookMins:    0,
			Portions:    6,
		}
		_, err = s.CreateRecipe(ctx, request)
		require.ErrorIs(t, err, ErrInvalidUnit)
	})
}

func TestCreateRecipe_NoSource(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()

		ingredientID := uuid.New()
		testIngredient := ingredient.Ingredient{
			ID:            ingredientID,
			FileKey:       "test_ingredient",
			Name:          "Test Ingredient",
			PreferredUnit: 1,
		}
		err := seeds.InsertIngredient(ctx, tx, &testIngredient)
		require.NoError(t, err, "Failed to seed test ingredient")

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		s := NewService(
			txRunner,
			NewRecipeRepo(),
			NewRecipeVersionRepo(),
			ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
			NewIngredientUsageRepo(),
			nil,
		)
		ingredientRequest := CreateIngredientUsageRequest{
			IngredientID: ingredientID.String(),
			Quantity:     200,
			Unit:         1,
		}

		request := CreateRecipeRequest{
			Name:        "Vanilla Ice Cream",
			Ingredients: []CreateIngredientUsageRequest{ingredientRequest},
			UserID:      testUser.ID.String(),
			PrepMins:    15,
			CookMins:    0,
			Portions:    6,
		}

		// Act
		_, err = s.CreateRecipe(ctx, request)

		// Assert
		require.NoError(t, err)
	})
}

func TestGetRecipes_PaginatesAcrossPages(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		txRunner := testutil.NewTestTxRunner(tx)
		s := NewService(
			txRunner,
			NewRecipeRepo(),
			NewRecipeVersionRepo(),
			ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
			NewIngredientUsageRepo(),
			nil,
		)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		newestCreatedAt := time.Date(2026, time.March, 17, 11, 42, 48, 147630000, time.UTC)
		middleCreatedAt := newestCreatedAt.Add(-1 * time.Minute)
		oldestCreatedAt := newestCreatedAt.Add(-2 * time.Minute)

		newest := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1"), uuid.New(), "Newest", newestCreatedAt, nil)
		middle := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2"), uuid.New(), "Middle", middleCreatedAt, nil)
		oldest := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("cccccccc-cccc-cccc-cccc-ccccccccccc3"), uuid.New(), "Oldest", oldestCreatedAt, nil)

		firstPage, nextCursor, err := s.GetRecipes(ctx, 2, nil)

		require.NoError(t, err)
		require.Len(t, firstPage, 2)
		require.Equal(t, newest.RecipeID, firstPage[0].ID)
		require.Equal(t, middle.RecipeID, firstPage[1].ID)
		require.NotNil(t, nextCursor)

		parsedCursor, err := ParseRecipeCursor(nextCursor)
		require.NoError(t, err)
		require.NotNil(t, parsedCursor)
		require.True(t, middle.CreatedAt.Equal(parsedCursor.CreatedAt))
		require.Equal(t, middle.RecipeID, parsedCursor.ID)

		secondPage, finalCursor, err := s.GetRecipes(ctx, 2, nextCursor)

		require.NoError(t, err)
		require.Len(t, secondPage, 1)
		require.Equal(t, oldest.RecipeID, secondPage[0].ID)
		require.Nil(t, finalCursor)
	})
}

func TestGetRecipes_InvalidCursor(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		txRunner := testutil.NewTestTxRunner(tx)
		s := NewService(
			txRunner,
			NewRecipeRepo(),
			NewRecipeVersionRepo(),
			ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
			NewIngredientUsageRepo(),
			nil,
		)

		invalidCursor := "not-a-valid-cursor"

		recipes, nextCursor, err := s.GetRecipes(ctx, 2, &invalidCursor)

		require.ErrorIs(t, err, ErrInvalidCursor)
		require.Nil(t, recipes)
		require.Nil(t, nextCursor)
	})
}
