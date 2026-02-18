package recipe

import (
	"context"
	"database/sql"
	"foodplanner/internal/db"
)

type Repo struct{}

func NewRepo() *Repo {
	return &Repo{}
}

// Creates a new recipe including ingredient usages
// Returns recipe fields, ingredient usages can be loaded lazily
// TODO: use transaction
func (r *Repo) CreateRecipe(ctx context.Context, tx *sql.Tx, recipe *Recipe) (*Recipe, error) {
	var dbRecipe Recipe
	query := `INSERT INTO recipes (id, name) VALUES ($1, $2) RETURNING id, name`
	err := tx.QueryRowContext(ctx, query, recipe.ID, recipe.Name).Scan(&dbRecipe.ID, &dbRecipe.Name)
	if err != nil {
		return nil, err
	}
	usageQuery := `INSERT INTO ingredient_usages (id, recipe_id, ingredient_id, quantity, unit) VALUES ($1, $2, $3, $4, $5)`
	for _, usage := range recipe.Ingredients {
		_, err := tx.ExecContext(ctx, usageQuery, usage.ID, recipe.ID, usage.IngredientID, usage.Quantity, usage.Unit)
		if err != nil {
			return nil, err
		}
	}
	return &dbRecipe, nil
}

func (r *Repo) GetRecipeByID(ctx context.Context, db db.DBTX, id string) (*Recipe, error) {
	var recipe Recipe
	row := db.QueryRowContext(ctx, "SELECT id, name FROM recipes WHERE id = $1", id)
	err := row.Scan(&recipe.ID, &recipe.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &recipe, err
}

func (r *Repo) GetAllRecipes(ctx context.Context, db db.DBTX) ([]*Recipe, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, name FROM recipes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []*Recipe
	for rows.Next() {
		var recipe Recipe
		if err := rows.Scan(&recipe.ID, &recipe.Name); err != nil {
			return nil, err
		}
		recipes = append(recipes, &recipe)
	}
	return recipes, nil
}

func (r *Repo) IngredientExists(ctx context.Context, db db.DBTX, ingredientID string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM ingredients WHERE id = $1)", ingredientID).Scan(&exists)
	return exists, err
}

func (r *Repo) GetIngredientUsagesForRecipe(ctx context.Context, db db.DBTX, recipeID string) ([]*IngredientUsage, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, ingredient_id, quantity, unit FROM ingredient_usages WHERE recipe_id = $1", recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usages []*IngredientUsage
	for rows.Next() {
		var usage IngredientUsage
		if err := rows.Scan(&usage.ID, &usage.IngredientID, &usage.Quantity, &usage.Unit); err != nil {
			return nil, err
		}
		usages = append(usages, &usage)
	}
	return usages, nil
}
