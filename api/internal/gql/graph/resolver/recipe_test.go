package resolver

import (
	"context"
	"database/sql"
	"foodplanner/internal/auth"
	"foodplanner/internal/gql/graph/model"
	"foodplanner/internal/recipe"
	"foodplanner/internal/testutil"
	"testing"

	_ "github.com/lib/pq"
)

func TestRecipeResolver_CreateAndGetRecipe(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		repo := recipe.NewRepo()
		service := recipe.NewService(tx, repo)
		r := &Resolver{
			RecipeService: service,
		}
		mutationResolver := &mutationResolver{r}

		input := model.NewRecipe{
			Name: "Chocolate Cake",
		}
		ctx := context.Background()
		claims := auth.Claims{UserID: "some-user-id"}
		ctx = auth.ContextWithClaims(ctx, &claims)
		recipeModel, err := mutationResolver.CreateRecipe(ctx, input)

		if err != nil {
			t.Fatalf("CreateRecipe failed: %v", err)
		}

		if recipeModel.Name != "Chocolate Cake" {
			t.Errorf("Expected recipe name %q, got %q", "Chocolate Cake", recipeModel.Name)
		}

		if recipeModel.ID == "" {
			t.Errorf("Expected recipe ID to be set, got empty string")
		}

		dbRecipe, err := repo.GetRecipeByID(recipeModel.ID, context.Background(), tx)
		if err != nil {
			t.Fatalf("GetRecipeByID failed: %v", err)
		}
		if dbRecipe == nil {
			t.Fatalf("Expected to find recipe in DB, got nil")
		}
		if dbRecipe.Name != "Chocolate Cake" {
			t.Errorf("Expected DB recipe name %q, got %q", "Chocolate Cake", dbRecipe.Name)
		}
	})
}
