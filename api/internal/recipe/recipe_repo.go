package recipe

import (
	"context"
	"database/sql"
	"foodplanner/internal/db"

	"github.com/google/uuid"
)

type recipeRepo struct{}

const (
	selectRecipeContainerWithVersionBaseQuery = `SELECT 
	rc.id, rc.user_id, rc.created_at, rc.current_version_id,
	rv.id, rv.recipe_id, rv.name, rv.prep_mins, rv.cook_mins, rv.portions, rv.created_at, rv.version
	FROM recipe_containers rc
	JOIN recipe_versions rv ON rc.current_version_id = rv.id`
	selectRecipeContainerWithVersionByIDQuery = selectRecipeContainerWithVersionBaseQuery + `
	WHERE rc.id = $1`
	selectRecipeContainerWithVersionByUserIDQuery = selectRecipeContainerWithVersionBaseQuery + `
	WHERE rc.user_id = $1`
)

func NewRecipeRepo() *recipeRepo {
	return &recipeRepo{}
}

func (r *recipeRepo) createRecipeContainer(ctx context.Context, tx *sql.Tx, recipeContainer *RecipeContainer) (*RecipeContainer, error) {
	var dbRecipeContainer RecipeContainer
	containerQuery := `INSERT INTO recipe_containers 
	(id, user_id) 
	VALUES ($1, $2) 
	RETURNING id, user_id, created_at`
	err := tx.QueryRowContext(
		ctx,
		containerQuery,
		recipeContainer.ID,
		recipeContainer.UserID,
	).Scan(
		&dbRecipeContainer.ID,
		&dbRecipeContainer.UserID,
		&dbRecipeContainer.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &dbRecipeContainer, nil
}

func (r *recipeRepo) updateRecipeCurrentVersion(ctx context.Context, tx *sql.Tx, recipeID, versionID uuid.UUID) error {
	updateQuery := `UPDATE recipe_containers SET current_version_id = $1 WHERE id = $2`
	_, err := tx.ExecContext(ctx, updateQuery, versionID, recipeID)
	return err
}

func (r *recipeRepo) getRecipeByID(ctx context.Context, db db.DBTX, id uuid.UUID) (*RecipeContainer, error) {
	row := db.QueryRowContext(ctx,
		selectRecipeContainerWithVersionByIDQuery,
		id,
	)
	rc, err := scanRecipeContainerWithVersion(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return rc, nil
}

func (r *recipeRepo) getRecipesByCreatedAt(ctx context.Context, db db.DBTX, limit int, cursor *RecipeCursor) ([]*RecipeListRow, error) {
	const baseQuery = selectRecipeContainerWithVersionBaseQuery + `
	WHERE rc.deleted_on IS NULL`

	var query string
	var args []any

	if cursor != nil {
		query = baseQuery + `
	AND (rc.created_at, rc.id) < ($1, $2)
	ORDER BY rc.created_at DESC, rc.id DESC
	LIMIT $3`
		args = []any{cursor.CreatedAt, cursor.ID, limit}
	} else {
		query = baseQuery + `
	ORDER BY rc.created_at DESC, rc.id DESC
	LIMIT $1`
		args = []any{limit}
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []*RecipeListRow
	for rows.Next() {
		rc, err := scanRecipeContainerWithVersion(rows)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, &RecipeListRow{Recipe: rc})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return recipes, nil
}

func (r *recipeRepo) getRecipesByRelevance(ctx context.Context, db db.DBTX, query string, limit int, cursor *RecipeCursor) ([]*RecipeListRow, error) {
	return nil, nil
}

func (r *recipeRepo) getRecipesByUserID(ctx context.Context, db db.DBTX, userID uuid.UUID) ([]*RecipeContainer, error) {
	rows, err := db.QueryContext(ctx,
		selectRecipeContainerWithVersionByUserIDQuery,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []*RecipeContainer
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for rows.Next() {
		rc, err := scanRecipeContainerWithVersion(rows)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, rc)
	}
	return recipes, nil
}

func scanRecipeContainerWithVersion(
	scanner interface{ Scan(dest ...any) error },
) (*RecipeContainer, error) {

	var rc RecipeContainer
	var rv RecipeVersion

	err := scanner.Scan(
		&rc.ID,
		&rc.UserID,
		&rc.CreatedAt,
		&rc.CurrentVersionID,
		&rv.ID,
		&rv.RecipeID,
		&rv.Name,
		&rv.PrepMins,
		&rv.CookMins,
		&rv.Portions,
		&rv.CreatedAt,
		&rv.Version,
	)
	if err != nil {
		return nil, err
	}

	rc.CurrentVersion = &rv
	return &rc, nil
}
