package recipe

import (
	"context"
	"database/sql"
	"fmt"
	"foodplanner/internal/db"
	"strings"
	"time"

	"github.com/google/uuid"
)

type recipeRepo struct {
	trigramWeight  float64
	fullTextWeight float64
}

const (
	selectRecipeContainerWithVersionBaseQuery = `SELECT 
	rc.id, rc.user_id, rc.created_at, rc.published_at, rc.current_version_id, rc.draft_version_id,
	rv.id, rv.recipe_id, rv.name, rv.prep_mins, rv.cook_mins, rv.portions, rv.created_at, rv.published_at, rv.version, rv.img_src, rv.description, rv.animal_product_level, rv.contains_gluten
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

// For creating an entirely new recipe we need a new container
func (r *recipeRepo) createRecipeContainer(ctx context.Context, tx *sql.Tx, recipeContainer *RecipeContainer) (*RecipeContainer, error) {
	var dbRecipeContainer RecipeContainer
	containerQuery := `INSERT INTO recipe_containers 
	(id, user_id, created_at, published_at) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id, user_id, created_at, published_at`
	err := tx.QueryRowContext(
		ctx,
		containerQuery,
		recipeContainer.ID,
		recipeContainer.UserID,
		recipeContainer.CreatedAt,
		recipeContainer.PublishedAt,
	).Scan(
		&dbRecipeContainer.ID,
		&dbRecipeContainer.UserID,
		&dbRecipeContainer.CreatedAt,
		&dbRecipeContainer.PublishedAt,
	)
	if err != nil {
		return nil, err
	}
	return &dbRecipeContainer, nil
}

// When a new version is created, call this to set the current version of the container to the new version id
func (r *recipeRepo) updateRecipeCurrentVersion(ctx context.Context, tx *sql.Tx, recipeID, versionID uuid.UUID) error {
	updateQuery := `UPDATE recipe_containers SET current_version_id = $1 WHERE id = $2`
	_, err := tx.ExecContext(ctx, updateQuery, versionID, recipeID)
	return err
}

// When a new draft version is created, call this to set the draft version id on the container. If publish is true, this should be set to null as there won't be a draft version
func (r *recipeRepo) updateRecipeDraftVersion(ctx context.Context, tx *sql.Tx, recipeID uuid.UUID, draftVersionID *uuid.UUID) error {
	updateQuery := `UPDATE recipe_containers SET draft_version_id = $1 WHERE id = $2`
	_, err := tx.ExecContext(ctx, updateQuery, draftVersionID, recipeID)
	return err
}

// When a recipe that has had drafts is published for the first time, need to set published on the container as well
// If a recipe is created published, then we handle that at instantiation
func (r *recipeRepo) setRecipePublishedAt(ctx context.Context, tx *sql.Tx, recipeID uuid.UUID, publishedAt *time.Time) error {
	updateQuery := `UPDATE recipe_containers SET published_at = $1 WHERE id = $2`
	_, err := tx.ExecContext(ctx, updateQuery, publishedAt, recipeID)
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

// For recipe listing ordered by created at
// Must only retrieve published and not deleted recipes
func (r *recipeRepo) getRecipesByCreatedAt(ctx context.Context, db db.DBTX, limit int, cursor *RecipeCursor, filter normalizedRecipeFilter, includeDrafts bool) ([]*RecipeListRow, error) {
	conditions := []string{"rc.deleted_on IS NULL"}
	if !includeDrafts {
		conditions = append(conditions, "rc.published_at IS NOT NULL")
	}
	args := make([]any, 0, 4)

	filterConditions, args := buildFilterConditions(filter, args)
	conditions = append(conditions, filterConditions...)

	if cursor != nil {
		conditions = append(conditions, fmt.Sprintf("(rc.created_at, rc.id) < ($%d, $%d)", len(args)+1, len(args)+2))
		args = append(args, cursor.CreatedAt, cursor.ID)
	}

	query := selectRecipeContainerWithVersionBaseQuery + `
	WHERE ` + strings.Join(conditions, " AND ") + `
	ORDER BY rc.created_at DESC, rc.id DESC
	LIMIT $` + fmt.Sprintf("%d", len(args)+1)
	args = append(args, limit)

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

// For recipe listing ordered by relevance to a search query
// Must only retrieve published and not deleted recipes
func (r *recipeRepo) getRecipesByRelevance(ctx context.Context, db db.DBTX, query string, limit int, cursor *RecipeCursor, filter normalizedRecipeFilter) ([]*RecipeListRow, error) {
	// Relevance score is a mix of full-text search ranking and trigram similarity, weighted towards full-text search. Sorting is done by relevance score first, then created_at and id to ensure a deterministic order.
	// this can be fine tuned in the future
	scoreExpression := `
		($1 * ts_rank_cd(to_tsvector('english', coalesce(rv.name, '')), websearch_to_tsquery('english', $3)) 
		+ $2 * similarity(lower(rv.name), lower($3))) AS relevance_score
	`

	args := []any{r.fullTextWeight, r.trigramWeight, query}
	filterConditions, args := buildFilterConditions(filter, args)
	userFilterClause := ""
	for _, c := range filterConditions {
		userFilterClause += "      AND " + c + "\n"
	}
	nextArg := len(args) + 1

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
        rv.description,
        rv.created_at AS version_created_at,
        rv.version,
		rv.img_src,
		rv.animal_product_level,
		rv.contains_gluten,
    ` + scoreExpression + `
    FROM recipe_containers rc
    JOIN recipe_versions rv ON rc.current_version_id = rv.id
    WHERE rc.deleted_on IS NULL
	AND rc.published_at IS NOT NULL
` + userFilterClause + `
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
	img_src,
    description,
	animal_product_level,
	contains_gluten,
    relevance_score
FROM ranked
`

	var q string

	if cursor != nil {
		if cursor.RelevanceScore == nil {
			return nil, ErrInvalidCursor
		}

		relevanceArg := nextArg
		createdAtArg := nextArg + 1
		idArg := nextArg + 2
		limitArg := nextArg + 3

		q = baseQuery + fmt.Sprintf(`
		WHERE (relevance_score, container_created_at, container_id) < ($%d, $%d, $%d)
		ORDER BY relevance_score DESC, container_created_at DESC, container_id DESC
		LIMIT $%d`, relevanceArg, createdAtArg, idArg, limitArg)

		args = append(args, *cursor.RelevanceScore, cursor.CreatedAt, cursor.ID, limit)
	} else {
		q = baseQuery + fmt.Sprintf(`
	ORDER BY relevance_score DESC, container_created_at DESC, container_id DESC
	LIMIT $%d`, nextArg)
		args = append(args, limit)
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
			&rv.ImgSrc,
			&rv.Description,
			&rv.AnimalProductLevel,
			&rv.ContainsGluten,
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

// buildFilterConditions appends SQL WHERE conditions and args for nf, starting
// from the current position in args. Returns the new conditions and extended args.
func buildFilterConditions(nf normalizedRecipeFilter, args []any) ([]string, []any) {
	var conditions []string
	if nf.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("rc.user_id = $%d", len(args)+1))
		args = append(args, *nf.UserID)
	}
	if nf.AnimalProductLevel != nil {
		conditions = append(conditions, fmt.Sprintf("rv.animal_product_level <= $%d", len(args)+1))
		args = append(args, *nf.AnimalProductLevel)
	}
	if nf.ContainsGluten != nil {
		conditions = append(conditions, fmt.Sprintf("rv.contains_gluten = $%d", len(args)+1))
		args = append(args, *nf.ContainsGluten)
	}
	return conditions, args
}

func scanRecipeContainerWithVersion(
	scanner interface{ Scan(dest ...any) error },
) (*RecipeContainer, error) {

	var rc RecipeContainer
	var rv RecipeVersion

	var draftVersionID uuid.NullUUID

	err := scanner.Scan(
		&rc.ID,
		&rc.UserID,
		&rc.CreatedAt,
		&rc.PublishedAt,
		&rc.CurrentVersionID,
		&draftVersionID,
		&rv.ID,
		&rv.RecipeID,
		&rv.Name,
		&rv.PrepMins,
		&rv.CookMins,
		&rv.Portions,
		&rv.CreatedAt,
		&rv.PublishedAt,
		&rv.Version,
		&rv.ImgSrc,
		&rv.Description,
		&rv.AnimalProductLevel,
		&rv.ContainsGluten,
	)
	if err != nil {
		return nil, err
	}

	rc.CurrentVersion = &rv
	if draftVersionID.Valid {
		rc.DraftVersionID = &draftVersionID.UUID
	}
	return &rc, nil
}
