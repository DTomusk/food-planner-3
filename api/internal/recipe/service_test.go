package recipe

import (
	"context"
	"database/sql"
	"foodplanner/internal/db"
	"foodplanner/internal/events"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/testutil"
	"foodplanner/internal/testutil/seeds"
	"foodplanner/internal/upload"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func newTestRecipeService(t *testing.T, tx *sql.Tx, txRunner db.TxRunner, repo *recipeRepo, uploadService *upload.UploadService) *Service {
	t.Helper()

	effectiveUploadService := uploadService
	if effectiveUploadService == nil {
		effectiveUploadService = upload.NewUploadServiceWithProvider(tx, upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com"), 0, upload.NewUploadRepo())
	}

	eventBus := events.NewInMemoryEventBus(1, 32, txRunner)
	t.Cleanup(func() {
		_ = eventBus.Close(context.Background())
	})

	return NewService(
		txRunner,
		repo,
		NewRecipeVersionRepo(),
		ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
		NewIngredientUsageRepo(),
		nil,
		effectiveUploadService,
		eventBus,
	)
}

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

		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, nil)
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

func TestCreateRecipe_PersistsDescription(t *testing.T) {
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
		require.NoError(t, err)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, nil)

		description := "Silky vanilla ice cream with a soft-set custard base."
		request := CreateRecipeRequest{
			Name: "Vanilla Ice Cream",
			Ingredients: []CreateIngredientUsageRequest{
				{
					IngredientID: ingredientID.String(),
					Quantity:     200,
					Unit:         1,
				},
			},
			UserID:      testUser.ID.String(),
			Description: description,
			PrepMins:    15,
			CookMins:    0,
			Portions:    6,
		}

		recipeContainer, err := s.CreateRecipe(ctx, request)
		require.NoError(t, err)
		require.NotNil(t, recipeContainer)
		require.Equal(t, description, recipeContainer.CurrentVersion.Description)

		var persistedDescription string
		err = tx.QueryRowContext(ctx, `SELECT description FROM recipe_versions WHERE id = $1`, recipeContainer.CurrentVersionID).Scan(&persistedDescription)
		require.NoError(t, err)
		require.Equal(t, description, persistedDescription)

		retrievedVersion, err := s.GetRecipeVersionByID(ctx, recipeContainer.CurrentVersionID)
		require.NoError(t, err)
		require.NotNil(t, retrievedVersion)
		require.Equal(t, description, retrievedVersion.Description)
	})
}

func TestCreateRecipe_SetsAnimalProductLevelFromIngredients(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()

		veganIngredientID := uuid.New()
		veganIngredient := ingredient.Ingredient{
			ID:                 veganIngredientID,
			FileKey:            "vegan_ingredient",
			Name:               "Tofu",
			PreferredUnit:      1,
			AnimalProductLevel: ingredient.Vegan,
		}
		err := seeds.InsertIngredient(ctx, tx, &veganIngredient)
		require.NoError(t, err)

		meatIngredientID := uuid.New()
		meatIngredient := ingredient.Ingredient{
			ID:                 meatIngredientID,
			FileKey:            "meat_ingredient",
			Name:               "Chicken",
			PreferredUnit:      1,
			AnimalProductLevel: ingredient.Meat,
		}
		err = seeds.InsertIngredient(ctx, tx, &meatIngredient)
		require.NoError(t, err)

		_, err = tx.ExecContext(ctx, `UPDATE reference.ingredients SET animal_product_level = $1 WHERE id = $2`, int(ingredient.Vegan), veganIngredientID)
		require.NoError(t, err)
		_, err = tx.ExecContext(ctx, `UPDATE reference.ingredients SET animal_product_level = $1 WHERE id = $2`, int(ingredient.Meat), meatIngredientID)
		require.NoError(t, err)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, nil)

		request := CreateRecipeRequest{
			Name: "Mixed Dish",
			Ingredients: []CreateIngredientUsageRequest{
				{IngredientID: veganIngredientID.String(), Quantity: 1, Unit: 1},
				{IngredientID: meatIngredientID.String(), Quantity: 1, Unit: 1},
			},
			UserID:   testUser.ID.String(),
			PrepMins: 10,
			CookMins: 15,
			Portions: 2,
		}

		recipeContainer, err := s.CreateRecipe(ctx, request)
		require.NoError(t, err)
		require.NotNil(t, recipeContainer)
		require.Equal(t, int(ingredient.Meat), recipeContainer.CurrentVersion.AnimalProductLevel)

		var persistedAnimalProductLevel int
		err = tx.QueryRowContext(ctx, `SELECT animal_product_level FROM recipe_versions WHERE id = $1`, recipeContainer.CurrentVersionID).Scan(&persistedAnimalProductLevel)
		require.NoError(t, err)
		require.Equal(t, int(ingredient.Meat), persistedAnimalProductLevel)
	})
}

func TestCreateRecipe_WithImageUploadID_PersistsAndRetrievesImageURL(t *testing.T) {
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

		uploadService := upload.NewUploadServiceWithProvider(
			tx,
			upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com"),
			0,
			upload.NewUploadRepo(),
		)

		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, uploadService)

		uploadRes, err := uploadService.CreateImageUploadURL(ctx, upload.CreateImageUploadURLRequest{
			OwnerUserID:   testUser.ID,
			FileName:      "dish.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       upload.UploadPurposeRecipeImage,
		})
		require.NoError(t, err)
		require.NotNil(t, uploadRes)

		uploadID := uploadRes.UploadID.String()
		request := CreateRecipeRequest{
			Name: "Vanilla Ice Cream",
			Ingredients: []CreateIngredientUsageRequest{
				{
					IngredientID: ingredientID.String(),
					Quantity:     200,
					Unit:         1,
				},
			},
			UserID:      testUser.ID.String(),
			PrepMins:    15,
			CookMins:    0,
			Portions:    6,
			ImgUploadID: &uploadID,
		}

		recipeContainer, err := s.CreateRecipe(ctx, request)
		require.NoError(t, err)
		require.NotNil(t, recipeContainer)

		var persistedImgSrc sql.NullString
		err = tx.QueryRowContext(ctx, `SELECT img_src FROM recipe_versions WHERE id = $1`, recipeContainer.CurrentVersionID).Scan(&persistedImgSrc)
		require.NoError(t, err)
		require.True(t, persistedImgSrc.Valid)
		require.Equal(t, uploadRes.FileURL, persistedImgSrc.String)

		retrievedByID, err := s.GetRecipeVersionByID(ctx, recipeContainer.CurrentVersionID)
		require.NoError(t, err)
		require.NotNil(t, retrievedByID)
		require.NotNil(t, retrievedByID.ImgSrc)
		require.Equal(t, uploadRes.FileURL, *retrievedByID.ImgSrc)

		retrievedByVersion, err := s.GetRecipeVersionByRecipeIDAndVersion(ctx, recipeContainer.ID, 1)
		require.NoError(t, err)
		require.NotNil(t, retrievedByVersion)
		require.NotNil(t, retrievedByVersion.ImgSrc)
		require.Equal(t, uploadRes.FileURL, *retrievedByVersion.ImgSrc)

		var usedAt sql.NullTime
		var linkedEntityID sql.NullString
		var linkedEntityType sql.NullString
		err = tx.QueryRowContext(
			ctx,
			`SELECT used_at, linked_entity_id, linked_entity_type FROM uploads WHERE id = $1`,
			uploadRes.UploadID,
		).Scan(&usedAt, &linkedEntityID, &linkedEntityType)
		require.NoError(t, err)
		require.True(t, usedAt.Valid)
		require.True(t, linkedEntityID.Valid)
		require.Equal(t, recipeContainer.CurrentVersionID.String(), linkedEntityID.String)
		require.True(t, linkedEntityType.Valid)
		require.Equal(t, "recipe-version", linkedEntityType.String)
	})
}

func TestCreateRecipe_WithUsedImageUploadID_ReturnsErrorAndDoesNotCreateRecipe(t *testing.T) {
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

		uploadService := upload.NewUploadServiceWithProvider(
			tx,
			upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com"),
			0,
			upload.NewUploadRepo(),
		)

		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, uploadService)

		uploadRes, err := uploadService.CreateImageUploadURL(ctx, upload.CreateImageUploadURLRequest{
			OwnerUserID:   testUser.ID,
			FileName:      "dish.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       upload.UploadPurposeRecipeImage,
		})
		require.NoError(t, err)

		err = uploadService.MarkUploadAsUsed(ctx, tx, upload.ClaimUploadRequest{
			UploadID:    uploadRes.UploadID,
			OwnerUserID: testUser.ID,
			EntityID:    uuid.New(),
			EntityType:  "recipe-version",
		})
		require.NoError(t, err)

		uploadID := uploadRes.UploadID.String()
		request := CreateRecipeRequest{
			Name: "Vanilla Ice Cream",
			Ingredients: []CreateIngredientUsageRequest{
				{
					IngredientID: ingredientID.String(),
					Quantity:     200,
					Unit:         1,
				},
			},
			UserID:      testUser.ID.String(),
			PrepMins:    15,
			CookMins:    0,
			Portions:    6,
			ImgUploadID: &uploadID,
		}

		recipeContainer, err := s.CreateRecipe(ctx, request)
		require.Error(t, err)
		require.ErrorIs(t, err, upload.ErrUploadAlreadyUsed)
		require.Nil(t, recipeContainer)

		var recipeCount int
		err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM recipe_containers WHERE user_id = $1`, testUser.ID).Scan(&recipeCount)
		require.NoError(t, err)
		require.Equal(t, 0, recipeCount)
	})
}

func TestUpdateRecipe_WithImageUploadID_PersistsAndRetrievesImageURL(t *testing.T) {
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

		uploadService := upload.NewUploadServiceWithProvider(
			tx,
			upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com"),
			0,
			upload.NewUploadRepo(),
		)

		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, uploadService)

		baseRequest := CreateRecipeRequest{
			Name: "Vanilla Ice Cream",
			Ingredients: []CreateIngredientUsageRequest{
				{
					IngredientID: ingredientID.String(),
					Quantity:     200,
					Unit:         1,
				},
			},
			UserID:   testUser.ID.String(),
			PrepMins: 15,
			CookMins: 0,
			Portions: 6,
		}

		recipeContainer, err := s.CreateRecipe(ctx, baseRequest)
		require.NoError(t, err)
		require.NotNil(t, recipeContainer)

		uploadRes, err := uploadService.CreateImageUploadURL(ctx, upload.CreateImageUploadURLRequest{
			OwnerUserID:   testUser.ID,
			FileName:      "updated-dish.png",
			FileType:      "image/png",
			FileSizeBytes: 2048,
			Purpose:       upload.UploadPurposeRecipeImage,
		})
		require.NoError(t, err)
		require.NotNil(t, uploadRes)

		uploadID := uploadRes.UploadID.String()
		updateReq := UpdateRecipeRequest{
			RecipeId: recipeContainer.ID.String(),
			Request: CreateRecipeRequest{
				Name: "Vanilla Ice Cream Updated",
				Ingredients: []CreateIngredientUsageRequest{
					{
						IngredientID: ingredientID.String(),
						Quantity:     200,
						Unit:         1,
					},
				},
				UserID:      testUser.ID.String(),
				PrepMins:    20,
				CookMins:    0,
				Portions:    6,
				ImgUploadID: &uploadID,
			},
		}

		updatedRecipe, err := s.UpdateRecipe(ctx, updateReq)
		require.NoError(t, err)
		require.NotNil(t, updatedRecipe)

		var persistedImgSrc sql.NullString
		err = tx.QueryRowContext(ctx, `SELECT img_src FROM recipe_versions WHERE id = $1`, updatedRecipe.CurrentVersionID).Scan(&persistedImgSrc)
		require.NoError(t, err)
		require.True(t, persistedImgSrc.Valid)
		require.Equal(t, uploadRes.FileURL, persistedImgSrc.String)

		retrievedByID, err := s.GetRecipeVersionByID(ctx, updatedRecipe.CurrentVersionID)
		require.NoError(t, err)
		require.NotNil(t, retrievedByID)
		require.NotNil(t, retrievedByID.ImgSrc)
		require.Equal(t, uploadRes.FileURL, *retrievedByID.ImgSrc)

		retrievedByVersion, err := s.GetRecipeVersionByRecipeIDAndVersion(ctx, updatedRecipe.ID, 2)
		require.NoError(t, err)
		require.NotNil(t, retrievedByVersion)
		require.NotNil(t, retrievedByVersion.ImgSrc)
		require.Equal(t, uploadRes.FileURL, *retrievedByVersion.ImgSrc)

		var usedAt sql.NullTime
		var linkedEntityID sql.NullString
		var linkedEntityType sql.NullString
		err = tx.QueryRowContext(
			ctx,
			`SELECT used_at, linked_entity_id, linked_entity_type FROM uploads WHERE id = $1`,
			uploadRes.UploadID,
		).Scan(&usedAt, &linkedEntityID, &linkedEntityType)
		require.NoError(t, err)
		require.True(t, usedAt.Valid)
		require.True(t, linkedEntityID.Valid)
		require.Equal(t, updatedRecipe.CurrentVersionID.String(), linkedEntityID.String)
		require.True(t, linkedEntityType.Valid)
		require.Equal(t, "recipe-version", linkedEntityType.String)
	})
}

func TestUpdateRecipe_PersistsDescriptionOnNewVersion(t *testing.T) {
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
		require.NoError(t, err)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, nil)

		createDescription := "Original vanilla ice cream description."
		baseRequest := CreateRecipeRequest{
			Name: "Vanilla Ice Cream",
			Ingredients: []CreateIngredientUsageRequest{
				{
					IngredientID: ingredientID.String(),
					Quantity:     200,
					Unit:         1,
				},
			},
			UserID:      testUser.ID.String(),
			Description: createDescription,
			PrepMins:    15,
			CookMins:    0,
			Portions:    6,
		}

		recipeContainer, err := s.CreateRecipe(ctx, baseRequest)
		require.NoError(t, err)
		require.NotNil(t, recipeContainer)

		updatedDescription := "Updated description with more detail about the new version."
		updateReq := UpdateRecipeRequest{
			RecipeId: recipeContainer.ID.String(),
			Request: CreateRecipeRequest{
				Name: "Vanilla Ice Cream Updated",
				Ingredients: []CreateIngredientUsageRequest{
					{
						IngredientID: ingredientID.String(),
						Quantity:     200,
						Unit:         1,
					},
				},
				UserID:      testUser.ID.String(),
				Description: updatedDescription,
				PrepMins:    20,
				CookMins:    0,
				Portions:    6,
			},
		}

		updatedRecipe, err := s.UpdateRecipe(ctx, updateReq)
		require.NoError(t, err)
		require.NotNil(t, updatedRecipe)
		require.Equal(t, updatedDescription, updatedRecipe.CurrentVersion.Description)

		var persistedDescription string
		err = tx.QueryRowContext(ctx, `SELECT description FROM recipe_versions WHERE id = $1`, updatedRecipe.CurrentVersionID).Scan(&persistedDescription)
		require.NoError(t, err)
		require.Equal(t, updatedDescription, persistedDescription)

		retrievedVersion, err := s.GetRecipeVersionByRecipeIDAndVersion(ctx, updatedRecipe.ID, 2)
		require.NoError(t, err)
		require.NotNil(t, retrievedVersion)
		require.Equal(t, updatedDescription, retrievedVersion.Description)

		originalVersion, err := s.GetRecipeVersionByRecipeIDAndVersion(ctx, updatedRecipe.ID, 1)
		require.NoError(t, err)
		require.NotNil(t, originalVersion)
		require.Equal(t, createDescription, originalVersion.Description)
	})
}

func TestUpdateRecipe_SetsAnimalProductLevelOnNewVersionFromIngredients(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()

		veganIngredientID := uuid.New()
		veganIngredient := ingredient.Ingredient{
			ID:                 veganIngredientID,
			FileKey:            "vegan_ingredient_update",
			Name:               "Spinach",
			PreferredUnit:      1,
			AnimalProductLevel: ingredient.Vegan,
		}
		err := seeds.InsertIngredient(ctx, tx, &veganIngredient)
		require.NoError(t, err)

		vegetarianIngredientID := uuid.New()
		vegetarianIngredient := ingredient.Ingredient{
			ID:                 vegetarianIngredientID,
			FileKey:            "vegetarian_ingredient_update",
			Name:               "Cheese",
			PreferredUnit:      1,
			AnimalProductLevel: ingredient.Vegetarian,
		}
		err = seeds.InsertIngredient(ctx, tx, &vegetarianIngredient)
		require.NoError(t, err)

		_, err = tx.ExecContext(ctx, `UPDATE reference.ingredients SET animal_product_level = $1 WHERE id = $2`, int(ingredient.Vegan), veganIngredientID)
		require.NoError(t, err)
		_, err = tx.ExecContext(ctx, `UPDATE reference.ingredients SET animal_product_level = $1 WHERE id = $2`, int(ingredient.Vegetarian), vegetarianIngredientID)
		require.NoError(t, err)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, nil)

		createReq := CreateRecipeRequest{
			Name: "Salad",
			Ingredients: []CreateIngredientUsageRequest{
				{IngredientID: veganIngredientID.String(), Quantity: 1, Unit: 1},
			},
			UserID:   testUser.ID.String(),
			PrepMins: 5,
			CookMins: 0,
			Portions: 1,
		}

		createdRecipe, err := s.CreateRecipe(ctx, createReq)
		require.NoError(t, err)
		require.NotNil(t, createdRecipe)
		require.Equal(t, int(ingredient.Vegan), createdRecipe.CurrentVersion.AnimalProductLevel)

		updateReq := UpdateRecipeRequest{
			RecipeId: createdRecipe.ID.String(),
			Request: CreateRecipeRequest{
				Name: "Salad With Cheese",
				Ingredients: []CreateIngredientUsageRequest{
					{IngredientID: veganIngredientID.String(), Quantity: 1, Unit: 1},
					{IngredientID: vegetarianIngredientID.String(), Quantity: 1, Unit: 1},
				},
				UserID:   testUser.ID.String(),
				PrepMins: 5,
				CookMins: 0,
				Portions: 1,
			},
		}

		updatedRecipe, err := s.UpdateRecipe(ctx, updateReq)
		require.NoError(t, err)
		require.NotNil(t, updatedRecipe)
		require.Equal(t, int(ingredient.Vegetarian), updatedRecipe.CurrentVersion.AnimalProductLevel)

		var persistedAnimalProductLevel int
		err = tx.QueryRowContext(ctx, `SELECT animal_product_level FROM recipe_versions WHERE id = $1`, updatedRecipe.CurrentVersionID).Scan(&persistedAnimalProductLevel)
		require.NoError(t, err)
		require.Equal(t, int(ingredient.Vegetarian), persistedAnimalProductLevel)

		originalVersion, err := s.GetRecipeVersionByRecipeIDAndVersion(ctx, updatedRecipe.ID, 1)
		require.NoError(t, err)
		require.NotNil(t, originalVersion)
		require.Equal(t, int(ingredient.Vegan), originalVersion.AnimalProductLevel)
	})
}

func TestCreateRecipe_SetsContainsGlutenFromIngredients(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()

		glutenFreeIngredientID := uuid.New()
		glutenFreeIngredient := ingredient.Ingredient{
			ID:            glutenFreeIngredientID,
			FileKey:       "gluten_free_ingredient",
			Name:          "Rice",
			PreferredUnit: 1,
		}
		err := seeds.InsertIngredient(ctx, tx, &glutenFreeIngredient)
		require.NoError(t, err)

		glutenIngredientID := uuid.New()
		glutenIngredient := ingredient.Ingredient{
			ID:            glutenIngredientID,
			FileKey:       "gluten_ingredient",
			Name:          "Flour",
			PreferredUnit: 1,
		}
		err = seeds.InsertIngredient(ctx, tx, &glutenIngredient)
		require.NoError(t, err)

		_, err = tx.ExecContext(ctx, `UPDATE reference.ingredients SET contains_gluten = $1 WHERE id = $2`, false, glutenFreeIngredientID)
		require.NoError(t, err)
		_, err = tx.ExecContext(ctx, `UPDATE reference.ingredients SET contains_gluten = $1 WHERE id = $2`, true, glutenIngredientID)
		require.NoError(t, err)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, nil)

		request := CreateRecipeRequest{
			Name: "Gluten Dish",
			Ingredients: []CreateIngredientUsageRequest{
				{IngredientID: glutenFreeIngredientID.String(), Quantity: 1, Unit: 1},
				{IngredientID: glutenIngredientID.String(), Quantity: 1, Unit: 1},
			},
			UserID:   testUser.ID.String(),
			PrepMins: 10,
			CookMins: 15,
			Portions: 2,
		}

		recipeContainer, err := s.CreateRecipe(ctx, request)
		require.NoError(t, err)
		require.NotNil(t, recipeContainer)
		require.True(t, recipeContainer.CurrentVersion.ContainsGluten)

		var persistedContainsGluten bool
		err = tx.QueryRowContext(ctx, `SELECT contains_gluten FROM recipe_versions WHERE id = $1`, recipeContainer.CurrentVersionID).Scan(&persistedContainsGluten)
		require.NoError(t, err)
		require.True(t, persistedContainsGluten)
	})
}

func TestUpdateRecipe_SetsContainsGlutenOnNewVersionFromIngredients(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()

		glutenFreeIngredientID := uuid.New()
		glutenFreeIngredient := ingredient.Ingredient{
			ID:            glutenFreeIngredientID,
			FileKey:       "gluten_free_ingredient_update",
			Name:          "Rice",
			PreferredUnit: 1,
		}
		err := seeds.InsertIngredient(ctx, tx, &glutenFreeIngredient)
		require.NoError(t, err)

		glutenIngredientID := uuid.New()
		glutenIngredient := ingredient.Ingredient{
			ID:            glutenIngredientID,
			FileKey:       "gluten_ingredient_update",
			Name:          "Wheat",
			PreferredUnit: 1,
		}
		err = seeds.InsertIngredient(ctx, tx, &glutenIngredient)
		require.NoError(t, err)

		_, err = tx.ExecContext(ctx, `UPDATE reference.ingredients SET contains_gluten = $1 WHERE id = $2`, false, glutenFreeIngredientID)
		require.NoError(t, err)
		_, err = tx.ExecContext(ctx, `UPDATE reference.ingredients SET contains_gluten = $1 WHERE id = $2`, true, glutenIngredientID)
		require.NoError(t, err)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, nil)

		createReq := CreateRecipeRequest{
			Name: "Gluten Free Dish",
			Ingredients: []CreateIngredientUsageRequest{
				{IngredientID: glutenFreeIngredientID.String(), Quantity: 1, Unit: 1},
			},
			UserID:   testUser.ID.String(),
			PrepMins: 5,
			CookMins: 10,
			Portions: 1,
		}

		createdRecipe, err := s.CreateRecipe(ctx, createReq)
		require.NoError(t, err)
		require.NotNil(t, createdRecipe)
		require.False(t, createdRecipe.CurrentVersion.ContainsGluten)

		updateReq := UpdateRecipeRequest{
			RecipeId: createdRecipe.ID.String(),
			Request: CreateRecipeRequest{
				Name: "Now With Wheat",
				Ingredients: []CreateIngredientUsageRequest{
					{IngredientID: glutenFreeIngredientID.String(), Quantity: 1, Unit: 1},
					{IngredientID: glutenIngredientID.String(), Quantity: 1, Unit: 1},
				},
				UserID:   testUser.ID.String(),
				PrepMins: 5,
				CookMins: 10,
				Portions: 1,
			},
		}

		updatedRecipe, err := s.UpdateRecipe(ctx, updateReq)
		require.NoError(t, err)
		require.NotNil(t, updatedRecipe)
		require.True(t, updatedRecipe.CurrentVersion.ContainsGluten)

		var persistedContainsGluten bool
		err = tx.QueryRowContext(ctx, `SELECT contains_gluten FROM recipe_versions WHERE id = $1`, updatedRecipe.CurrentVersionID).Scan(&persistedContainsGluten)
		require.NoError(t, err)
		require.True(t, persistedContainsGluten)
	})
}

func TestUpdateRecipe_WithUsedImageUploadID_ReturnsErrorAndDoesNotCreateNewVersion(t *testing.T) {
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

		uploadService := upload.NewUploadServiceWithProvider(
			tx,
			upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com"),
			0,
			upload.NewUploadRepo(),
		)

		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, uploadService)

		baseRequest := CreateRecipeRequest{
			Name: "Vanilla Ice Cream",
			Ingredients: []CreateIngredientUsageRequest{
				{
					IngredientID: ingredientID.String(),
					Quantity:     200,
					Unit:         1,
				},
			},
			UserID:   testUser.ID.String(),
			PrepMins: 15,
			CookMins: 0,
			Portions: 6,
		}

		recipeContainer, err := s.CreateRecipe(ctx, baseRequest)
		require.NoError(t, err)
		require.NotNil(t, recipeContainer)
		originalCurrentVersionID := recipeContainer.CurrentVersionID

		uploadRes, err := uploadService.CreateImageUploadURL(ctx, upload.CreateImageUploadURLRequest{
			OwnerUserID:   testUser.ID,
			FileName:      "already-used.png",
			FileType:      "image/png",
			FileSizeBytes: 2048,
			Purpose:       upload.UploadPurposeRecipeImage,
		})
		require.NoError(t, err)

		err = uploadService.MarkUploadAsUsed(ctx, tx, upload.ClaimUploadRequest{
			UploadID:    uploadRes.UploadID,
			OwnerUserID: testUser.ID,
			EntityID:    uuid.New(),
			EntityType:  "recipe-version",
		})
		require.NoError(t, err)

		uploadID := uploadRes.UploadID.String()
		updateReq := UpdateRecipeRequest{
			RecipeId: recipeContainer.ID.String(),
			Request: CreateRecipeRequest{
				Name: "Vanilla Ice Cream Updated",
				Ingredients: []CreateIngredientUsageRequest{
					{
						IngredientID: ingredientID.String(),
						Quantity:     220,
						Unit:         1,
					},
				},
				UserID:      testUser.ID.String(),
				PrepMins:    20,
				CookMins:    0,
				Portions:    6,
				ImgUploadID: &uploadID,
			},
		}

		updatedRecipe, err := s.UpdateRecipe(ctx, updateReq)
		require.Error(t, err)
		require.ErrorIs(t, err, upload.ErrUploadAlreadyUsed)
		require.Nil(t, updatedRecipe)

		var currentVersionID uuid.UUID
		err = tx.QueryRowContext(ctx, `SELECT current_version_id FROM recipe_containers WHERE id = $1`, recipeContainer.ID).Scan(&currentVersionID)
		require.NoError(t, err)
		require.Equal(t, originalCurrentVersionID, currentVersionID)

		var versionCount int
		err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM recipe_versions WHERE recipe_id = $1`, recipeContainer.ID).Scan(&versionCount)
		require.NoError(t, err)
		require.Equal(t, 1, versionCount)
	})
}

func TestUpdateRecipe_WithoutImageChanges_PreservesImageFromPreviousVersion(t *testing.T) {
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

		uploadService := upload.NewUploadServiceWithProvider(
			tx,
			upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com"),
			0,
			upload.NewUploadRepo(),
		)

		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, uploadService)

		// Create recipe with initial image
		initialUploadRes, err := uploadService.CreateImageUploadURL(ctx, upload.CreateImageUploadURLRequest{
			OwnerUserID:   testUser.ID,
			FileName:      "initial-dish.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       upload.UploadPurposeRecipeImage,
		})
		require.NoError(t, err)
		require.NotNil(t, initialUploadRes)

		initialUploadID := initialUploadRes.UploadID.String()
		createReq := CreateRecipeRequest{
			Name: "Chocolate Cake",
			Ingredients: []CreateIngredientUsageRequest{
				{
					IngredientID: ingredientID.String(),
					Quantity:     200,
					Unit:         1,
				},
			},
			UserID:      testUser.ID.String(),
			PrepMins:    30,
			CookMins:    45,
			Portions:    8,
			ImgUploadID: &initialUploadID,
		}

		recipeContainer, err := s.CreateRecipe(ctx, createReq)
		require.NoError(t, err)
		require.NotNil(t, recipeContainer)
		require.NotNil(t, recipeContainer.CurrentVersion.ImgSrc)
		require.Equal(t, initialUploadRes.FileURL, *recipeContainer.CurrentVersion.ImgSrc)

		// Update recipe without providing new image ID and without removeImage flag
		// Expected: image should be preserved from previous version
		updateReq := UpdateRecipeRequest{
			RecipeId: recipeContainer.ID.String(),
			Request: CreateRecipeRequest{
				Name: "Chocolate Cake Updated",
				Ingredients: []CreateIngredientUsageRequest{
					{
						IngredientID: ingredientID.String(),
						Quantity:     250,
						Unit:         1,
					},
				},
				UserID:      testUser.ID.String(),
				PrepMins:    35,
				CookMins:    45,
				Portions:    10,
				ImgUploadID: nil,
			},
			RemoveImage: nil,
		}

		updatedRecipe, err := s.UpdateRecipe(ctx, updateReq)
		require.NoError(t, err)
		require.NotNil(t, updatedRecipe)

		// Verify image is preserved in new version
		require.NotNil(t, updatedRecipe.CurrentVersion.ImgSrc)
		require.Equal(t, initialUploadRes.FileURL, *updatedRecipe.CurrentVersion.ImgSrc)

		// Verify in database
		var persistedImgSrc sql.NullString
		err = tx.QueryRowContext(ctx, `SELECT img_src FROM recipe_versions WHERE id = $1`, updatedRecipe.CurrentVersionID).Scan(&persistedImgSrc)
		require.NoError(t, err)
		require.True(t, persistedImgSrc.Valid)
		require.Equal(t, initialUploadRes.FileURL, persistedImgSrc.String)

		// Verify via service retrieval
		retrievedVersion, err := s.GetRecipeVersionByRecipeIDAndVersion(ctx, updatedRecipe.ID, 2)
		require.NoError(t, err)
		require.NotNil(t, retrievedVersion)
		require.NotNil(t, retrievedVersion.ImgSrc)
		require.Equal(t, initialUploadRes.FileURL, *retrievedVersion.ImgSrc)
	})
}

func TestUpdateRecipe_WithRemoveImageFlag_RemovesImage(t *testing.T) {
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

		uploadService := upload.NewUploadServiceWithProvider(
			tx,
			upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com"),
			0,
			upload.NewUploadRepo(),
		)

		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, uploadService)

		// Create recipe with initial image
		initialUploadRes, err := uploadService.CreateImageUploadURL(ctx, upload.CreateImageUploadURLRequest{
			OwnerUserID:   testUser.ID,
			FileName:      "initial-dish.png",
			FileType:      "image/png",
			FileSizeBytes: 1024,
			Purpose:       upload.UploadPurposeRecipeImage,
		})
		require.NoError(t, err)
		require.NotNil(t, initialUploadRes)

		initialUploadID := initialUploadRes.UploadID.String()
		createReq := CreateRecipeRequest{
			Name: "Pancakes",
			Ingredients: []CreateIngredientUsageRequest{
				{
					IngredientID: ingredientID.String(),
					Quantity:     100,
					Unit:         1,
				},
			},
			UserID:      testUser.ID.String(),
			PrepMins:    10,
			CookMins:    15,
			Portions:    4,
			ImgUploadID: &initialUploadID,
		}

		recipeContainer, err := s.CreateRecipe(ctx, createReq)
		require.NoError(t, err)
		require.NotNil(t, recipeContainer)
		require.NotNil(t, recipeContainer.CurrentVersion.ImgSrc)

		// Update recipe with removeImage flag set to true
		// Expected: image should be removed (imgSrc = nil)
		trueFlag := true
		updateReq := UpdateRecipeRequest{
			RecipeId: recipeContainer.ID.String(),
			Request: CreateRecipeRequest{
				Name: "Pancakes Updated",
				Ingredients: []CreateIngredientUsageRequest{
					{
						IngredientID: ingredientID.String(),
						Quantity:     100,
						Unit:         1,
					},
				},
				UserID:      testUser.ID.String(),
				PrepMins:    10,
				CookMins:    20,
				Portions:    4,
				ImgUploadID: nil,
			},
			RemoveImage: &trueFlag,
		}

		updatedRecipe, err := s.UpdateRecipe(ctx, updateReq)
		require.NoError(t, err)
		require.NotNil(t, updatedRecipe)

		// Verify image is removed in new version
		require.Nil(t, updatedRecipe.CurrentVersion.ImgSrc)

		// Verify in database
		var persistedImgSrc sql.NullString
		err = tx.QueryRowContext(ctx, `SELECT img_src FROM recipe_versions WHERE id = $1`, updatedRecipe.CurrentVersionID).Scan(&persistedImgSrc)
		require.NoError(t, err)
		require.False(t, persistedImgSrc.Valid, "Expected img_src to be NULL in database")

		// Verify via service retrieval
		retrievedVersion, err := s.GetRecipeVersionByRecipeIDAndVersion(ctx, updatedRecipe.ID, 2)
		require.NoError(t, err)
		require.NotNil(t, retrievedVersion)
		require.Nil(t, retrievedVersion.ImgSrc)
	})
}

func TestUpdateRecipe_RecipeWithoutImage_RemainsWithoutImage(t *testing.T) {
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

		uploadService := upload.NewUploadServiceWithProvider(
			tx,
			upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com"),
			0,
			upload.NewUploadRepo(),
		)

		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, uploadService)

		// Create recipe without image
		createReq := CreateRecipeRequest{
			Name: "Simple Salad",
			Ingredients: []CreateIngredientUsageRequest{
				{
					IngredientID: ingredientID.String(),
					Quantity:     150,
					Unit:         1,
				},
			},
			UserID:      testUser.ID.String(),
			PrepMins:    5,
			CookMins:    0,
			Portions:    2,
			ImgUploadID: nil,
		}

		recipeContainer, err := s.CreateRecipe(ctx, createReq)
		require.NoError(t, err)
		require.NotNil(t, recipeContainer)
		require.Nil(t, recipeContainer.CurrentVersion.ImgSrc)

		// Update recipe without providing image ID and without removeImage flag
		// Expected: should remain without image
		updateReq := UpdateRecipeRequest{
			RecipeId: recipeContainer.ID.String(),
			Request: CreateRecipeRequest{
				Name: "Simple Salad Updated",
				Ingredients: []CreateIngredientUsageRequest{
					{
						IngredientID: ingredientID.String(),
						Quantity:     200,
						Unit:         1,
					},
				},
				UserID:      testUser.ID.String(),
				PrepMins:    5,
				CookMins:    0,
				Portions:    3,
				ImgUploadID: nil,
			},
			RemoveImage: nil,
		}

		updatedRecipe, err := s.UpdateRecipe(ctx, updateReq)
		require.NoError(t, err)
		require.NotNil(t, updatedRecipe)

		// Verify image is still nil
		require.Nil(t, updatedRecipe.CurrentVersion.ImgSrc)

		// Verify in database
		var persistedImgSrc sql.NullString
		err = tx.QueryRowContext(ctx, `SELECT img_src FROM recipe_versions WHERE id = $1`, updatedRecipe.CurrentVersionID).Scan(&persistedImgSrc)
		require.NoError(t, err)
		require.False(t, persistedImgSrc.Valid, "Expected img_src to be NULL in database")

		// Verify via service retrieval
		retrievedVersion, err := s.GetRecipeVersionByRecipeIDAndVersion(ctx, updatedRecipe.ID, 2)
		require.NoError(t, err)
		require.NotNil(t, retrievedVersion)
		require.Nil(t, retrievedVersion.ImgSrc)
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

		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, nil)

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
		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, nil)

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
		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, nil)
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
		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, nil)
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

		repo, err := NewRecipeRepo(0.15, 0.85)

		s := newTestRecipeService(t, tx, txRunner, repo, nil)
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
		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, nil)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		newestCreatedAt := time.Date(2026, time.March, 17, 11, 42, 48, 147630000, time.UTC)
		middleCreatedAt := newestCreatedAt.Add(-1 * time.Minute)
		oldestCreatedAt := newestCreatedAt.Add(-2 * time.Minute)

		newest := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1"), uuid.New(), "Newest", newestCreatedAt, nil, nil)
		middle := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2"), uuid.New(), "Middle", middleCreatedAt, nil, nil)
		oldest := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("cccccccc-cccc-cccc-cccc-ccccccccccc3"), uuid.New(), "Oldest", oldestCreatedAt, nil, nil)

		params := RecipeListParams{
			Pagination: RecipePagination{
				First: 2,
				After: nil,
			},
			Filter: RecipeFilter{},
		}

		firstPage, nextCursor, err := s.GetRecipes(ctx, params)

		require.NoError(t, err)
		require.Len(t, firstPage, 2)
		require.Equal(t, newest.RecipeID, firstPage[0].Recipe.ID)
		require.Equal(t, middle.RecipeID, firstPage[1].Recipe.ID)

		firstEdgeCursor, err := ParseRecipeCursor(&firstPage[0].Cursor)
		require.NoError(t, err)
		require.NotNil(t, firstEdgeCursor)
		require.True(t, newest.CreatedAt.Equal(firstEdgeCursor.CreatedAt))
		require.Equal(t, newest.RecipeID, firstEdgeCursor.ID)

		secondEdgeCursor, err := ParseRecipeCursor(&firstPage[1].Cursor)
		require.NoError(t, err)
		require.NotNil(t, secondEdgeCursor)
		require.True(t, middle.CreatedAt.Equal(secondEdgeCursor.CreatedAt))
		require.Equal(t, middle.RecipeID, secondEdgeCursor.ID)
		require.NotNil(t, nextCursor)

		parsedCursor, err := ParseRecipeCursor(nextCursor)
		require.NoError(t, err)
		require.NotNil(t, parsedCursor)
		require.True(t, middle.CreatedAt.Equal(parsedCursor.CreatedAt))
		require.Equal(t, middle.RecipeID, parsedCursor.ID)

		params.Pagination.After = nextCursor
		secondPage, finalCursor, err := s.GetRecipes(ctx, params)

		require.NoError(t, err)
		require.Len(t, secondPage, 1)
		require.Equal(t, oldest.RecipeID, secondPage[0].Recipe.ID)
		require.Nil(t, finalCursor)
	})
}

func TestGetRecipes_InvalidCursor(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		txRunner := testutil.NewTestTxRunner(tx)
		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, nil)

		invalidCursor := "not-a-valid-cursor"

		params := RecipeListParams{
			Pagination: RecipePagination{
				First: 2,
				After: &invalidCursor,
			},
			Filter: RecipeFilter{},
		}

		recipes, nextCursor, err := s.GetRecipes(ctx, params)

		require.ErrorIs(t, err, ErrInvalidCursor)
		require.Nil(t, recipes)
		require.Nil(t, nextCursor)
	})
}

func TestGetRecipes_CursorIncludesModeAndFilterHashForNewest(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx, s, newest, middle, _ := setupRecipeListFixture(t, tx)

		params := RecipeListParams{
			Pagination: RecipePagination{First: 2},
			Filter:     RecipeFilter{},
		}

		recipes, nextCursor, err := s.GetRecipes(ctx, params)

		require.NoError(t, err)
		require.Len(t, recipes, 2)
		require.Equal(t, newest.RecipeID, recipes[0].Recipe.ID)
		require.Equal(t, middle.RecipeID, recipes[1].Recipe.ID)
		require.NotNil(t, nextCursor)

		edgeCursor, err := ParseRecipeCursor(&recipes[0].Cursor)
		require.NoError(t, err)
		require.NotNil(t, edgeCursor)
		require.Equal(t, RecipeCursorModeNewest, edgeCursor.Mode)
		require.Equal(t, filterHash(RecipeCursorModeNewest, nil, normalizedRecipeFilter{}), edgeCursor.FilterHash)
		require.Nil(t, edgeCursor.RelevanceScore)

		pageCursor, err := ParseRecipeCursor(nextCursor)
		require.NoError(t, err)
		require.NotNil(t, pageCursor)
		require.Equal(t, RecipeCursorModeNewest, pageCursor.Mode)
		require.Equal(t, filterHash(RecipeCursorModeNewest, nil, normalizedRecipeFilter{}), pageCursor.FilterHash)
	})
}

func TestGetRecipes_StaleCursorModeIsIgnored(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx, s, newest, middle, _ := setupRecipeListFixture(t, tx)

		searchQuery := "pasta"
		score := 0.8
		staleCursor := (&RecipeCursor{
			Mode:           RecipeCursorModeRelevance,
			FilterHash:     filterHash(RecipeCursorModeRelevance, &searchQuery, normalizedRecipeFilter{}),
			CreatedAt:      newest.CreatedAt,
			ID:             newest.RecipeID,
			RelevanceScore: &score,
		}).String()
		require.NotEmpty(t, staleCursor)

		params := RecipeListParams{
			Pagination: RecipePagination{First: 2, After: &staleCursor},
			Filter:     RecipeFilter{},
		}

		recipes, _, err := s.GetRecipes(ctx, params)

		require.NoError(t, err)
		require.Len(t, recipes, 2)
		// If stale cursor mode is ignored, results should start from the first page.
		require.Equal(t, newest.RecipeID, recipes[0].Recipe.ID)
		require.Equal(t, middle.RecipeID, recipes[1].Recipe.ID)
	})
}

func TestGetRecipes_StaleCursorFilterHashIsIgnored(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx, s, newest, middle, _ := setupRecipeListFixture(t, tx)

		staleCursor := (&RecipeCursor{
			Mode:       RecipeCursorModeNewest,
			FilterHash: "stale-filter-hash",
			CreatedAt:  newest.CreatedAt,
			ID:         newest.RecipeID,
		}).String()
		require.NotEmpty(t, staleCursor)

		params := RecipeListParams{
			Pagination: RecipePagination{First: 2, After: &staleCursor},
			Filter:     RecipeFilter{},
		}

		recipes, _, err := s.GetRecipes(ctx, params)

		require.NoError(t, err)
		require.Len(t, recipes, 2)
		// If stale cursor hash is ignored, results should start from the first page.
		require.Equal(t, newest.RecipeID, recipes[0].Recipe.ID)
		require.Equal(t, middle.RecipeID, recipes[1].Recipe.ID)
	})
}

func TestGetRecipes_ValidCursorWithMatchingHashAppliesBoundary(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx, s, newest, middle, oldest := setupRecipeListFixture(t, tx)

		validCursor := (&RecipeCursor{
			Mode:       RecipeCursorModeNewest,
			FilterHash: filterHash(RecipeCursorModeNewest, nil, normalizedRecipeFilter{}),
			CreatedAt:  newest.CreatedAt,
			ID:         newest.RecipeID,
		}).String()
		require.NotEmpty(t, validCursor)

		params := RecipeListParams{
			Pagination: RecipePagination{First: 2, After: &validCursor},
			Filter:     RecipeFilter{},
		}

		recipes, nextCursor, err := s.GetRecipes(ctx, params)

		require.NoError(t, err)
		require.Len(t, recipes, 2)
		require.Equal(t, middle.RecipeID, recipes[0].Recipe.ID)
		require.Equal(t, oldest.RecipeID, recipes[1].Recipe.ID)
		require.Nil(t, nextCursor)
	})
}

func TestGetRecipes_SearchQueryReturnsRelevanceCursorsAndPaginates(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx, s, exactHigh, exactLow, fuzzy := setupRecipeSearchFixture(t, tx)

		query := "chicken soup"
		params := RecipeListParams{
			Pagination: RecipePagination{First: 2},
			Filter: RecipeFilter{
				Query: &query,
			},
		}

		firstPage, nextCursor, err := s.GetRecipes(ctx, params)

		require.NoError(t, err)
		require.Len(t, firstPage, 2)
		require.NotNil(t, nextCursor)
		require.Equal(t, exactHigh.RecipeID, firstPage[0].Recipe.ID)
		require.Equal(t, exactLow.RecipeID, firstPage[1].Recipe.ID)

		expectedHash := filterHash(RecipeCursorModeRelevance, &query, normalizedRecipeFilter{})
		for _, edge := range firstPage {
			parsed, err := ParseRecipeCursor(&edge.Cursor)
			require.NoError(t, err)
			require.NotNil(t, parsed)
			require.Equal(t, RecipeCursorModeRelevance, parsed.Mode)
			require.Equal(t, expectedHash, parsed.FilterHash)
			require.NotNil(t, parsed.RelevanceScore)
		}

		parsedNext, err := ParseRecipeCursor(nextCursor)
		require.NoError(t, err)
		require.NotNil(t, parsedNext)
		require.Equal(t, RecipeCursorModeRelevance, parsedNext.Mode)
		require.Equal(t, expectedHash, parsedNext.FilterHash)
		require.NotNil(t, parsedNext.RelevanceScore)

		params.Pagination.After = nextCursor
		secondPage, finalCursor, err := s.GetRecipes(ctx, params)

		require.NoError(t, err)
		require.Len(t, secondPage, 1)
		require.Equal(t, fuzzy.RecipeID, secondPage[0].Recipe.ID)
		require.Nil(t, finalCursor)

		lastParsed, err := ParseRecipeCursor(&secondPage[0].Cursor)
		require.NoError(t, err)
		require.NotNil(t, lastParsed)
		require.Equal(t, RecipeCursorModeRelevance, lastParsed.Mode)
		require.Equal(t, expectedHash, lastParsed.FilterHash)
		require.NotNil(t, lastParsed.RelevanceScore)
	})
}

func TestGetRecipes_SearchQueryWithStaleNewestCursorIsIgnored(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx, s, exactHigh, exactLow, _ := setupRecipeSearchFixture(t, tx)

		staleCursor := (&RecipeCursor{
			Mode:       RecipeCursorModeNewest,
			FilterHash: filterHash(RecipeCursorModeNewest, nil, normalizedRecipeFilter{}),
			CreatedAt:  exactHigh.CreatedAt,
			ID:         exactHigh.RecipeID,
		}).String()
		require.NotEmpty(t, staleCursor)

		query := "chicken soup"
		params := RecipeListParams{
			Pagination: RecipePagination{First: 2, After: &staleCursor},
			Filter: RecipeFilter{
				Query: &query,
			},
		}

		recipes, _, err := s.GetRecipes(ctx, params)

		require.NoError(t, err)
		require.Len(t, recipes, 2)
		// If stale mode cursor is ignored, we start at the first page of search results.
		require.Equal(t, exactHigh.RecipeID, recipes[0].Recipe.ID)
		require.Equal(t, exactLow.RecipeID, recipes[1].Recipe.ID)
	})
}

func TestGetRecipes_SearchQueryWithStaleRelevanceHashIsIgnored(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx, s, exactHigh, exactLow, _ := setupRecipeSearchFixture(t, tx)

		otherQuery := "pasta"
		score := 0.9
		staleCursor := (&RecipeCursor{
			Mode:           RecipeCursorModeRelevance,
			FilterHash:     filterHash(RecipeCursorModeRelevance, &otherQuery, normalizedRecipeFilter{}),
			CreatedAt:      exactHigh.CreatedAt,
			ID:             exactHigh.RecipeID,
			RelevanceScore: &score,
		}).String()
		require.NotEmpty(t, staleCursor)

		query := "chicken soup"
		params := RecipeListParams{
			Pagination: RecipePagination{First: 2, After: &staleCursor},
			Filter: RecipeFilter{
				Query: &query,
			},
		}

		recipes, _, err := s.GetRecipes(ctx, params)

		require.NoError(t, err)
		require.Len(t, recipes, 2)
		// If stale hash is ignored, we start at the first page for the current search query.
		require.Equal(t, exactHigh.RecipeID, recipes[0].Recipe.ID)
		require.Equal(t, exactLow.RecipeID, recipes[1].Recipe.ID)
	})
}

func TestGetRecipes_FiltersByUserID(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		txRunner := testutil.NewTestTxRunner(tx)
		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, nil)

		userA, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)
		userB, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		base := time.Date(2026, time.March, 17, 11, 42, 48, 147630000, time.UTC)
		aNewest := seedRecipeForListTests(t, ctx, tx, userA.ID, uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1"), uuid.New(), "A Newest", base, nil, nil)
		aOlder := seedRecipeForListTests(t, ctx, tx, userA.ID, uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2"), uuid.New(), "A Older", base.Add(-1*time.Minute), nil, nil)
		_ = seedRecipeForListTests(t, ctx, tx, userB.ID, uuid.MustParse("cccccccc-cccc-cccc-cccc-ccccccccccc3"), uuid.New(), "B Recipe", base.Add(-2*time.Minute), nil, nil)

		userAID := userA.ID
		params := RecipeListParams{
			Pagination: RecipePagination{First: 10},
			Filter:     RecipeFilter{UserID: &userAID},
		}

		recipes, nextCursor, err := s.GetRecipes(ctx, params)
		require.NoError(t, err)
		require.Len(t, recipes, 2)
		require.Nil(t, nextCursor)
		require.Equal(t, aNewest.RecipeID, recipes[0].Recipe.ID)
		require.Equal(t, aOlder.RecipeID, recipes[1].Recipe.ID)

		for _, row := range recipes {
			require.Equal(t, userA.ID, row.Recipe.UserID)
		}

		parsed, err := ParseRecipeCursor(&recipes[0].Cursor)
		require.NoError(t, err)
		require.NotNil(t, parsed)
		require.Equal(t, filterHash(RecipeCursorModeNewest, nil, normalizedRecipeFilter{UserID: &userAID}), parsed.FilterHash)
	})
}

func TestGetRecipes_FiltersByAnimalProductLevel(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		txRunner := testutil.NewTestTxRunner(tx)
		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := newTestRecipeService(t, tx, txRunner, repo, nil)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		base := time.Date(2026, time.March, 17, 11, 42, 48, 147630000, time.UTC)
		veganLevel := 0
		vegetarianLevel := 1
		meatLevel := 2

		vegan := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("12345678-aaaa-aaaa-aaaa-aaaaaaaaaaa1"), uuid.New(), "Vegan Bowl", base, nil, &veganLevel)
		vegetarian := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("12345678-bbbb-bbbb-bbbb-bbbbbbbbbbb2"), uuid.New(), "Vegetarian Bowl", base.Add(-1*time.Minute), nil, &vegetarianLevel)
		meat := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("12345678-cccc-cccc-cccc-ccccccccccc3"), uuid.New(), "Meat Bowl", base.Add(-2*time.Minute), nil, &meatLevel)

		veganFilter := 0
		recipes, nextCursor, err := s.GetRecipes(ctx, RecipeListParams{
			Pagination: RecipePagination{First: 10},
			Filter:     RecipeFilter{AnimalProductLevel: &veganFilter},
		})
		require.NoError(t, err)
		require.Nil(t, nextCursor)
		require.Len(t, recipes, 1)
		require.Equal(t, vegan.RecipeID, recipes[0].Recipe.ID)

		vegetarianFilter := 1
		recipes, nextCursor, err = s.GetRecipes(ctx, RecipeListParams{
			Pagination: RecipePagination{First: 10},
			Filter:     RecipeFilter{AnimalProductLevel: &vegetarianFilter},
		})
		require.NoError(t, err)
		require.Nil(t, nextCursor)
		require.Len(t, recipes, 2)
		require.Equal(t, vegan.RecipeID, recipes[0].Recipe.ID)
		require.Equal(t, vegetarian.RecipeID, recipes[1].Recipe.ID)

		anyFilter := 2
		recipes, nextCursor, err = s.GetRecipes(ctx, RecipeListParams{
			Pagination: RecipePagination{First: 10},
			Filter:     RecipeFilter{AnimalProductLevel: &anyFilter},
		})
		require.NoError(t, err)
		require.Nil(t, nextCursor)
		require.Len(t, recipes, 3)
		require.Equal(t, vegan.RecipeID, recipes[0].Recipe.ID)
		require.Equal(t, vegetarian.RecipeID, recipes[1].Recipe.ID)
		require.Equal(t, meat.RecipeID, recipes[2].Recipe.ID)
	})
}

func TestGetRecipes_UnsupportedAnimalProductLevelUsesUnfilteredCursorHash(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx, s, newest, middle, _ := setupRecipeListFixture(t, tx)

		unsupportedAnimalLevel := 2
		recipes, _, err := s.GetRecipes(ctx, RecipeListParams{
			Pagination: RecipePagination{First: 2},
			Filter:     RecipeFilter{AnimalProductLevel: &unsupportedAnimalLevel},
		})
		require.NoError(t, err)
		require.Len(t, recipes, 2)
		require.Equal(t, newest.RecipeID, recipes[0].Recipe.ID)
		require.Equal(t, middle.RecipeID, recipes[1].Recipe.ID)

		parsed, err := ParseRecipeCursor(&recipes[0].Cursor)
		require.NoError(t, err)
		require.NotNil(t, parsed)
		require.Equal(t, filterHash(RecipeCursorModeNewest, nil, normalizedRecipeFilter{}), parsed.FilterHash)
	})
}

func TestGetRecipes_FiltersByContainsGluten(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx, s, newest, middle, _ := setupRecipeListFixture(t, tx)

		_, err := tx.ExecContext(ctx, `UPDATE recipe_versions SET contains_gluten = $1 WHERE id = $2`, false, newest.VersionID)
		require.NoError(t, err)
		_, err = tx.ExecContext(ctx, `UPDATE recipe_versions SET contains_gluten = $1 WHERE id = $2`, true, middle.VersionID)
		require.NoError(t, err)

		containsGluten := true
		recipes, nextCursor, err := s.GetRecipes(ctx, RecipeListParams{
			Pagination: RecipePagination{First: 10},
			Filter:     RecipeFilter{ContainsGluten: &containsGluten},
		})
		require.NoError(t, err)
		require.Nil(t, nextCursor)
		require.Len(t, recipes, 1)
		require.Equal(t, middle.RecipeID, recipes[0].Recipe.ID)

		glutenFree := false
		recipes, nextCursor, err = s.GetRecipes(ctx, RecipeListParams{
			Pagination: RecipePagination{First: 10},
			Filter:     RecipeFilter{ContainsGluten: &glutenFree},
		})
		require.NoError(t, err)
		require.Nil(t, nextCursor)
		require.Len(t, recipes, 2)
		require.Equal(t, newest.RecipeID, recipes[0].Recipe.ID)
	})
}

func TestGetRecipes_ContainsGlutenFilterAffectsCursorHash(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx, s, newest, _, _ := setupRecipeListFixture(t, tx)

		_, err := tx.ExecContext(ctx, `UPDATE recipe_versions SET contains_gluten = $1 WHERE id = $2`, true, newest.VersionID)
		require.NoError(t, err)

		containsGluten := true
		recipes, _, err := s.GetRecipes(ctx, RecipeListParams{
			Pagination: RecipePagination{First: 2},
			Filter:     RecipeFilter{ContainsGluten: &containsGluten},
		})
		require.NoError(t, err)
		require.NotEmpty(t, recipes)

		parsed, err := ParseRecipeCursor(&recipes[0].Cursor)
		require.NoError(t, err)
		require.NotNil(t, parsed)
		require.Equal(t, filterHash(RecipeCursorModeNewest, nil, normalizedRecipeFilter{ContainsGluten: &containsGluten}), parsed.FilterHash)
	})
}

func setupRecipeListFixture(t *testing.T, tx *sql.Tx) (context.Context, *Service, listedRecipeSeed, listedRecipeSeed, listedRecipeSeed) {
	t.Helper()

	ctx := context.Background()
	txRunner := testutil.NewTestTxRunner(tx)
	repo, err := NewRecipeRepo(0.15, 0.85)
	require.NoError(t, err)
	s := newTestRecipeService(t, tx, txRunner, repo, nil)

	testUser, err := seeds.SeedTestUser(ctx, tx)
	require.NoError(t, err, "Failed to seed test user")

	newestCreatedAt := time.Date(2026, time.March, 17, 11, 42, 48, 147630000, time.UTC)
	middleCreatedAt := newestCreatedAt.Add(-1 * time.Minute)
	oldestCreatedAt := newestCreatedAt.Add(-2 * time.Minute)

	newest := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1"), uuid.New(), "Newest", newestCreatedAt, nil, nil)
	middle := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2"), uuid.New(), "Middle", middleCreatedAt, nil, nil)
	oldest := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("cccccccc-cccc-cccc-cccc-ccccccccccc3"), uuid.New(), "Oldest", oldestCreatedAt, nil, nil)

	return ctx, s, newest, middle, oldest
}

func setupRecipeSearchFixture(t *testing.T, tx *sql.Tx) (context.Context, *Service, listedRecipeSeed, listedRecipeSeed, listedRecipeSeed) {
	t.Helper()

	ctx := context.Background()
	txRunner := testutil.NewTestTxRunner(tx)
	repo, err := NewRecipeRepo(0.15, 0.85)
	require.NoError(t, err)
	s := newTestRecipeService(t, tx, txRunner, repo, nil)

	testUser, err := seeds.SeedTestUser(ctx, tx)
	require.NoError(t, err, "Failed to seed test user")

	sameCreatedAt := time.Date(2026, time.March, 17, 11, 40, 48, 147630000, time.UTC)
	olderCreatedAt := sameCreatedAt.Add(-1 * time.Minute)

	exactHigh := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("ffffffff-ffff-ffff-ffff-fffffffffff1"), uuid.New(), "Chicken Soup", sameCreatedAt, nil, nil)
	exactLow := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("00000000-0000-0000-0000-000000000002"), uuid.New(), "Chicken Soup", sameCreatedAt, nil, nil)
	fuzzy := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("11111111-1111-1111-1111-111111111111"), uuid.New(), "Chikcen Soup", olderCreatedAt, nil, nil)
	seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("22222222-2222-2222-2222-222222222222"), uuid.New(), "Beef Chili", olderCreatedAt, nil, nil)

	return ctx, s, exactHigh, exactLow, fuzzy
}
