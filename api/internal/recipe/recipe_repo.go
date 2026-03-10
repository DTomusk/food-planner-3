package recipe

import (
	"context"
	"database/sql"
	"fmt"
	"foodplanner/internal/db"
	"strings"

	"github.com/google/uuid"
)

type RecipeRepo struct{}

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

func NewRecipeRepo() *RecipeRepo {
	return &RecipeRepo{}
}

// Creates a new recipe including ingredient usages
// Returns recipe fields, ingredient usages can be loaded lazily
func (r *RecipeRepo) CreateRecipe(ctx context.Context, tx *sql.Tx, recipeContainer *RecipeContainer) (*RecipeContainer, error) {
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
	versionQuery := `INSERT INTO recipe_versions 
	(id, recipe_id, name, prep_mins, cook_mins, portions, version) 
	VALUES ($1, $2, $3, $4, $5, $6, 1) 
	RETURNING id, recipe_id, name, prep_mins, cook_mins, portions, created_at`
	recipeVersion := recipeContainer.CurrentVersion
	err = tx.QueryRowContext(
		ctx,
		versionQuery,
		recipeVersion.ID,
		recipeContainer.ID,
		recipeVersion.Name,
		recipeVersion.PrepMins,
		recipeVersion.CookMins,
		recipeVersion.Portions,
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

	err = r.InsertIngredientUsages(ctx, tx, recipeVersion.Ingredients, dbRecipeVersion.ID)
	if err != nil {
		return nil, err
	}

	if recipeVersion.Source != nil {
		sourceQuery := `INSERT INTO recipe_sources (recipe_version_id, type, url, book_title, book_page, instructions) VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = tx.ExecContext(ctx, sourceQuery, dbRecipeVersion.ID, int(recipeVersion.Source.Type), recipeVersion.Source.URL, recipeVersion.Source.BookTitle, recipeVersion.Source.BookPage, recipeVersion.Source.Instructions)
		if err != nil {
			return nil, err
		}
	}
	return &dbRecipeContainer, nil
}

func (r *RecipeRepo) InsertIngredientUsages(ctx context.Context, tx *sql.Tx, ingredientUsages []*IngredientUsage, recipeVersionId uuid.UUID) error {
	if len(ingredientUsages) == 0 {
		return nil
	}
	var (
		query  = "INSERT INTO ingredient_usages (id, recipe_version_id, ingredient_id, quantity, unit) VALUES"
		args   []any
		values []string
	)

	for i, u := range ingredientUsages {
		start := i*5 + 1
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", start, start+1, start+2, start+3, start+4))
		args = append(args, u.ID, recipeVersionId, u.IngredientID, u.Quantity, int(u.Unit))
	}

	query += " " + strings.Join(values, ", ") + ";"

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *RecipeRepo) GetRecipeByID(ctx context.Context, db db.DBTX, id string) (*RecipeContainer, error) {
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

func (r *RecipeRepo) GetAllRecipes(ctx context.Context, db db.DBTX) ([]*RecipeContainer, error) {
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

func (r *RecipeRepo) GetIngredientUsagesForRecipeVersion(ctx context.Context, db db.DBTX, recipeVersionID string) ([]*IngredientUsage, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, ingredient_id, quantity, unit FROM ingredient_usages WHERE recipe_version_id = $1", recipeVersionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usages []*IngredientUsage
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for rows.Next() {
		var usage IngredientUsage
		if err := rows.Scan(&usage.ID, &usage.IngredientID, &usage.Quantity, &usage.Unit); err != nil {
			return nil, err
		}
		usages = append(usages, &usage)
	}
	return usages, nil
}

func (r *RecipeRepo) GetRecipesByUserID(ctx context.Context, db db.DBTX, userID string) ([]*RecipeContainer, error) {
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
