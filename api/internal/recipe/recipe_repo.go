package recipe

import (
	"context"
	"database/sql"
	"fmt"
	"foodplanner/internal/db"

	"github.com/google/uuid"
)

type recipeRepo struct {
	trigramWeight  float64
	fullTextWeight float64
}

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

func NewRecipeRepo(trigramWeight, fullTextWeight float64) (*recipeRepo, error) {
	if trigramWeight < 0 || fullTextWeight < 0 {
		return nil, fmt.Errorf("weights must be non-negative")
	}
	if trigramWeight+fullTextWeight != 1 {
		return nil, fmt.Errorf("weights must sum to 1")
	}
	return &recipeRepo{
		trigramWeight:  trigramWeight,
		fullTextWeight: fullTextWeight,
	}, nil
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
	// Relevance score is a mix of full-text search ranking and trigram similarity, weighted towards full-text search. Sorting is done by relevance score first, then created_at and id to ensure a deterministic order.
	// this can be fine tuned in the future
	scoreExpression := `
		($1 * ts_rank_cd(to_tsvector('english', coalesce(rv.name, '')), websearch_to_tsquery('english', $3)) 
		+ $2 * similarity(lower(rv.name), lower($3))) AS relevance_score
	`

	baseQuery := `
WITH ranked AS (
    SELECT
		rc.id AS container_id,
        rc.user_id,
		rc.created_at AS container_created_at,
        rc.current_version_id,
        rv.id AS version_id,
		rv.recipe_id AS version_recipe_id,
        rv.name,
        rv.prep_mins,
        rv.cook_mins,
        rv.portions,
        rv.created_at AS version_created_at,
        rv.version,
    ` + scoreExpression + `
    FROM recipe_containers rc
    JOIN recipe_versions rv ON rc.current_version_id = rv.id
    WHERE rc.deleted_on IS NULL
      AND (
        to_tsvector('english', coalesce(rv.name, '')) @@ websearch_to_tsquery('english', $3)
        OR lower(rv.name) % lower($3)
      )
)
SELECT
	container_id,
    user_id,
	container_created_at,
    current_version_id,
    version_id,
	version_recipe_id,
    name,
    prep_mins,
    cook_mins,
    portions,
    version_created_at,
    version,
    relevance_score
FROM ranked
`

	var q string
	var args []any

	if cursor != nil {
		if cursor.RelevanceScore == nil {
			return nil, ErrInvalidCursor
		}

		q = baseQuery + `
		WHERE (relevance_score, container_created_at, container_id) < ($4, $5, $6)
		ORDER BY relevance_score DESC, container_created_at DESC, container_id DESC
		LIMIT $7`

		args = []any{r.fullTextWeight, r.trigramWeight, query, *cursor.RelevanceScore, cursor.CreatedAt, cursor.ID, limit}
	} else {
		q = baseQuery + `
	ORDER BY relevance_score DESC, container_created_at DESC, container_id DESC
	LIMIT $4`
		args = []any{r.fullTextWeight, r.trigramWeight, query, limit}
	}

	rows, err := db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*RecipeListRow, 0, limit)
	for rows.Next() {
		var rc RecipeContainer
		var rv RecipeVersion
		score := new(float64)
		err := rows.Scan(
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
			score,
		)
		if err != nil {
			return nil, err
		}
		rc.CurrentVersion = &rv
		results = append(results, &RecipeListRow{Recipe: &rc, RelevanceScore: score})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
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
