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

		s := NewService(txRunner, NewRepo(), ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100))
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

		s := NewService(txRunner, NewRepo(), ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100))
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
		s := NewService(txRunner, NewRepo(), ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100))
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
		s := NewService(txRunner, NewRepo(), ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100))
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
		s := NewService(txRunner, NewRepo(), ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100))
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

		s := NewService(txRunner, NewRepo(), ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100))
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

		s := NewService(txRunner, NewRepo(), ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100))
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
		deletedRecipes, err := s.GetDeletedRecipesByUserID(ctx, testUser.ID.String())

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
		deletedRecipes, err = s.GetDeletedRecipesByUserID(ctx, testUser.ID.String())

		// Assert
		require.NoError(t, err, "Expected no error when getting deleted recipes by user ID")
		require.Len(t, deletedRecipes, 0, "Expected to find no deleted recipes for user after undeleting")
	})
}
