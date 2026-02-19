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
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestRecipeResolver_CreateAndGetRecipe(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		repo := recipe.NewRepo()
		txRunner := testutil.NewTestTxRunner(tx)
		ctx := context.Background()

		testIngredient, err := seeds.SeedTestIngredient(ctx, tx)
		require.NoError(t, err, "Failed to seed test ingredient")

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err, "Failed to seed test user")

		ingredientService := ingredient.NewIngredientService(txRunner, ingredient.NewIngredientRepo(), 100)
		service := recipe.NewService(txRunner, repo, ingredientService)
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
		}
		claims := auth.Claims{UserID: testUser.ID.String()}
		ctx = auth.ContextWithClaims(ctx, &claims)
		recipeModel, err := mutationResolver.CreateRecipe(ctx, input)

		require.NoError(t, err, "CreateRecipe failed")
		require.Equal(t, "Chocolate Cake", recipeModel.Name)

		dbRecipe, err := repo.GetRecipeByID(ctx, tx, recipeModel.ID)

		require.NoError(t, err)
		require.NotNil(t, dbRecipe, "Expected to find recipe in DB, got nil")
		require.Equal(t, "Chocolate Cake", dbRecipe.Name)
	})
}
