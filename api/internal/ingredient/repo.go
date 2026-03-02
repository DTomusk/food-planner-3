package ingredient

import (
	"context"
	"database/sql"
	"fmt"
	"foodplanner/internal/db"
	"foodplanner/internal/unit"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
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
	Counter       *string
	Plural        *string
	CounterPlural *string
}

func (r *IngredientRepo) IngredientExists(ctx context.Context, db db.DBTX, ingredientID string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM reference.ingredients WHERE id = $1)", ingredientID).Scan(&exists)
	return exists, err
}

func (r *IngredientRepo) GetAllIngredients(ctx context.Context, db db.DBTX) ([]*Ingredient, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, name, preferred_unit, file_key, counter, plural, counter_plural FROM reference.ingredients")
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
		if err := rows.Scan(&ingredientRow.ID, &ingredientRow.Name, &ingredientRow.PreferredUnit, &ingredientRow.FileKey, &ingredientRow.Counter, &ingredientRow.Plural, &ingredientRow.CounterPlural); err != nil {
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
			Counter:       ingredientRow.Counter,
			Plural:        ingredientRow.Plural,
			CounterPlural: ingredientRow.CounterPlural,
		})
	}
	return ingredients, nil
}

func (r *IngredientRepo) GetIngredientsByIDs(ctx context.Context, db db.DBTX, ingredientIDs []string) ([]*Ingredient, error) {
	if len(ingredientIDs) == 0 {
		return []*Ingredient{}, nil
	}
	query := "SELECT id, name, preferred_unit, file_key, counter, plural, counter_plural FROM reference.ingredients WHERE id = ANY($1)"
	rows, err := db.QueryContext(ctx, query, pq.Array(ingredientIDs))
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
		if err := rows.Scan(&ingredientRow.ID, &ingredientRow.Name, &ingredientRow.PreferredUnit, &ingredientRow.FileKey, &ingredientRow.Counter, &ingredientRow.Plural, &ingredientRow.CounterPlural); err != nil {
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
			Counter:       ingredientRow.Counter,
			Plural:        ingredientRow.Plural,
			CounterPlural: ingredientRow.CounterPlural,
		})
	}
	return ingredients, nil
}

func (r *IngredientRepo) UpsertIngredients(ctx context.Context, db db.DBTX, ingredients []*Ingredient) error {
	if len(ingredients) == 0 {
		return nil
	}

	var (
		query  = "INSERT INTO reference.ingredients (id, name, preferred_unit, file_key, counter, plural, counter_plural) VALUES"
		args   []any
		values []string
	)

	for i, ingredient := range ingredients {
		start := i*7 + 1
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", start, start+1, start+2, start+3, start+4, start+5, start+6))
		args = append(args, ingredient.ID, ingredient.Name, int(ingredient.PreferredUnit), ingredient.FileKey, ingredient.Counter, ingredient.Plural, ingredient.CounterPlural)
	}

	query += " " + strings.Join(values, ", ") + `
	ON CONFLICT (file_key)
	DO UPDATE SET
	name = EXCLUDED.name,
	preferred_unit = EXCLUDED.preferred_unit,
	counter = EXCLUDED.counter,
	plural = EXCLUDED.plural,
	counter_plural = EXCLUDED.counter_plural;
	`

	_, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
