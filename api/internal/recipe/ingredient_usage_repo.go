package recipe

import (
	"context"
	"database/sql"
	"fmt"
	"foodplanner/internal/db"
	"strings"

	"github.com/google/uuid"
)

type ingredientUsageRepo struct{}

func NewIngredientUsageRepo() *ingredientUsageRepo {
	return &ingredientUsageRepo{}
}

func (r *ingredientUsageRepo) insertIngredientUsages(ctx context.Context, tx *sql.Tx, ingredientUsages []*IngredientUsage, recipeVersionId uuid.UUID) error {
	if len(ingredientUsages) == 0 {
		return nil
	}
	var (
		query  = "INSERT INTO ingredient_usages (id, recipe_version_id, ingredient_id, quantity, unit) VALUES"
		args   []any
		values []string
	)

	for i, u := range ingredientUsages {
		start := i*5 + 1
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", start, start+1, start+2, start+3, start+4))
		args = append(args, u.ID, recipeVersionId, u.IngredientID, u.Quantity, int(u.Unit))
	}

	query += " " + strings.Join(values, ", ") + ";"

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *ingredientUsageRepo) getIngredientUsagesForRecipeVersion(ctx context.Context, db db.DBTX, recipeVersionID uuid.UUID) ([]*IngredientUsage, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, ingredient_id, quantity, unit FROM ingredient_usages WHERE recipe_version_id = $1", recipeVersionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usages []*IngredientUsage
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for rows.Next() {
		var usage IngredientUsage
		if err := rows.Scan(&usage.ID, &usage.IngredientID, &usage.Quantity, &usage.Unit); err != nil {
			return nil, err
		}
		usages = append(usages, &usage)
	}
	return usages, nil
}
