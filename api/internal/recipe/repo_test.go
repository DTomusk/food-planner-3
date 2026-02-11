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
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		r := NewRepo()
		ingredientUsage, err := NewIngredientUsage(CreateIngredientUsageRequest{
			IngredientID: "04061e4e-6d4c-41d1-abcf-8b214927e1ed",
			Quantity:     200,
			Unit:         1,
		})
		if err != nil {
			t.Fatalf("Failed to create ingredient usage entity: %v", err)
		}
		recipe, err := NewRecipe("Chocolate Cake", []*IngredientUsage{ingredientUsage})
		if err != nil {
			t.Fatalf("Failed to create recipe entity: %v", err)
		}

		// Act
		_, err = r.CreateRecipe(context.Background(), tx, recipe)
		if err != nil {
			t.Fatalf("Failed to create recipe: %v", err)
		}

		got, err := r.GetRecipeByID(context.Background(), tx, recipe.ID.String())
		if err != nil {
			t.Fatalf("Failed to get recipe: %v", err)
		}

		// Assert
		if got.Name != "Chocolate Cake" {
			t.Errorf("Expected name %q, got %q", "Chocolate Cake", got.Name)
		}
	})
}

func TestGetRecipe_DoesNotErrorWhenNotFound(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		r := NewRepo()

		// Act
		recipe, err := r.GetRecipeByID(context.Background(), tx, "04061e4e-6d4c-41d1-abcf-8b214927e1ed")

		// Assert
		require.NoError(t, err)
		require.Nil(t, recipe)
	})
}
