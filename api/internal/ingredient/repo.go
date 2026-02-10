package ingredient

import (
	"context"
	"foodplanner/internal/db"
)

type IngredientRepo struct{}

func NewIngredientRepo() *IngredientRepo {
	return &IngredientRepo{}
}

func (r *IngredientRepo) IngredientExists(ctx context.Context, db db.DBTX, ingredientID string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM ingredients WHERE id = $1)", ingredientID).Scan(&exists)
	return exists, err
}
