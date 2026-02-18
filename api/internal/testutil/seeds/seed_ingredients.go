package seeds

import (
	"context"
	"foodplanner/internal/db"
	"foodplanner/internal/ingredient"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func InsertIngredient(ctx context.Context, db db.DBTX, ingredient *ingredient.Ingredient) error {
	query := `INSERT INTO reference.ingredients 
	(id, name, preferred_unit, file_key) 
	VALUES
	($1, $2, $3, $4)`
	_, err := db.ExecContext(ctx, query, ingredient.ID, ingredient.Name, ingredient.PreferredUnit, ingredient.FileKey)
	return err
}

func SeedTestIngredient(ctx context.Context, db db.DBTX, t *testing.T) (*ingredient.Ingredient, error) {
	ingredientID := uuid.New()
	testIngredient := ingredient.Ingredient{
		ID:            ingredientID,
		FileKey:       "test_ingredient",
		Name:          "Test Ingredient",
		PreferredUnit: 1,
	}
	err := InsertIngredient(context.Background(), db, &testIngredient)
	require.NoError(t, err, "Failed to seed test ingredient")
	return &testIngredient, nil
}
