package recipe

import (
	"context"
	"database/sql"
	"testing"

	"food-planner/internal/testutil"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func TestCreateAndGetRecipe(t *testing.T) {
	r := NewRepo()

	testutil.WithTx(t, func(tx *sql.Tx) {
		id := uuid.New()
		err := r.CreateRecipe(Recipe{
			ID:   id,
			Name: "Chocolate Cake",
		}, context.Background(), tx)
		if err != nil {
			t.Fatalf("Failed to create recipe: %v", err)
		}

		got, err := r.GetRecipeByID(id.String(), context.Background(), tx)
		if err != nil {
			t.Fatalf("Failed to get recipe: %v", err)
		}

		if got.Name != "Chocolate Cake" {
			t.Errorf("Expected name %q, got %q", "Chocolate Cake", got.Name)
		}
	})
}
