package recipe

import (
	"context"
	"database/sql"
	"food-planner/internal/testutil"
	"testing"
)

func TestCreateRecipe(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		s := NewService(tx, NewRepo())
		recipe, err := s.CreateRecipe("Vanilla Ice Cream", context.Background())
		if err != nil {
			t.Fatalf("Failed to create recipe: %v", err)
		}
		if recipe.Name != "Vanilla Ice Cream" {
			t.Errorf("Expected name %q, got %q", "Vanilla Ice Cream", recipe.Name)
		}
	})
}
