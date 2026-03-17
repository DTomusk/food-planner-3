package resolver

import (
	"context"
	"database/sql"
	"foodplanner/internal/auth"
	"foodplanner/internal/gql/graph/model"
	"foodplanner/internal/ingredient"
	"foodplanner/internal/recipe"
	"foodplanner/internal/testutil"
	"foodplanner/internal/testutil/seeds"
	"foodplanner/internal/user"
	"testing"

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

func TestRecipeResolver_CreateAndGetRecipe(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()

		testIngredient, err := seeds.SeedTestIngredient(ctx, tx)
		require.NoError(t, err, "Failed to seed test ingredient")

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)
		service := recipe.NewService(txRunner, recipe.NewRecipeRepo(), recipe.NewRecipeVersionRepo(), ingredientService, recipe.NewIngredientUsageRepo(), nil)
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
		ctx := context.Background()

		testIngredient, err := seeds.SeedTestIngredient(ctx, tx)
		require.NoError(t, err, "Failed to seed test ingredient")

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)
		service := recipe.NewService(txRunner, recipe.NewRecipeRepo(), recipe.NewRecipeVersionRepo(), ingredientService, recipe.NewIngredientUsageRepo(), nil)
		r := &Resolver{
			IngredientsService: ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
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
		ctx := context.Background()

		testIngredient, err := seeds.SeedTestIngredient(ctx, tx)
		require.NoError(t, err, "Failed to seed test ingredient")

		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)
		service := recipe.NewService(txRunner, recipe.NewRecipeRepo(), recipe.NewRecipeVersionRepo(), ingredientService, recipe.NewIngredientUsageRepo(), nil)
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
		ctx := context.Background()

		testIngredient, err := seeds.SeedTestIngredient(ctx, tx)
		require.NoError(t, err, "Failed to seed test ingredient")

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)
		service := recipe.NewService(txRunner, recipe.NewRecipeRepo(), recipe.NewRecipeVersionRepo(), ingredientService, recipe.NewIngredientUsageRepo(), nil)
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
		ctx := context.Background()

		testIngredient, err := seeds.SeedTestIngredient(ctx, tx)
		require.NoError(t, err, "Failed to seed test ingredient")

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)
		service := recipe.NewService(txRunner, recipe.NewRecipeRepo(), recipe.NewRecipeVersionRepo(), ingredientService, recipe.NewIngredientUsageRepo(), nil)
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

func TestRecipeResolver_RecipesAndRecipe_NoData(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		r := &Resolver{
			RecipeService: recipe.NewService(txRunner, recipe.NewRecipeRepo(), recipe.NewRecipeVersionRepo(), ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100), recipe.NewIngredientUsageRepo(), nil),
		}
		queryResolver := &queryResolver{r}

		ctx := context.Background()

		// Act & Assert
		recipes, err := queryResolver.Recipes(ctx)
		require.NoError(t, err)
		require.Len(t, recipes, 0)

		fetchedRecipe, err := queryResolver.Recipe(ctx, uuid.New().String())
		require.NoError(t, err)
		require.Nil(t, fetchedRecipe)
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

		err = seeds.InsertRecipeContainer(ctx, tx, recipeID, testUser.ID)
		require.NoError(t, err)

		err = seeds.InsertRecipeVersion(ctx, tx, versionID, recipeID, "Recipe Without Source", 15, 25, 2, 1)
		require.NoError(t, err)

		err = seeds.SetRecipeContainerCurrentVersion(ctx, tx, recipeID, versionID)
		require.NoError(t, err)

		r := &Resolver{
			RecipeService:      recipe.NewService(txRunner, recipe.NewRecipeRepo(), recipe.NewRecipeVersionRepo(), ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100), recipe.NewIngredientUsageRepo(), nil),
			IngredientsService: ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
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
