package resolver

import (
	"context"
	"database/sql"
	"foodplanner/internal/auth"
	"foodplanner/internal/events"
	"foodplanner/internal/gql/graph/model"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/recipe"
	"foodplanner/internal/testutil"
	"foodplanner/internal/testutil/seeds"
	"foodplanner/internal/upload"
	"foodplanner/internal/user"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func ptrString(s string) *string {
	return &s
}

func ptrInt32(i int32) *int32 {
	return &i
}

func newTestRecipeService(t *testing.T, tx *sql.Tx, txRunner *testutil.TestTxRunner, ingredientService *ingredient.IngredientService) *recipe.Service {
	t.Helper()

	effectiveIngredientService := ingredientService
	if effectiveIngredientService == nil {
		effectiveIngredientService = ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)
	}

	repo, err := recipe.NewRecipeRepo(0.15, 0.85)
	require.NoError(t, err)

	eventBus := events.NewInMemoryEventBus(1, 32, txRunner)
	t.Cleanup(func() {
		_ = eventBus.Close(context.Background())
	})

	return recipe.NewService(
		txRunner,
		repo,
		recipe.NewRecipeVersionRepo(),
		effectiveIngredientService,
		recipe.NewIngredientUsageRepo(),
		nil,
		upload.NewUploadServiceWithProvider(tx, upload.NewStaticUploadProvider("https://upload.example.com", "https://cdn.example.com"), 0, upload.NewUploadRepo()),
		eventBus,
	)
}

func TestRecipeResolver_CreateAndGetRecipe(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := authContext(nil)

		testIngredient, err := seeds.SeedTestIngredient(ctx, tx)
		require.NoError(t, err, "Failed to seed test ingredient")

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)
		service := newTestRecipeService(t, tx, txRunner, ingredientService)
		r := &Resolver{
			RecipeService: service,
		}
		mutationResolver := &mutationResolver{r}

		input := model.CreateRecipeInput{
			Name: "Chocolate Cake",
			IngredientUsages: []*model.CreateIngredientUsageInput{
				{
					IngredientID: testIngredient.ID.String(),
					Quantity:     200,
					Unit:         1,
				},
			},
			PrepMins: 30,
			CookMins: 60,
			Portions: 2,
			RecipeSource: &model.CreateRecipeSourceInput{
				Type: 1,
				URL:  ptrString("https://example.com/chocolate-cake"),
			},
		}
		claims := auth.Claims{UserID: testUser.ID.String()}
		ctx = auth.ContextWithClaims(ctx, &claims)

		// Act
		recipeModel, err := mutationResolver.CreateRecipe(ctx, input)

		// Assert
		require.NoError(t, err, "CreateRecipe failed")
		require.Equal(t, "Chocolate Cake", recipeModel.CurrentVersion.Name)

		// Fetch the recipe directly from the DB to ensure it was saved correctly
		dbRecipe, err := service.GetRecipeByID(ctx, uuid.MustParse(recipeModel.ID))

		require.NoError(t, err)
		require.NotNil(t, dbRecipe, "Expected to find recipe in DB, got nil")
		require.Equal(t, "Chocolate Cake", dbRecipe.CurrentVersion.Name)
		require.Equal(t, 30, dbRecipe.CurrentVersion.PrepMins, "Recipe prep minutes mismatch")
		require.Equal(t, 60, dbRecipe.CurrentVersion.CookMins, "Recipe cook minutes mismatch")
		require.Equal(t, 2, dbRecipe.CurrentVersion.Portions, "Recipe portions mismatch")
	})
}

func TestRecipeResolver_CreateAndGetRecipe_WithResolver(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := authContext(nil)

		testIngredient, err := seeds.SeedTestIngredient(ctx, tx)
		require.NoError(t, err, "Failed to seed test ingredient")

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)
		service := newTestRecipeService(t, tx, txRunner, ingredientService)
		r := &Resolver{
			IngredientsService: ingredientService,
			RecipeService:      service,
			UserService:        user.NewUserService(tx, user.NewUserRepo()),
		}
		mutationResolver := &mutationResolver{r}

		input := model.CreateRecipeInput{
			Name: "Chocolate Cake",
			IngredientUsages: []*model.CreateIngredientUsageInput{
				{
					IngredientID: testIngredient.ID.String(),
					Quantity:     200,
					Unit:         1,
				},
			},
			PrepMins: 30,
			CookMins: 60,
			Portions: 2,
			RecipeSource: &model.CreateRecipeSourceInput{
				Type:      2,
				BookTitle: ptrString("Blah blah"),
				BookPage:  ptrInt32(42),
			},
		}
		claims := auth.Claims{UserID: testUser.ID.String()}
		ctx = auth.ContextWithClaims(ctx, &claims)

		// Act
		recipeModel, err := mutationResolver.CreateRecipe(ctx, input)

		// Assert
		require.NoError(t, err, "CreateRecipe failed")
		require.Equal(t, "Chocolate Cake", recipeModel.CurrentVersion.Name)

		queryResolver := &queryResolver{r}

		// Fetch the recipe using the resolver to ensure all nested resolvers work as expected
		fetchedRecipe, err := queryResolver.Recipe(ctx, recipeModel.ID)

		require.NoError(t, err, "Failed to fetch recipe with resolver")
		require.NotNil(t, fetchedRecipe, "Expected to fetch recipe with resolver, got nil")
		require.Equal(t, recipeModel.ID, fetchedRecipe.ID, "Recipe ID mismatch")
		require.Equal(t, recipeModel.CurrentVersion.ID, fetchedRecipe.CurrentVersion.ID, "Recipe version ID mismatch")
		require.Equal(t, "Chocolate Cake", fetchedRecipe.CurrentVersion.Name, "Recipe name mismatch")
		require.Equal(t, int32(30), fetchedRecipe.CurrentVersion.PrepMins, "Recipe prep minutes mismatch")
		require.Equal(t, int32(60), fetchedRecipe.CurrentVersion.CookMins, "Recipe cook minutes mismatch")
		require.Equal(t, int32(2), fetchedRecipe.CurrentVersion.Portions, "Recipe portions mismatch")

		recipeResolver := &recipeResolver{r}
		user, err := recipeResolver.Author(ctx, fetchedRecipe)
		require.NoError(t, err, "Failed to fetch user with resolver")
		require.NotNil(t, user, "Expected to fetch user with resolver, got nil")
		require.Equal(t, testUser.ID.String(), user.ID, "User ID mismatch")

		recipeVersionResolver := &recipeVersionResolver{r}
		ingredientUsages, err := recipeVersionResolver.IngredientUsages(ctx, fetchedRecipe.CurrentVersion)
		require.NoError(t, err, "Failed to fetch ingredient usages with resolver")
		require.NotNil(t, ingredientUsages, "Expected to fetch ingredient usages with resolver, got nil")
		require.Len(t, ingredientUsages, 1, "Expected exactly 1 ingredient usage")
		require.Equal(t, testIngredient.ID.String(), ingredientUsages[0].Ingredient.ID, "Ingredient ID mismatch")
		require.Equal(t, 200.0, ingredientUsages[0].Quantity, "Ingredient quantity mismatch")
		require.Equal(t, int32(1), ingredientUsages[0].Unit.Val, "Ingredient unit mismatch")

		recipeSource, err := recipeVersionResolver.Source(ctx, fetchedRecipe.CurrentVersion)
		require.NoError(t, err, "Failed to fetch recipe source with resolver")
		require.NotNil(t, recipeSource, "Expected to fetch recipe source with resolver, got nil")
		require.Equal(t, int32(2), recipeSource.Type, "Recipe source type mismatch")
		require.Equal(t, "Blah blah", *recipeSource.BookTitle, "Recipe source book title mismatch")
		require.Equal(t, int32(42), *recipeSource.BookPage, "Recipe source book page mismatch")
	})
}

func TestRecipeResolver_CreateRecipe_Unauthenticated(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := authContext(nil)

		testIngredient, err := seeds.SeedTestIngredient(ctx, tx)
		require.NoError(t, err, "Failed to seed test ingredient")

		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)
		service := newTestRecipeService(t, tx, txRunner, ingredientService)

		r := &Resolver{
			RecipeService: service,
		}
		mutationResolver := &mutationResolver{r}

		input := model.CreateRecipeInput{
			Name: "Chocolate Cake",
			IngredientUsages: []*model.CreateIngredientUsageInput{
				{
					IngredientID: testIngredient.ID.String(),
					Quantity:     200,
					Unit:         1,
				},
			},
			PrepMins: 30,
			CookMins: 60,
			Portions: 2,
			RecipeSource: &model.CreateRecipeSourceInput{
				Type: 1,
				URL:  ptrString("https://example.com/chocolate-cake"),
			},
		}

		// Act
		recipeModel, err := mutationResolver.CreateRecipe(ctx, input)

		// Assert
		require.Error(t, err, "Expected unauthenticated error")
		require.Nil(t, recipeModel)

		gqlErr, ok := err.(*gqlerror.Error)
		require.True(t, ok, "Expected GraphQL error type")
		require.Equal(t, "user is not authenticated", gqlErr.Message)
		require.Equal(t, "UNAUTHENTICATED", gqlErr.Extensions["code"])
	})
}

func TestRecipeResolver_UpdateRecipe_Unauthenticated(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := authContext(nil)

		testIngredient, err := seeds.SeedTestIngredient(ctx, tx)
		require.NoError(t, err, "Failed to seed test ingredient")

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)
		service := newTestRecipeService(t, tx, txRunner, ingredientService)
		r := &Resolver{
			RecipeService: service,
		}
		mutationResolver := &mutationResolver{r}

		createInput := model.CreateRecipeInput{
			Name: "Chocolate Cake",
			IngredientUsages: []*model.CreateIngredientUsageInput{
				{
					IngredientID: testIngredient.ID.String(),
					Quantity:     200,
					Unit:         1,
				},
			},
			PrepMins: 30,
			CookMins: 60,
			Portions: 2,
			RecipeSource: &model.CreateRecipeSourceInput{
				Type: 1,
				URL:  ptrString("https://example.com/chocolate-cake"),
			},
		}

		authenticatedCtx := auth.ContextWithClaims(ctx, &auth.Claims{UserID: testUser.ID.String()})
		createdRecipe, err := mutationResolver.CreateRecipe(authenticatedCtx, createInput)
		require.NoError(t, err, "CreateRecipe failed")

		updateInput := model.UpdateRecipeInput{
			ID: createdRecipe.ID,
			Details: &model.CreateRecipeInput{
				Name: "Chocolate Cake Updated",
				IngredientUsages: []*model.CreateIngredientUsageInput{
					{
						IngredientID: testIngredient.ID.String(),
						Quantity:     250,
						Unit:         1,
					},
				},
				PrepMins: 35,
				CookMins: 55,
				Portions: 3,
				RecipeSource: &model.CreateRecipeSourceInput{
					Type:         3,
					Instructions: ptrString("Grandma's notes"),
				},
			},
		}

		// Act
		updatedRecipe, err := mutationResolver.UpdateRecipe(ctx, updateInput)

		// Assert
		require.Error(t, err, "Expected unauthenticated error")
		require.Nil(t, updatedRecipe)

		gqlErr, ok := err.(*gqlerror.Error)
		require.True(t, ok, "Expected GraphQL error type")
		require.Equal(t, "user is not authenticated", gqlErr.Message)
		require.Equal(t, "UNAUTHENTICATED", gqlErr.Extensions["code"])
	})
}

func TestRecipeResolver_UpdateRecipe_WithVersionResolvers(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := authContext(nil)

		testIngredient, err := seeds.SeedTestIngredient(ctx, tx)
		require.NoError(t, err, "Failed to seed test ingredient")

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)
		service := newTestRecipeService(t, tx, txRunner, ingredientService)
		r := &Resolver{
			IngredientsService: ingredientService,
			RecipeService:      service,
		}
		mutationResolver := &mutationResolver{r}

		createInput := model.CreateRecipeInput{
			Name: "Chocolate Cake",
			IngredientUsages: []*model.CreateIngredientUsageInput{
				{
					IngredientID: testIngredient.ID.String(),
					Quantity:     200,
					Unit:         1,
				},
			},
			PrepMins: 30,
			CookMins: 60,
			Portions: 2,
			RecipeSource: &model.CreateRecipeSourceInput{
				Type: 1,
				URL:  ptrString("https://example.com/chocolate-cake"),
			},
			Publish: true,
		}

		authenticatedCtx := auth.ContextWithClaims(ctx, &auth.Claims{UserID: testUser.ID.String()})
		createdRecipe, err := mutationResolver.CreateRecipe(authenticatedCtx, createInput)
		require.NoError(t, err, "CreateRecipe failed")

		updateInput := model.UpdateRecipeInput{
			ID: createdRecipe.ID,
			Details: &model.CreateRecipeInput{
				Name: "Chocolate Cake Updated",
				IngredientUsages: []*model.CreateIngredientUsageInput{
					{
						IngredientID: testIngredient.ID.String(),
						Quantity:     250,
						Unit:         1,
					},
				},
				PrepMins: 35,
				CookMins: 55,
				Portions: 3,
				RecipeSource: &model.CreateRecipeSourceInput{
					Type:         3,
					Instructions: ptrString("Grandma's notes"),
				},
				Publish: true,
			},
		}

		// Act
		updatedRecipe, err := mutationResolver.UpdateRecipe(authenticatedCtx, updateInput)

		// Assert
		require.NoError(t, err, "UpdateRecipe failed")
		require.NotNil(t, updatedRecipe.CurrentVersion)
		require.Equal(t, "Chocolate Cake Updated", updatedRecipe.CurrentVersion.Name)
		require.Equal(t, int32(2), updatedRecipe.CurrentVersion.Version)

		// Additional assertions for resolvers
		queryResolver := &queryResolver{r}
		fetchedRecipe, err := queryResolver.Recipe(authenticatedCtx, createdRecipe.ID)
		require.NoError(t, err, "Failed to fetch updated recipe")
		require.NotNil(t, fetchedRecipe)
		require.Equal(t, int32(2), fetchedRecipe.CurrentVersion.Version)

		recipeResolver := &recipeResolver{r}
		versions, err := recipeResolver.Versions(authenticatedCtx, fetchedRecipe)
		require.NoError(t, err, "Failed to fetch recipe versions")
		require.Len(t, versions, 2, "Expected exactly 2 versions")

		versionOne, err := recipeResolver.Version(authenticatedCtx, fetchedRecipe, 1)
		require.NoError(t, err, "Failed to fetch version 1")
		require.NotNil(t, versionOne)
		require.Equal(t, "Chocolate Cake", versionOne.Name)

		versionTwo, err := recipeResolver.Version(authenticatedCtx, fetchedRecipe, 2)
		require.NoError(t, err, "Failed to fetch version 2")
		require.NotNil(t, versionTwo)
		require.Equal(t, "Chocolate Cake Updated", versionTwo.Name)

		resolvedCurrentVersion, err := recipeResolver.CurrentVersion(authenticatedCtx, &model.Recipe{CurrentVersionID: fetchedRecipe.CurrentVersionID})
		require.NoError(t, err, "Failed to resolve current version by ID")
		require.NotNil(t, resolvedCurrentVersion)
		require.Equal(t, fetchedRecipe.CurrentVersionID, resolvedCurrentVersion.ID)

		recipeVersionResolver := &recipeVersionResolver{r}
		resolvedSource, err := recipeVersionResolver.Source(authenticatedCtx, versionTwo)
		require.NoError(t, err, "Failed to resolve updated version source")
		require.NotNil(t, resolvedSource)
		require.Equal(t, int32(3), resolvedSource.Type)
		require.Equal(t, "Grandma's notes", *resolvedSource.Instructions)
	})
}

func TestRecipeResolver_DraftVisibility_RestrictedToOwner(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := authContext(nil)

		testIngredient, err := seeds.SeedTestIngredient(ctx, tx)
		require.NoError(t, err, "Failed to seed test ingredient")

		ownerUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed owner user")

		otherUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed non-owner user")

		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)
		service := newTestRecipeService(t, tx, txRunner, ingredientService)
		r := &Resolver{
			RecipeService: service,
		}

		mutationResolver := &mutationResolver{r}
		queryResolver := &queryResolver{r}
		recipeResolver := &recipeResolver{r}

		ownerCtx := auth.ContextWithClaims(ctx, &auth.Claims{UserID: ownerUser.ID.String()})
		nonOwnerCtx := auth.ContextWithClaims(ctx, &auth.Claims{UserID: otherUser.ID.String()})

		createdRecipe, err := mutationResolver.CreateRecipe(ownerCtx, model.CreateRecipeInput{
			Name: "Published Cake",
			IngredientUsages: []*model.CreateIngredientUsageInput{{
				IngredientID: testIngredient.ID.String(),
				Quantity:     200,
				Unit:         1,
			}},
			PrepMins: 30,
			CookMins: 45,
			Portions: 4,
			RecipeSource: &model.CreateRecipeSourceInput{
				Type: 1,
				URL:  ptrString("https://example.com/published-cake"),
			},
			Publish: true,
		})
		require.NoError(t, err)
		require.NotNil(t, createdRecipe)

		_, err = mutationResolver.UpdateRecipe(ownerCtx, model.UpdateRecipeInput{
			ID: createdRecipe.ID,
			Details: &model.CreateRecipeInput{
				Name: "Draft Cake",
				IngredientUsages: []*model.CreateIngredientUsageInput{{
					IngredientID: testIngredient.ID.String(),
					Quantity:     210,
					Unit:         1,
				}},
				PrepMins: 32,
				CookMins: 46,
				Portions: 4,
				RecipeSource: &model.CreateRecipeSourceInput{
					Type: 1,
					URL:  ptrString("https://example.com/draft-cake"),
				},
				Publish: false,
			},
		})
		require.NoError(t, err)

		fetchedRecipe, err := queryResolver.Recipe(ctx, createdRecipe.ID)
		require.NoError(t, err)
		require.NotNil(t, fetchedRecipe)

		nonOwnerDraft, err := recipeResolver.DraftVersion(nonOwnerCtx, fetchedRecipe)
		require.NoError(t, err)
		require.Nil(t, nonOwnerDraft, "Expected non-owner to not see draftVersion")

		ownerDraft, err := recipeResolver.DraftVersion(ownerCtx, fetchedRecipe)
		require.NoError(t, err)
		require.NotNil(t, ownerDraft, "Expected owner to see draftVersion")
		require.Equal(t, "Draft Cake", ownerDraft.Name)

		nonOwnerVersions, err := recipeResolver.Versions(nonOwnerCtx, fetchedRecipe)
		require.NoError(t, err)
		require.Len(t, nonOwnerVersions, 1, "Expected non-owner to only see published versions")
		require.Equal(t, int32(1), nonOwnerVersions[0].Version)

		ownerVersions, err := recipeResolver.Versions(ownerCtx, fetchedRecipe)
		require.NoError(t, err)
		require.Len(t, ownerVersions, 2, "Expected owner to see published and draft versions")

		nonOwnerDraftByVersion, err := recipeResolver.Version(nonOwnerCtx, fetchedRecipe, 2)
		require.NoError(t, err)
		require.Nil(t, nonOwnerDraftByVersion, "Expected non-owner to not resolve draft version by number")

		ownerDraftByVersion, err := recipeResolver.Version(ownerCtx, fetchedRecipe, 2)
		require.NoError(t, err)
		require.NotNil(t, ownerDraftByVersion, "Expected owner to resolve draft version by number")
		require.Equal(t, "Draft Cake", ownerDraftByVersion.Name)
	})
}

func TestRecipeResolver_CurrentVersion_UsesPreloadedValue(t *testing.T) {
	// Arrange
	r := &recipeResolver{&Resolver{}}

	preloadedVersion := &model.RecipeVersion{
		ID:   uuid.New().String(),
		Name: "Already Loaded",
	}

	// Act
	resolvedVersion, err := r.CurrentVersion(context.Background(), &model.Recipe{
		CurrentVersionID: uuid.New().String(),
		CurrentVersion:   preloadedVersion,
	})

	// Assert
	require.NoError(t, err)
	require.Same(t, preloadedVersion, resolvedVersion)
}

func TestRecipeResolver_Recipes_EmptyConnection(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		txRunner := testutil.NewTestTxRunner(tx)

		r := &Resolver{
			RecipeService: newTestRecipeService(t, tx, txRunner, nil),
		}
		ctx := context.Background()

		connection, err := (&queryResolver{r}).Recipes(ctx, nil, nil)

		require.NoError(t, err)
		require.NotNil(t, connection)
		require.Empty(t, connection.Edges)
		require.NotNil(t, connection.PageInfo)
		require.False(t, connection.PageInfo.HasNextPage)
		require.Nil(t, connection.PageInfo.EndCursor)
	})
}

func TestRecipeResolver_Recipes_PaginatesEdges(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		txRunner := testutil.NewTestTxRunner(tx)

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		// Seed 3 recipes with descending created_at: Alpha newest, Gamma oldest.
		now := time.Date(2026, time.March, 17, 12, 0, 0, 0, time.UTC)
		names := []string{"Recipe Alpha", "Recipe Beta", "Recipe Gamma"}
		for i, name := range names {
			recipeID := uuid.New()
			versionID := uuid.New()
			createdAt := now.Add(-time.Duration(i) * time.Minute)
			require.NoError(t, seeds.InsertPublishedRecipeContainer(ctx, tx, recipeID, testUser.ID))
			require.NoError(t, seeds.InsertPublishedRecipeVersion(ctx, tx, versionID, recipeID, name, 10, 20, 2, 1))
			require.NoError(t, seeds.SetRecipeContainerCurrentVersion(ctx, tx, recipeID, versionID))
			_, err = tx.ExecContext(ctx, `UPDATE recipe_containers SET created_at = $1 WHERE id = $2`, createdAt, recipeID)
			require.NoError(t, err)
			_, err = tx.ExecContext(ctx, `UPDATE recipe_versions SET created_at = $1 WHERE id = $2`, createdAt, versionID)
			require.NoError(t, err)
		}

		r := &Resolver{
			RecipeService: newTestRecipeService(t, tx, txRunner, nil),
		}
		qr := &queryResolver{r}

		// First page: 2 of 3 items.
		firstPage, err := qr.Recipes(ctx, &model.PaginationInput{First: ptrInt32(2)}, nil)
		require.NoError(t, err)
		require.NotNil(t, firstPage)
		require.Len(t, firstPage.Edges, 2)
		require.True(t, firstPage.PageInfo.HasNextPage)
		require.NotNil(t, firstPage.PageInfo.EndCursor)
		require.Equal(t, "Recipe Alpha", firstPage.Edges[0].Node.CurrentVersion.Name)
		require.Equal(t, "Recipe Beta", firstPage.Edges[1].Node.CurrentVersion.Name)
		require.NotEmpty(t, firstPage.Edges[0].Cursor)
		require.NotEmpty(t, firstPage.Edges[1].Cursor)
		// EndCursor is the cursor of the last edge on the page.
		require.Equal(t, *firstPage.PageInfo.EndCursor, firstPage.Edges[1].Cursor)

		// Second page: remaining item, no next page.
		secondPage, err := qr.Recipes(ctx, &model.PaginationInput{
			First: ptrInt32(2),
			After: firstPage.PageInfo.EndCursor,
		}, nil)
		require.NoError(t, err)
		require.NotNil(t, secondPage)
		require.Len(t, secondPage.Edges, 1)
		require.Equal(t, "Recipe Gamma", secondPage.Edges[0].Node.CurrentVersion.Name)
		require.False(t, secondPage.PageInfo.HasNextPage)
		require.Nil(t, secondPage.PageInfo.EndCursor)
	})
}

func TestRecipeVersionResolver_NoDataPaths(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		recipeID := uuid.New()
		versionID := uuid.New()

		err = seeds.InsertPublishedRecipeContainer(ctx, tx, recipeID, testUser.ID)
		require.NoError(t, err)

		err = seeds.InsertPublishedRecipeVersion(ctx, tx, versionID, recipeID, "Recipe Without Source", 15, 25, 2, 1)
		require.NoError(t, err)

		err = seeds.SetRecipeContainerCurrentVersion(ctx, tx, recipeID, versionID)
		require.NoError(t, err)

		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)

		r := &Resolver{
			RecipeService:      newTestRecipeService(t, tx, txRunner, ingredientService),
			IngredientsService: ingredientService,
		}
		recipeVersionResolver := &recipeVersionResolver{r}

		recipeVersionModel := &model.RecipeVersion{ID: versionID.String(), RecipeID: recipeID.String()}

		// Act & Assert
		resolvedRecipe, err := recipeVersionResolver.Recipe(ctx, recipeVersionModel)
		require.NoError(t, err, "Failed to resolve recipe from recipe version")
		require.NotNil(t, resolvedRecipe)
		require.Equal(t, recipeID.String(), resolvedRecipe.ID)

		resolvedSource, err := recipeVersionResolver.Source(ctx, recipeVersionModel)
		require.NoError(t, err, "Failed to resolve source from recipe version")
		require.NotNil(t, resolvedSource)
		require.Equal(t, int32(0), resolvedSource.Type)

		resolvedIngredientUsages, err := recipeVersionResolver.IngredientUsages(ctx, recipeVersionModel)
		require.NoError(t, err, "Failed to resolve ingredient usages from recipe version")
		require.Nil(t, resolvedIngredientUsages)

		resolvedRecipeWithNoRecipeID, err := recipeVersionResolver.Recipe(ctx, &model.RecipeVersion{ID: versionID.String(), RecipeID: ""})
		require.NoError(t, err)
		require.Nil(t, resolvedRecipeWithNoRecipeID)
	})
}
