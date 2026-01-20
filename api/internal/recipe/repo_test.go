package recipe

import (
	"context"
	"database/sql"
	"testing"

	"foodplanner/internal/testutil"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
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

func TestGetRecipe_DoesNotErrorWhenNotFound(t *testing.T) {
	r := NewRepo()

	testutil.WithTx(t, func(tx *sql.Tx) {
		_, err := r.GetRecipeByID("04061e4e-6d4c-41d1-abcf-8b214927e1ed", context.Background(), tx)
		require.NoError(t, err)
	})
}
