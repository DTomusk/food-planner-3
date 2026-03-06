package resolver

import (
	_ "github.com/lib/pq"
)

func ptrString(s string) *string {
	return &s
}

func ptrInt32(i int32) *int32 {
	return &i
}

// func TestRecipeResolver_CreateAndGetRecipe(t *testing.T) {
// 	testutil.WithTx(t, func(tx *sql.Tx) {
// 		repo := recipe.NewRepo()
// 		txRunner := testutil.NewTestTxRunner(tx)
// 		ctx := context.Background()

// 		testIngredient, err := seeds.SeedTestIngredient(ctx, tx)
// 		require.NoError(t, err, "Failed to seed test ingredient")

// 		testUser, err := seeds.SeedTestUser(ctx, tx)
// 		require.NoError(t, err, "Failed to seed test user")

// 		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)
// 		service := recipe.NewService(txRunner, repo, ingredientService, nil)
// 		r := &Resolver{
// 			RecipeService: service,
// 		}
// 		mutationResolver := &mutationResolver{r}

// 		input := model.CreateRecipeInput{
// 			Name: "Chocolate Cake",
// 			IngredientUsages: []*model.CreateIngredientUsageInput{
// 				{
// 					IngredientID: testIngredient.ID.String(),
// 					Quantity:     200,
// 					Unit:         1,
// 				},
// 			},
// 			PrepMins: 30,
// 			CookMins: 60,
// 			Portions: 2,
// 			RecipeSource: &model.CreateRecipeSourceInput{
// 				Type: 1,
// 				URL:  ptrString("https://example.com/chocolate-cake"),
// 			},
// 		}
// 		claims := auth.Claims{UserID: testUser.ID.String()}
// 		ctx = auth.ContextWithClaims(ctx, &claims)
// 		recipeModel, err := mutationResolver.CreateRecipe(ctx, input)

// 		require.NoError(t, err, "CreateRecipe failed")
// 		require.Equal(t, "Chocolate Cake", recipeModel.Name)

// 		dbRecipe, err := repo.GetRecipeByID(ctx, tx, recipeModel.ID)

// 		require.NoError(t, err)
// 		require.NotNil(t, dbRecipe, "Expected to find recipe in DB, got nil")
// 		require.Equal(t, "Chocolate Cake", dbRecipe.Name)
// 		require.Equal(t, 30, dbRecipe.PrepMins, "Recipe prep minutes mismatch")
// 		require.Equal(t, 60, dbRecipe.CookMins, "Recipe cook minutes mismatch")
// 		require.Equal(t, 2, dbRecipe.Portions, "Recipe portions mismatch")
// 	})
// }

// func TestRecipeResolver_CreateAndGetRecipe_WithResolver(t *testing.T) {
// 	testutil.WithTx(t, func(tx *sql.Tx) {
// 		repo := recipe.NewRepo()
// 		txRunner := testutil.NewTestTxRunner(tx)
// 		ctx := context.Background()

// 		testIngredient, err := seeds.SeedTestIngredient(ctx, tx)
// 		require.NoError(t, err, "Failed to seed test ingredient")

// 		testUser, err := seeds.SeedTestUser(ctx, tx)
// 		require.NoError(t, err, "Failed to seed test user")

// 		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)
// 		service := recipe.NewService(txRunner, repo, ingredientService, nil)
// 		r := &Resolver{
// 			IngredientsService: ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100),
// 			RecipeService:      service,
// 			UserService:        user.NewUserService(tx, user.NewUserRepo()),
// 		}
// 		mutationResolver := &mutationResolver{r}

// 		input := model.CreateRecipeInput{
// 			Name: "Chocolate Cake",
// 			IngredientUsages: []*model.CreateIngredientUsageInput{
// 				{
// 					IngredientID: testIngredient.ID.String(),
// 					Quantity:     200,
// 					Unit:         1,
// 				},
// 			},
// 			PrepMins: 30,
// 			CookMins: 60,
// 			Portions: 2,
// 			RecipeSource: &model.CreateRecipeSourceInput{
// 				Type:      2,
// 				BookTitle: ptrString("Blah blah"),
// 				BookPage:  ptrInt32(42),
// 			},
// 		}
// 		claims := auth.Claims{UserID: testUser.ID.String()}
// 		ctx = auth.ContextWithClaims(ctx, &claims)
// 		recipeModel, err := mutationResolver.CreateRecipe(ctx, input)

// 		require.NoError(t, err, "CreateRecipe failed")
// 		require.Equal(t, "Chocolate Cake", recipeModel.Name)

// 		queryResolver := &queryResolver{r}
// 		fetchedRecipe, err := queryResolver.Recipe(ctx, recipeModel.ID)

// 		require.NoError(t, err, "Failed to fetch recipe with resolver")
// 		require.NotNil(t, fetchedRecipe, "Expected to fetch recipe with resolver, got nil")
// 		require.Equal(t, "Chocolate Cake", fetchedRecipe.Name)
// 		require.Equal(t, recipeModel.ID, fetchedRecipe.ID, "Recipe ID mismatch")
// 		require.Equal(t, "Chocolate Cake", fetchedRecipe.Name, "Recipe name mismatch")
// 		require.Equal(t, int32(30), fetchedRecipe.PrepMins, "Recipe prep minutes mismatch")
// 		require.Equal(t, int32(60), fetchedRecipe.CookMins, "Recipe cook minutes mismatch")
// 		require.Equal(t, int32(2), fetchedRecipe.Portions, "Recipe portions mismatch")

// 		recipeResolver := &recipeResolver{r}
// 		ingredientUsages, err := recipeResolver.IngredientUsages(ctx, fetchedRecipe)
// 		require.NoError(t, err, "Failed to fetch ingredient usages with resolver")
// 		require.NotNil(t, ingredientUsages, "Expected to fetch ingredient usages with resolver, got nil")
// 		require.Len(t, ingredientUsages, 1, "Expected exactly 1 ingredient usage")
// 		require.Equal(t, testIngredient.ID.String(), ingredientUsages[0].Ingredient.ID, "Ingredient ID mismatch")
// 		require.Equal(t, 200.0, ingredientUsages[0].Quantity, "Ingredient quantity mismatch")
// 		require.Equal(t, int32(1), ingredientUsages[0].Unit.Val, "Ingredient unit mismatch")

// 		user, err := recipeResolver.User(ctx, fetchedRecipe)
// 		require.NoError(t, err, "Failed to fetch user with resolver")
// 		require.NotNil(t, user, "Expected to fetch user with resolver, got nil")
// 		require.Equal(t, testUser.ID.String(), user.ID, "User ID mismatch")

// 		recipeSource, err := recipeResolver.Source(ctx, fetchedRecipe)
// 		require.NoError(t, err, "Failed to fetch recipe source with resolver")
// 		require.NotNil(t, recipeSource, "Expected to fetch recipe source with resolver, got nil")
// 		require.Equal(t, int32(2), recipeSource.Type, "Recipe source type mismatch")
// 		require.Equal(t, "Blah blah", *recipeSource.BookTitle, "Recipe source book title mismatch")
// 		require.Equal(t, int32(42), *recipeSource.BookPage, "Recipe source book page mismatch")
// 	})
//}
