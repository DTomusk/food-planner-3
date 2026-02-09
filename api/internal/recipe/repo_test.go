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
		_, err = r.CreateRecipe(context.Background(), tx, entity)
		if err != nil {
			t.Fatalf("Failed to create recipe: %v", err)
		}

		got, err := r.GetRecipeByID(context.Background(), tx, entity.ID.String())
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
		_, err := r.GetRecipeByID(context.Background(), tx, "04061e4e-6d4c-41d1-abcf-8b214927e1ed")
		require.NoError(t, err)
	})
}
