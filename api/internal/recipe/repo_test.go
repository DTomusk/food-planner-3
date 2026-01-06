package recipe

import (
	"context"
	"database/sql"
	"testing"

	"food-planner/internal/testutil"

	_ "github.com/lib/pq"
)

func TestCreateAndGetRecipe(t *testing.T) {
	r := NewRepo()

	testutil.WithTx(t, func(tx *sql.Tx) {
		entity, err := NewRecipe("Chocolate Cake")
		if err != nil {
			t.Fatalf("Failed to create recipe entity: %v", err)
		}
		_, err = r.CreateRecipe(entity, context.Background(), tx)
		if err != nil {
			t.Fatalf("Failed to create recipe: %v", err)
		}

		got, err := r.GetRecipeByID(entity.ID.String(), context.Background(), tx)
		if err != nil {
			t.Fatalf("Failed to get recipe: %v", err)
		}

		if got.Name != "Chocolate Cake" {
			t.Errorf("Expected name %q, got %q", "Chocolate Cake", got.Name)
		}
	})
}
