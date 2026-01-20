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

func (r *Repo) CreateRecipe(recipe *Recipe, ctx context.Context, db db.DBTX) (*Recipe, error) {
	var dbRecipe Recipe
	query := `INSERT INTO recipes (id, name) VALUES ($1, $2) RETURNING id, name`
	err := db.QueryRowContext(ctx, query, recipe.ID, recipe.Name).Scan(&dbRecipe.ID, &dbRecipe.Name)
	if err != nil {
		return nil, err
	}
	return &dbRecipe, nil
}

func (r *Repo) GetRecipeByID(id string, ctx context.Context, db db.DBTX) (*Recipe, error) {
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
