package seeds

import (
	"context"
	"foodplanner/internal/db"
	"foodplanner/internal/ingredient"
)

func InsertIngredient(ctx context.Context, db db.DBTX, ingredient *ingredient.Ingredient) error {
	query := `INSERT INTO reference.ingredients 
	(id, name, preferred_unit, file_key) 
	VALUES
	($1, $2, $3, $4)`
	_, err := db.ExecContext(ctx, query, ingredient.ID, ingredient.Name, ingredient.PreferredUnit, ingredient.FileKey)
	return err
}
