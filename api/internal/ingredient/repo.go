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

const (
	selectIngredientColumnsBaseQuery = "SELECT id, name, preferred_unit, file_key, counter, plural, counter_plural, animal_product_level, contains_gluten, taxonomy_parent_id, processing_level, is_searchable FROM reference.ingredients"
	selectIngredientsByIDsQuery      = selectIngredientColumnsBaseQuery + " WHERE id = ANY($1)"
)

func NewIngredientRepo() *IngredientRepo {
	return &IngredientRepo{}
}

type IngredientRow struct {
	ID                 uuid.UUID
	Name               string
	PreferredUnit      int
	FileKey            string
	Counter            *string
	Plural             *string
	CounterPlural      *string
	AnimalProductLevel int
	ContainsGluten     bool
	TaxonomyParentID   *uuid.UUID
	ProcessingLevel    int
	IsSearchable       bool
}

func (r *IngredientRepo) IngredientExists(ctx context.Context, db db.DBTX, ingredientID string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM reference.ingredients WHERE id = $1)", ingredientID).Scan(&exists)
	return exists, err
}

func (r *IngredientRepo) GetAllIngredients(ctx context.Context, db db.DBTX) ([]*Ingredient, error) {
	rows, err := db.QueryContext(ctx, selectIngredientColumnsBaseQuery)
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
		if err := rows.Scan(&ingredientRow.ID, &ingredientRow.Name, &ingredientRow.PreferredUnit, &ingredientRow.FileKey, &ingredientRow.Counter, &ingredientRow.Plural, &ingredientRow.CounterPlural, &ingredientRow.AnimalProductLevel, &ingredientRow.ContainsGluten, &ingredientRow.TaxonomyParentID, &ingredientRow.ProcessingLevel, &ingredientRow.IsSearchable); err != nil {
			return nil, err
		}

		// Skip non-searchable ingredients
		if !ingredientRow.IsSearchable {
			continue
		}

		resolvedUnit := unit.Unit(ingredientRow.PreferredUnit)
		if !isPreferredUnitAllowed(ingredientRow.PreferredUnit, ingredientRow.IsSearchable) {
			return nil, fmt.Errorf("invalid preferred unit %d for ingredient %s", ingredientRow.PreferredUnit, ingredientRow.ID)
		}

		ingredients = append(ingredients, &Ingredient{
			ID:                 ingredientRow.ID,
			Name:               ingredientRow.Name,
			PreferredUnit:      resolvedUnit,
			FileKey:            ingredientRow.FileKey,
			Counter:            ingredientRow.Counter,
			Plural:             ingredientRow.Plural,
			CounterPlural:      ingredientRow.CounterPlural,
			AnimalProductLevel: AnimalProductLevel(ingredientRow.AnimalProductLevel),
			ContainsGluten:     ingredientRow.ContainsGluten,
			TaxonomyParentID:   ingredientRow.TaxonomyParentID,
			ProcessedLevel:     ProcessedLevel(ingredientRow.ProcessingLevel),
			IsSearchable:       ingredientRow.IsSearchable,
		})
	}
	return ingredients, nil
}

func (r *IngredientRepo) GetIngredientsByIDs(ctx context.Context, db db.DBTX, ingredientIDs []string) ([]*Ingredient, error) {
	if len(ingredientIDs) == 0 {
		return []*Ingredient{}, nil
	}
	rows, err := db.QueryContext(ctx, selectIngredientsByIDsQuery, pq.Array(ingredientIDs))
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
		if err := rows.Scan(&ingredientRow.ID, &ingredientRow.Name, &ingredientRow.PreferredUnit, &ingredientRow.FileKey, &ingredientRow.Counter, &ingredientRow.Plural, &ingredientRow.CounterPlural, &ingredientRow.AnimalProductLevel, &ingredientRow.ContainsGluten, &ingredientRow.TaxonomyParentID, &ingredientRow.ProcessingLevel, &ingredientRow.IsSearchable); err != nil {
			return nil, err
		}

		resolvedUnit := unit.Unit(ingredientRow.PreferredUnit)
		if !isPreferredUnitAllowed(ingredientRow.PreferredUnit, ingredientRow.IsSearchable) {
			return nil, fmt.Errorf("invalid preferred unit %d for ingredient %s", ingredientRow.PreferredUnit, ingredientRow.ID)
		}

		ingredients = append(ingredients, &Ingredient{
			ID:                 ingredientRow.ID,
			Name:               ingredientRow.Name,
			PreferredUnit:      resolvedUnit,
			FileKey:            ingredientRow.FileKey,
			Counter:            ingredientRow.Counter,
			Plural:             ingredientRow.Plural,
			CounterPlural:      ingredientRow.CounterPlural,
			AnimalProductLevel: AnimalProductLevel(ingredientRow.AnimalProductLevel),
			ContainsGluten:     ingredientRow.ContainsGluten,
			TaxonomyParentID:   ingredientRow.TaxonomyParentID,
			ProcessedLevel:     ProcessedLevel(ingredientRow.ProcessingLevel),
			IsSearchable:       ingredientRow.IsSearchable,
		})
	}
	return ingredients, nil
}

func (r *IngredientRepo) UpsertIngredients(ctx context.Context, db db.DBTX, ingredients []*Ingredient) error {
	if len(ingredients) == 0 {
		return nil
	}

	var (
		query  = "INSERT INTO reference.ingredients (id, name, preferred_unit, file_key, counter, plural, counter_plural, animal_product_level, contains_gluten, taxonomy_parent_id, processing_level, is_searchable) VALUES"
		args   []any
		values []string
	)

	for i, ingredient := range ingredients {
		start := i*12 + 1
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", start, start+1, start+2, start+3, start+4, start+5, start+6, start+7, start+8, start+9, start+10, start+11))
		args = append(args, ingredient.ID, ingredient.Name, int(ingredient.PreferredUnit), ingredient.FileKey, ingredient.Counter, ingredient.Plural, ingredient.CounterPlural, int(ingredient.AnimalProductLevel), ingredient.ContainsGluten, ingredient.TaxonomyParentID, int(ingredient.ProcessedLevel), ingredient.IsSearchable)
	}

	query += " " + strings.Join(values, ", ") + `
	ON CONFLICT (file_key)
	DO UPDATE SET
	name = EXCLUDED.name,
	preferred_unit = EXCLUDED.preferred_unit,
	counter = EXCLUDED.counter,
	plural = EXCLUDED.plural,
	counter_plural = EXCLUDED.counter_plural,
	animal_product_level = EXCLUDED.animal_product_level,
	contains_gluten = EXCLUDED.contains_gluten,
	taxonomy_parent_id = EXCLUDED.taxonomy_parent_id,
	processing_level = EXCLUDED.processing_level,
	is_searchable = EXCLUDED.is_searchable;
	`

	_, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
