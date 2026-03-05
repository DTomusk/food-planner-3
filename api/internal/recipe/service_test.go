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

		s := NewService(txRunner, NewRepo(), ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100), nil)
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
		recipe, err := s.CreateRecipe(ctx, request)

		// Assert
		require.NoError(t, err, "Expected no error when creating recipe")
		require.Equal(t, "Vanilla Ice Cream", recipe.Name, "Expected recipe name to match the request")
		require.Equal(t, 15, recipe.PrepMins, "Expected prep minutes to match the request")
		require.Equal(t, 0, recipe.CookMins, "Expected cook minutes to match the request")
		require.Equal(t, 6, recipe.Portions, "Expected portions to match the request")
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

		s := NewService(txRunner, NewRepo(), ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100), nil)
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
		s := NewService(txRunner, NewRepo(), ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100), nil)
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
		s := NewService(txRunner, NewRepo(), ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100), nil)
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
		s := NewService(txRunner, NewRepo(), ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100), nil)
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

		s := NewService(txRunner, NewRepo(), ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100), nil)
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

func TestCreateAndDeleteRecipe(t *testing.T) {
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

		testUserIdPtr := testutil.PtrString(testUser.ID.String())

		s := NewService(txRunner, NewRepo(), ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100), nil)
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

		// Act 1 - Create recipe
		recipe, err := s.CreateRecipe(ctx, request)

		// Assert
		require.NoError(t, err, "Expected no error when creating recipe")
		require.Nil(t, recipe.DeletedOn, "Expected DeletedOn to be nil for newly created recipe")

		// Act 2 - Delete recipe
		recipe, err = s.DeleteRecipe(ctx, recipe.ID.String(), testUser.ID.String())

		// Assert
		require.NoError(t, err, "Expected no error when deleting recipe")
		require.NotNil(t, recipe.DeletedOn, "Expected DeletedOn to be set after deleting recipe")

		// Act 3 - Try to get deleted recipe
		got, err := s.GetRecipeByID(ctx, recipe.ID.String())

		// Assert
		require.NoError(t, err, "Expected no error when getting recipe by ID")
		require.Nil(t, got, "Expected to not find deleted recipe by ID")

		// Act 4 - Get deleted recipes for user
		deletedRecipes, err := s.GetRecipesByUserID(ctx, testUser.ID.String(), StatusDeleted, testUserIdPtr)

		// Assert
		require.NoError(t, err, "Expected no error when getting deleted recipes by user ID")
		require.Len(t, deletedRecipes, 1, "Expected to find one deleted recipe for user")
		require.Equal(t, recipe.ID, deletedRecipes[0].ID, "Expected deleted recipe ID to match created recipe ID")

		// Act 5 - Undelete recipe
		undeletedRecipe, err := s.UndeleteRecipe(ctx, recipe.ID.String(), testUser.ID.String())

		// Assert
		require.NoError(t, err, "Expected no error when undeleting recipe")
		require.Nil(t, undeletedRecipe.DeletedOn, "Expected DeletedOn to be nil after undeleting recipe")

		// Act 6 - Get undeleted recipe by ID
		got, err = s.GetRecipeByID(ctx, recipe.ID.String())

		// Assert
		require.NoError(t, err, "Expected no error when getting recipe by ID")
		require.NotNil(t, got, "Expected to find undeleted recipe by ID")
		require.Equal(t, recipe.ID, got.ID, "Expected undeleted recipe ID to match created recipe ID")

		// Act 7 - Get deleted recipes for user again
		deletedRecipes, err = s.GetRecipesByUserID(ctx, testUser.ID.String(), StatusDeleted, testUserIdPtr)

		// Assert
		require.NoError(t, err, "Expected no error when getting deleted recipes by user ID")
		require.Len(t, deletedRecipes, 0, "Expected to find no deleted recipes for user after undeleting")
	})
}

func TestDeleteOldRecipes(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()
		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")
		retentionPeriod := 30
		s := NewService(txRunner, NewRepo(), ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100), &retentionPeriod)

		// Act 1 - delete old recipes with empty repo
		deletedCount, err := s.DeleteOldRecipes(ctx)

		// Assert
		require.NoError(t, err, "Expected no error when deleting old recipes with empty repo")
		require.Equal(t, int64(0), deletedCount, "Expected to delete 0 old recipes with empty repo")

		// Arrange - create not deleted recipe
		recipeID := uuid.New()
		err = seeds.InsertRecipe(ctx, tx, recipeID.String(), testUser.ID.String(), "Test Recipe", 10, 20, 2, nil)
		require.NoError(t, err, "Failed to insert test recipe")

		// Act 2 - delete old recipes with no deleted recipes
		deletedCount, err = s.DeleteOldRecipes(ctx)

		// Assert
		require.NoError(t, err, "Expected no error when deleting old recipes with no deleted recipes")
		require.Equal(t, int64(0), deletedCount, "Expected to delete 0 old recipes when there are no deleted recipes")

		// Arrange - delete recipe with deleted_on > retention period
		deletedRecipeID := uuid.New()
		deletedOn := time.Now().AddDate(0, 0, -1) // 1 day ago
		err = seeds.InsertRecipe(ctx, tx, deletedRecipeID.String(), testUser.ID.String(), "Deleted Recipe", 10, 20, 2, &deletedOn)
		require.NoError(t, err, "Failed to insert deleted test recipe")

		// Act 3 - delete old recipes with a deleted recipe that is not old enough
		deletedCount, err = s.DeleteOldRecipes(ctx)

		// Assert
		require.NoError(t, err, "Expected no error when deleting old recipes with a deleted recipe that is not old enough")
		require.Equal(t, int64(0), deletedCount, "Expected to delete 0 old recipes when the deleted recipe is not old enough")

		// Arrange - delete recipe with deleted_on < retention period
		oldDeletedRecipeID := uuid.New()
		oldDeletedOn := time.Now().AddDate(0, 0, -31)
		err = seeds.InsertRecipe(ctx, tx, oldDeletedRecipeID.String(), testUser.ID.String(), "Old Deleted Recipe", 10, 20, 2, &oldDeletedOn)
		require.NoError(t, err, "Failed to insert old deleted test recipe")

		// Act 4 - delete old recipes with a deleted recipe that is old enough
		deletedCount, err = s.DeleteOldRecipes(ctx)

		// Assert
		require.NoError(t, err, "Expected no error when deleting old recipes with a deleted recipe that is old enough")
		require.Equal(t, int64(1), deletedCount, "Expected to delete 1 old recipe when there is one deleted recipe that is old enough")
	})
}
