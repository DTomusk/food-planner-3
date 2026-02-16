package ingredient

import (
	"context"
	"database/sql"
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

func (r *IngredientRepo) GetAllIngredients(ctx context.Context, db db.DBTX) ([]*Ingredient, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, name, preferred_unit, file_key FROM reference.ingredients")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var ingredients []*Ingredient
	for rows.Next() {
		var ingredient Ingredient
		if err := rows.Scan(&ingredient.ID, &ingredient.Name, &ingredient.PreferredUnit, &ingredient.FileKey); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, &ingredient)
	}
	return ingredients, nil
}

func (r *IngredientRepo) UpsertIngredients(ctx context.Context, db db.DBTX, ingredients []*Ingredient) error {
	return nil
}
