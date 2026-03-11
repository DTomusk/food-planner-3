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

// Creates a new recipe including ingredient usages
// Returns recipe fields, ingredient usages can be loaded lazily
// TODO: orchestration should happen in service, we're leaking recipe version stuff into here
func (r *recipeRepo) createRecipe(ctx context.Context, tx *sql.Tx, recipeContainer *RecipeContainer) (*RecipeContainer, error) {
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

	var dbRecipeVersion RecipeVersion
	recipeVersion := recipeContainer.CurrentVersion
	err = tx.QueryRowContext(
		ctx,
		insertRecipeVersionQuery,
		recipeVersion.ID,
		recipeContainer.ID,
		recipeVersion.Name,
		recipeVersion.PrepMins,
		recipeVersion.CookMins,
		recipeVersion.Portions,
		recipeVersion.Version,
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

	setCurrentVersionQuery := `UPDATE recipe_containers SET current_version_id = $1 WHERE id = $2`
	_, err = tx.ExecContext(ctx, setCurrentVersionQuery, dbRecipeVersion.ID, dbRecipeContainer.ID)
	if err != nil {
		return nil, err
	}

	dbRecipeContainer.CurrentVersionID = dbRecipeVersion.ID
	dbRecipeContainer.CurrentVersion = &dbRecipeVersion

	// TODO: recipe sources should have its own repo
	if recipeVersion.Source != nil {
		sourceQuery := `INSERT INTO recipe_sources (recipe_version_id, type, url, book_title, book_page, instructions) VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = tx.ExecContext(ctx, sourceQuery, dbRecipeVersion.ID, int(recipeVersion.Source.Type), recipeVersion.Source.URL, recipeVersion.Source.BookTitle, recipeVersion.Source.BookPage, recipeVersion.Source.Instructions)
		if err != nil {
			return nil, err
		}
	}
	return &dbRecipeContainer, nil
}

func (r *recipeRepo) updateRecipeCurrentVersion(ctx context.Context, tx *sql.Tx, recipeID, versionID uuid.UUID) error {
	updateQuery := `UPDATE recipe_containers SET current_version_id = $1 WHERE id = $2`
	_, err := tx.ExecContext(ctx, updateQuery, versionID, recipeID)
	return err
}

func (r *recipeRepo) getRecipeByID(ctx context.Context, db db.DBTX, id string) (*RecipeContainer, error) {
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

func (r *recipeRepo) getAllRecipes(ctx context.Context, db db.DBTX) ([]*RecipeContainer, error) {
	rows, err := db.QueryContext(ctx,
		selectRecipeContainerWithVersionBaseQuery,
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

func (r *recipeRepo) getRecipesByUserID(ctx context.Context, db db.DBTX, userID string) ([]*RecipeContainer, error) {
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
