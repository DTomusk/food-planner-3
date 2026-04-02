package recipe

import (
	"context"
	"database/sql"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/testutil"
	"foodplanner/internal/testutil/seeds"
	"foodplanner/internal/upload"
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

		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := NewService(
			txRunner,
			repo,
			NewRecipeVersionRepo(),
			ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
			NewIngredientUsageRepo(),
			nil,
			upload.NewUploadServiceWithProvider(tx, upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com"), 0, upload.NewUploadRepo()),
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

		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := NewService(
			txRunner,
			repo,
			NewRecipeVersionRepo(),
			ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
			NewIngredientUsageRepo(),
			nil,
			upload.NewUploadServiceWithProvider(tx, upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com"), 0, upload.NewUploadRepo()),
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
		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := NewService(
			txRunner,
			repo,
			NewRecipeVersionRepo(),
			ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
			NewIngredientUsageRepo(),
			nil,
			upload.NewUploadServiceWithProvider(tx, upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com"), 0, upload.NewUploadRepo()),
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
		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := NewService(
			txRunner,
			repo,
			NewRecipeVersionRepo(),
			ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
			NewIngredientUsageRepo(),
			nil,
			upload.NewUploadServiceWithProvider(tx, upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com"), 0, upload.NewUploadRepo()),
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
		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := NewService(
			txRunner,
			repo,
			NewRecipeVersionRepo(),
			ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
			NewIngredientUsageRepo(),
			nil,
			upload.NewUploadServiceWithProvider(tx, upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com"), 0, upload.NewUploadRepo()),
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

		repo, err := NewRecipeRepo(0.15, 0.85)

		s := NewService(
			txRunner,
			repo,
			NewRecipeVersionRepo(),
			ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
			NewIngredientUsageRepo(),
			nil,
			upload.NewUploadServiceWithProvider(tx, upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com"), 0, upload.NewUploadRepo()),
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
		repo, err := NewRecipeRepo(0.15, 0.85)
		require.NoError(t, err)
		s := NewService(
			txRunner,
			repo,
			NewRecipeVersionRepo(),
			ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
			NewIngredientUsageRepo(),
			nil,
			upload.NewUploadServiceWithProvider(tx, upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com"), 0, upload.NewUploadRepo()),
		)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		newestCreatedAt := time.Date(2026, time.March, 17, 11, 42, 48, 147630000, time.UTC)
		middleCreatedAt := newestCreatedAt.Add(-1 * time.Minute)
		oldestCreatedAt := newestCreatedAt.Add(-2 * time.Minute)

		newest := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1"), uuid.New(), "Newest", newestCreatedAt, nil)
		middle := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2"), uuid.New(), "Middle", middleCreatedAt, nil)
		oldest := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("cccccccc-cccc-cccc-cccc-ccccccccccc3"), uuid.New(), "Oldest", oldestCreatedAt, nil)

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
		s := NewService(
			txRunner,
			repo,
			NewRecipeVersionRepo(),
			ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
			NewIngredientUsageRepo(),
			nil,
			upload.NewUploadServiceWithProvider(tx, upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com"), 0, upload.NewUploadRepo()),
		)

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
		require.Equal(t, filterHashForParams(RecipeCursorModeNewest, nil, nil), edgeCursor.FilterHash)
		require.Nil(t, edgeCursor.RelevanceScore)

		pageCursor, err := ParseRecipeCursor(nextCursor)
		require.NoError(t, err)
		require.NotNil(t, pageCursor)
		require.Equal(t, RecipeCursorModeNewest, pageCursor.Mode)
		require.Equal(t, filterHashForParams(RecipeCursorModeNewest, nil, nil), pageCursor.FilterHash)
	})
}

func TestGetRecipes_StaleCursorModeIsIgnored(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx, s, newest, middle, _ := setupRecipeListFixture(t, tx)

		searchQuery := "pasta"
		score := 0.8
		staleCursor := (&RecipeCursor{
			Mode:           RecipeCursorModeRelevance,
			FilterHash:     filterHashForParams(RecipeCursorModeRelevance, &searchQuery, nil),
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
			FilterHash: filterHashForParams(RecipeCursorModeNewest, nil, nil),
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

		expectedHash := filterHashForParams(RecipeCursorModeRelevance, &query, nil)
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
			FilterHash: filterHashForParams(RecipeCursorModeNewest, nil, nil),
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
			FilterHash:     filterHashForParams(RecipeCursorModeRelevance, &otherQuery, nil),
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
		s := NewService(
			txRunner,
			repo,
			NewRecipeVersionRepo(),
			ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
			NewIngredientUsageRepo(),
			nil,
			upload.NewUploadServiceWithProvider(tx, upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com"), 0, upload.NewUploadRepo()),
		)

		userA, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)
		userB, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		base := time.Date(2026, time.March, 17, 11, 42, 48, 147630000, time.UTC)
		aNewest := seedRecipeForListTests(t, ctx, tx, userA.ID, uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1"), uuid.New(), "A Newest", base, nil)
		aOlder := seedRecipeForListTests(t, ctx, tx, userA.ID, uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2"), uuid.New(), "A Older", base.Add(-1*time.Minute), nil)
		_ = seedRecipeForListTests(t, ctx, tx, userB.ID, uuid.MustParse("cccccccc-cccc-cccc-cccc-ccccccccccc3"), uuid.New(), "B Recipe", base.Add(-2*time.Minute), nil)

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
		require.Equal(t, filterHashForParams(RecipeCursorModeNewest, nil, &userAID), parsed.FilterHash)
	})
}

func setupRecipeListFixture(t *testing.T, tx *sql.Tx) (context.Context, *Service, listedRecipeSeed, listedRecipeSeed, listedRecipeSeed) {
	t.Helper()

	ctx := context.Background()
	txRunner := testutil.NewTestTxRunner(tx)
	repo, err := NewRecipeRepo(0.15, 0.85)
	require.NoError(t, err)
	s := NewService(
		txRunner,
		repo,
		NewRecipeVersionRepo(),
		ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
		NewIngredientUsageRepo(),
		nil,
		upload.NewUploadServiceWithProvider(tx, upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com"), 0, upload.NewUploadRepo()),
	)

	testUser, err := seeds.SeedTestUser(ctx, tx)
	require.NoError(t, err, "Failed to seed test user")

	newestCreatedAt := time.Date(2026, time.March, 17, 11, 42, 48, 147630000, time.UTC)
	middleCreatedAt := newestCreatedAt.Add(-1 * time.Minute)
	oldestCreatedAt := newestCreatedAt.Add(-2 * time.Minute)

	newest := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1"), uuid.New(), "Newest", newestCreatedAt, nil)
	middle := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2"), uuid.New(), "Middle", middleCreatedAt, nil)
	oldest := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("cccccccc-cccc-cccc-cccc-ccccccccccc3"), uuid.New(), "Oldest", oldestCreatedAt, nil)

	return ctx, s, newest, middle, oldest
}

func setupRecipeSearchFixture(t *testing.T, tx *sql.Tx) (context.Context, *Service, listedRecipeSeed, listedRecipeSeed, listedRecipeSeed) {
	t.Helper()

	ctx := context.Background()
	txRunner := testutil.NewTestTxRunner(tx)
	repo, err := NewRecipeRepo(0.15, 0.85)
	require.NoError(t, err)
	s := NewService(
		txRunner,
		repo,
		NewRecipeVersionRepo(),
		ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
		NewIngredientUsageRepo(),
		nil,
		upload.NewUploadServiceWithProvider(tx, upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com"), 0, upload.NewUploadRepo()),
	)

	testUser, err := seeds.SeedTestUser(ctx, tx)
	require.NoError(t, err, "Failed to seed test user")

	sameCreatedAt := time.Date(2026, time.March, 17, 11, 40, 48, 147630000, time.UTC)
	olderCreatedAt := sameCreatedAt.Add(-1 * time.Minute)

	exactHigh := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("ffffffff-ffff-ffff-ffff-fffffffffff1"), uuid.New(), "Chicken Soup", sameCreatedAt, nil)
	exactLow := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("00000000-0000-0000-0000-000000000002"), uuid.New(), "Chicken Soup", sameCreatedAt, nil)
	fuzzy := seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("11111111-1111-1111-1111-111111111111"), uuid.New(), "Chikcen Soup", olderCreatedAt, nil)
	seedRecipeForListTests(t, ctx, tx, testUser.ID, uuid.MustParse("22222222-2222-2222-2222-222222222222"), uuid.New(), "Beef Chili", olderCreatedAt, nil)

	return ctx, s, exactHigh, exactLow, fuzzy
}
