package recipe

import (
	"context"
	"database/sql"
	"foodplanner/internal/db"

	"github.com/google/uuid"
)

type recipeVersionRepo struct{}

func NewRecipeVersionRepo() *recipeVersionRepo {
	return &recipeVersionRepo{}
}

const (
	recipeVersionColumns       = `id, recipe_id, name, description, prep_mins, cook_mins, portions, created_at, version, img_src, animal_product_level`
	recipeVersionInsertColumns = `id, recipe_id, name, description, prep_mins, cook_mins, portions, version, img_src, animal_product_level`

	insertRecipeVersionQuery = `INSERT INTO recipe_versions (` + recipeVersionInsertColumns + `)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING ` + recipeVersionColumns
	selectRecipeVersionByIDQuery                 = `SELECT ` + recipeVersionColumns + ` FROM recipe_versions WHERE id = $1`
	selectRecipeVersionByRecipeIDAndVersionQuery = `SELECT ` + recipeVersionColumns + `
	FROM recipe_versions
	WHERE recipe_id = $1 AND version = $2`
	selectRecipeVersionsByRecipeIDQuery = `SELECT ` + recipeVersionColumns + `
	FROM recipe_versions
	WHERE recipe_id = $1
	ORDER BY created_at DESC`
)

func recipeVersionScanDest(version *RecipeVersion) []any {
	return []any{
		&version.ID,
		&version.RecipeID,
		&version.Name,
		&version.Description,
		&version.PrepMins,
		&version.CookMins,
		&version.Portions,
		&version.CreatedAt,
		&version.Version,
		&version.ImgSrc,
		&version.AnimalProductLevel,
	}
}

func (r *recipeVersionRepo) createRecipeVersion(ctx context.Context, tx *sql.Tx, version *RecipeVersion) (*RecipeVersion, error) {
	var dbRecipeVersion RecipeVersion
	err := tx.QueryRowContext(
		ctx,
		insertRecipeVersionQuery,
		version.ID,
		version.RecipeID,
		version.Name,
		version.Description,
		version.PrepMins,
		version.CookMins,
		version.Portions,
		version.Version,
		version.ImgSrc,
		version.AnimalProductLevel,
	).Scan(recipeVersionScanDest(&dbRecipeVersion)...)
	if err != nil {
		return nil, err
	}
	return &dbRecipeVersion, nil
}

func (r *recipeVersionRepo) insertRecipeSource(ctx context.Context, tx *sql.Tx, recipeVersionID uuid.UUID, source *RecipeSource) error {
	sourceQuery := `INSERT INTO recipe_sources (recipe_version_id, type, url, book_title, book_page, instructions) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := tx.ExecContext(ctx, sourceQuery, recipeVersionID, int(source.Type), source.URL, source.BookTitle, source.BookPage, source.Instructions)
	if err != nil {
		return err
	}
	return nil
}

func (r *recipeVersionRepo) getRecipeVersionByID(ctx context.Context, db db.DBTX, id uuid.UUID) (*RecipeVersion, error) {
	var recipeVersion RecipeVersion
	row := db.QueryRowContext(ctx, selectRecipeVersionByIDQuery, id)
	err := row.Scan(recipeVersionScanDest(&recipeVersion)...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &recipeVersion, nil
}

func (r *recipeVersionRepo) getRecipeVersionByRecipeIDAndVersion(ctx context.Context, db db.DBTX, id uuid.UUID, version int) (*RecipeVersion, error) {
	var recipeVersion RecipeVersion
	row := db.QueryRowContext(ctx, selectRecipeVersionByRecipeIDAndVersionQuery, id, version)
	err := row.Scan(recipeVersionScanDest(&recipeVersion)...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &recipeVersion, nil
}

func (r *recipeVersionRepo) getRecipeVersionsByRecipeID(ctx context.Context, db db.DBTX, recipeID uuid.UUID) ([]*RecipeVersion, error) {
	rows, err := db.QueryContext(ctx, selectRecipeVersionsByRecipeIDQuery, recipeID)
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
		if err := rows.Scan(recipeVersionScanDest(&version)...); err != nil {
			return nil, err
		}
		versions = append(versions, &version)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return versions, nil
}

// TODO: could have a recipe source repo
func (r *recipeVersionRepo) getRecipeSourceByRecipeVersionID(ctx context.Context, db db.DBTX, recipeVersionID uuid.UUID) (*RecipeSource, error) {
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
