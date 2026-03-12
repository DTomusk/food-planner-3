package recipe

// func TestGetIngredientUsagesForRecipeVersion(t *testing.T) {
// 	testutil.WithTx(t, func(tx *sql.Tx) {
// 		// Arrange
// 		ctx := context.Background()
// 		recipeRepo := NewRecipeRepo()
// 		r := NewIngredientUsageRepo()
// 		ingredientID := uuid.New()
// 		testIngredient := ingredient.Ingredient{
// 			ID:            ingredientID,
// 			FileKey:       "test_ingredient",
// 			Name:          "Test Ingredient",
// 			PreferredUnit: 1,
// 		}
// 		err := seeds.InsertIngredient(ctx, tx, &testIngredient)
// 		require.NoError(t, err)

// 		testUser, err := seeds.SeedTestUser(ctx, tx)
// 		require.NoError(t, err)

// 		ingredientUsage, err := NewIngredientUsage(CreateIngredientUsageRequest{
// 			IngredientID: ingredientID.String(),
// 			Quantity:     200,
// 			Unit:         1,
// 		})
// 		require.NoError(t, err)

// 		recipeContainer, err := NewRecipe(
// 			"Recipe With Ingredients",
// 			testUser.ID,
// 			[]*IngredientUsage{ingredientUsage},
// 			30,
// 			60,
// 			8,
// 			nil,
// 		)
// 		require.NoError(t, err)
// 		created, err := recipeRepo.createRecipe(ctx, tx, recipeContainer)
// 		require.NoError(t, err)

// 		// Act
// 		err = r.insertIngredientUsages(ctx, tx, recipeContainer.CurrentVersion.Ingredients, created.CurrentVersion.ID)
// 		require.NoError(t, err, "Failed to insert ingredient usages")
// 		usages, err := r.getIngredientUsagesForRecipeVersion(ctx, tx, created.CurrentVersion.ID.String())

// 		// Assert
// 		require.NoError(t, err)
// 		require.Len(t, usages, 1, "Expected to find 1 ingredient usage")
// 		require.Equal(t, ingredientID, usages[0].IngredientID)
// 		require.Equal(t, float64(200), usages[0].Quantity)
// 	})
// }

// func TestInsertIngredientUsages(t *testing.T) {
// 	testutil.WithTx(t, func(tx *sql.Tx) {
// 		// Arrange
// 		ctx := context.Background()
// 		recipeRepo := NewRecipeRepo()
// 		r := NewIngredientUsageRepo()
// 		ingredientID1 := uuid.New()
// 		ingredientID2 := uuid.New()

// 		ingredient1 := ingredient.Ingredient{
// 			ID:            ingredientID1,
// 			FileKey:       "ingredient_1",
// 			Name:          "Ingredient 1",
// 			PreferredUnit: 1,
// 		}
// 		ingredient2 := ingredient.Ingredient{
// 			ID:            ingredientID2,
// 			FileKey:       "ingredient_2",
// 			Name:          "Ingredient 2",
// 			PreferredUnit: 2,
// 		}
// 		err := seeds.InsertIngredient(ctx, tx, &ingredient1)
// 		require.NoError(t, err)
// 		err = seeds.InsertIngredient(ctx, tx, &ingredient2)
// 		require.NoError(t, err)

// 		testUser, err := seeds.SeedTestUser(ctx, tx)
// 		require.NoError(t, err)

// 		usage1, err := NewIngredientUsage(CreateIngredientUsageRequest{
// 			IngredientID: ingredientID1.String(),
// 			Quantity:     200,
// 			Unit:         1,
// 		})
// 		require.NoError(t, err)

// 		usage2, err := NewIngredientUsage(CreateIngredientUsageRequest{
// 			IngredientID: ingredientID2.String(),
// 			Quantity:     300,
// 			Unit:         2,
// 		})
// 		require.NoError(t, err)

// 		recipeContainer, err := NewRecipe(
// 			"Recipe With Multiple Ingredients",
// 			testUser.ID,
// 			[]*IngredientUsage{usage1, usage2},
// 			30,
// 			60,
// 			8,
// 			nil,
// 		)
// 		require.NoError(t, err)
// 		created, err := recipeRepo.createRecipe(ctx, tx, recipeContainer)
// 		require.NoError(t, err)

// 		// Act
// 		err = r.insertIngredientUsages(ctx, tx, recipeContainer.CurrentVersion.Ingredients, created.CurrentVersion.ID)
// 		require.NoError(t, err, "Failed to insert ingredient usages")
// 		retrievedUsages, err := r.getIngredientUsagesForRecipeVersion(ctx, tx, created.CurrentVersion.ID.String())

// 		// Assert
// 		require.NoError(t, err)
// 		require.Len(t, retrievedUsages, 2, "Expected to find 2 ingredient usages")

// 		// Verify both ingredients are present
// 		ingredientIDs := map[uuid.UUID]bool{}
// 		for _, u := range retrievedUsages {
// 			ingredientIDs[u.IngredientID] = true
// 		}
// 		require.True(t, ingredientIDs[ingredientID1])
// 		require.True(t, ingredientIDs[ingredientID2])
// 	})
// }
