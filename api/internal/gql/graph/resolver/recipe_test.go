package resolver

import (
	"context"
	"database/sql"
	"food-planner/internal/gql/graph/model"
	"food-planner/internal/recipe"
	"food-planner/internal/testutil"
	"testing"

	_ "github.com/lib/pq"
)

func TestRecipeResolver_CreateRecipe(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		repo := recipe.NewRepo()
		service := recipe.NewService(repo)
		r := &Resolver{
			DB:            tx,
			RecipeRepo:    repo,
			RecipeService: service,
		}
		mutationResolver := &mutationResolver{r}

		input := model.NewRecipe{
			Name: "Chocolate Cake",
		}
		recipeModel, err := mutationResolver.CreateRecipe(context.Background(), input)

		if err != nil {
			t.Fatalf("CreateRecipe failed: %v", err)
		}

		if recipeModel.Name != "Chocolate Cake" {
			t.Errorf("Expected recipe name %q, got %q", "Chocolate Cake", recipeModel.Name)
		}
	})
}
