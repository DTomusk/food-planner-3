package recipe

import (
	"context"
	"database/sql"
	"foodplanner/internal/db"
)

type recipeVersionRepo struct{}

func NewRecipeVersionRepo() *recipeVersionRepo {
	return &recipeVersionRepo{}
}

const (
	insertRecipeVersionQuery = `INSERT INTO recipe_versions 
	(id, recipe_id, name, prep_mins, cook_mins, portions, version) 
	VALUES ($1, $2, $3, $4, $5, $6, $7) 
	RETURNING id, recipe_id, name, prep_mins, cook_mins, portions, created_at`
)

func (r *recipeVersionRepo) createRecipeVersion(ctx context.Context, tx *sql.Tx, version *RecipeVersion) (*RecipeVersion, error) {
	var dbRecipeVersion RecipeVersion
	err := tx.QueryRowContext(
		ctx,
		insertRecipeVersionQuery,
		version.ID,
		version.RecipeID,
		version.Name,
		version.PrepMins,
		version.CookMins,
		version.Portions,
		version.Version,
	).Scan(
		&dbRecipeVersion.ID,
		&dbRecipeVersion.RecipeID,
		&dbRecipeVersion.Name,
		&dbRecipeVersion.PrepMins,
		&dbRecipeVersion.CookMins,
		&dbRecipeVersion.Portions,
		&dbRecipeVersion.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &dbRecipeVersion, nil
}

func (r *recipeVersionRepo) getRecipeVersionByID(ctx context.Context, db db.DBTX, id string) (*RecipeVersion, error) {
	var recipeVersion RecipeVersion
	row := db.QueryRowContext(ctx, `SELECT id, recipe_id, name, prep_mins, cook_mins, portions, created_at, version FROM recipe_versions WHERE id = $1`, id)
	err := row.Scan(
		&recipeVersion.ID,
		&recipeVersion.RecipeID,
		&recipeVersion.Name,
		&recipeVersion.PrepMins,
		&recipeVersion.CookMins,
		&recipeVersion.Portions,
		&recipeVersion.CreatedAt,
		&recipeVersion.Version,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &recipeVersion, nil
}

func (r *recipeVersionRepo) getRecipeVersionsByRecipeID(ctx context.Context, db db.DBTX, recipeID string) ([]*RecipeVersion, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, recipe_id, name, prep_mins, cook_mins, portions, created_at, version FROM recipe_versions WHERE recipe_id = $1 ORDER BY created_at DESC`, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []*RecipeVersion
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for rows.Next() {
		var version RecipeVersion
		if err := rows.Scan(
			&version.ID,
			&version.RecipeID,
			&version.Name,
			&version.PrepMins,
			&version.CookMins,
			&version.Portions,
			&version.CreatedAt,
			&version.Version,
		); err != nil {
			return nil, err
		}
		versions = append(versions, &version)
	}
	return versions, nil
}

// TODO: could have a recipe source repo
func (r *recipeVersionRepo) getRecipeSourceByRecipeVersionID(ctx context.Context, db db.DBTX, recipeVersionID string) (*RecipeSource, error) {
	var source RecipeSource
	row := db.QueryRowContext(ctx, "SELECT type, url, book_title, book_page, instructions FROM recipe_sources WHERE recipe_version_id = $1", recipeVersionID)
	err := row.Scan(&source.Type, &source.URL, &source.BookTitle, &source.BookPage, &source.Instructions)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &source, nil
}
