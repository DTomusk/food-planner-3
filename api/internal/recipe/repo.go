package recipe

import (
	"context"
	"food-planner/internal/db"
)

type Repo struct{}

func NewRepo() *Repo {
	return &Repo{}
}

func (r *Repo) CreateRecipe(recipe Recipe, ctx context.Context, db db.DBTX) error {
	_, err := db.ExecContext(ctx, "INSERT INTO recipes (id, name) VALUES ($1, $2)", recipe.ID, recipe.Name)
	return err
}

func (r *Repo) GetRecipeByID(id string, ctx context.Context, db db.DBTX) (Recipe, error) {
	var recipe Recipe
	row := db.QueryRowContext(ctx, "SELECT id, name FROM recipes WHERE id = $1", id)
	err := row.Scan(&recipe.ID, &recipe.Name)
	return recipe, err
}
