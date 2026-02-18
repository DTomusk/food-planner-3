package ingredient

import (
	"context"
	"database/sql"
	"fmt"
	"foodplanner/internal/db"
	"foodplanner/internal/unit"
	"strings"

	"github.com/google/uuid"
)

type IngredientRepo struct{}

func NewIngredientRepo() *IngredientRepo {
	return &IngredientRepo{}
}

type IngredientRow struct {
	ID            uuid.UUID
	Name          string
	PreferredUnit int
	FileKey       string
}

func (r *IngredientRepo) IngredientExists(ctx context.Context, db db.DBTX, ingredientID string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM reference.ingredients WHERE id = $1)", ingredientID).Scan(&exists)
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
		var ingredientRow IngredientRow
		if err := rows.Scan(&ingredientRow.ID, &ingredientRow.Name, &ingredientRow.PreferredUnit, &ingredientRow.FileKey); err != nil {
			return nil, err
		}

		unit := unit.Unit(ingredientRow.PreferredUnit)
		if !unit.IsValid() {
			return nil, fmt.Errorf("invalid preferred unit %d for ingredient %s", ingredientRow.PreferredUnit, ingredientRow.ID)
		}

		ingredients = append(ingredients, &Ingredient{
			ID:            ingredientRow.ID,
			Name:          ingredientRow.Name,
			PreferredUnit: unit,
			FileKey:       ingredientRow.FileKey,
		})
	}
	return ingredients, nil
}

func (r *IngredientRepo) UpsertIngredients(ctx context.Context, db db.DBTX, ingredients []*Ingredient) error {
	if len(ingredients) == 0 {
		return nil
	}

	var (
		query  = "INSERT INTO reference.ingredients (id, name, preferred_unit, file_key) VALUES"
		args   []any
		values []string
	)

	for i, ingredient := range ingredients {
		start := i*4 + 1
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d)", start, start+1, start+2, start+3))
		args = append(args, ingredient.ID, ingredient.Name, int(ingredient.PreferredUnit), ingredient.FileKey)
	}

	query += " " + strings.Join(values, ", ") + `
	ON CONFLICT (file_key)
	DO UPDATE SET
	name = EXCLUDED.name,
	preferred_unit = EXCLUDED.preferred_unit;
	`

	_, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
