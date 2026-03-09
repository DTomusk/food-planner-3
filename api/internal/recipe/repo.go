package recipe

import (
	"context"
	"database/sql"
	"fmt"
	"foodplanner/internal/db"
	"strings"

	"github.com/google/uuid"
)

type Repo struct{}

func NewRepo() *Repo {
	return &Repo{}
}

// Creates a new recipe including ingredient usages
// Returns recipe fields, ingredient usages can be loaded lazily
func (r *Repo) CreateRecipe(ctx context.Context, tx *sql.Tx, recipeContainer *RecipeContainer) (*RecipeContainer, error) {
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
	RETURNING id, name, prep_mins, cook_mins, portions, created_at`
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

func (r *Repo) InsertIngredientUsages(ctx context.Context, tx *sql.Tx, ingredientUsages []*IngredientUsage, recipeVersionId uuid.UUID) error {
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

func (r *Repo) GetRecipeByID(ctx context.Context, db db.DBTX, id string) (*RecipeContainer, error) {
	row := db.QueryRowContext(ctx,
		`SELECT 
		rc.id, rc.user_id, rc.created_at, rc.current_version_id,
		rv.id, rv.name, rv.prep_mins, rv.cook_mins, rv.portions, rv.created_at, rv.version
		FROM recipe_containers rc
		JOIN recipe_versions rv ON rc.current_version_id = rv.id
		WHERE rc.id = $1`,
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

func (r *Repo) GetAllRecipes(ctx context.Context, db db.DBTX) ([]*RecipeContainer, error) {
	rows, err := db.QueryContext(ctx,
		`SELECT 
		rc.id, rc.user_id, rc.created_at, rc.current_version_id,
		rv.id, rv.name, rv.prep_mins, rv.cook_mins, rv.portions, rv.created_at, rv.version
		FROM recipe_containers rc
		JOIN recipe_versions rv ON rc.current_version_id = rv.id
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []*RecipeContainer
	for rows.Next() {
		rc, err := scanRecipeContainerWithVersion(rows)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, rc)
	}
	return recipes, nil
}

func (r *Repo) IngredientExists(ctx context.Context, db db.DBTX, ingredientID string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM ingredients WHERE id = $1)", ingredientID).Scan(&exists)
	return exists, err
}

func (r *Repo) GetIngredientUsagesForRecipe(ctx context.Context, db db.DBTX, recipeID string) ([]*IngredientUsage, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, ingredient_id, quantity, unit FROM ingredient_usages WHERE recipe_version_id = $1", recipeID)
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

func (r *Repo) GetRecipeSourceByRecipeID(ctx context.Context, db db.DBTX, recipeID string) (*RecipeSource, error) {
	var source RecipeSource
	row := db.QueryRowContext(ctx, "SELECT type, url, book_title, book_page, instructions FROM recipe_sources WHERE recipe_version_id = $1", recipeID)
	err := row.Scan(&source.Type, &source.URL, &source.BookTitle, &source.BookPage, &source.Instructions)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &source, nil
}

func (r *Repo) GetRecipesByUserID(ctx context.Context, db db.DBTX, userID string) ([]*RecipeContainer, error) {
	rows, err := db.QueryContext(ctx,
		`SELECT 
		rc.id, rc.user_id, rc.created_at, rc.current_version_id,
		rv.id, rv.name, rv.prep_mins, rv.cook_mins, rv.portions, rv.created_at, rv.version
		FROM recipe_containers rc
		JOIN recipe_versions rv ON rc.current_version_id = rv.id
		WHERE rc.user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []*RecipeContainer
	for rows.Next() {
		rc, err := scanRecipeContainerWithVersion(rows)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, rc)
	}
	return recipes, nil
}

func (r *Repo) GetRecipeVersionByID(ctx context.Context, db db.DBTX, id string) (*RecipeVersion, error) {
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
