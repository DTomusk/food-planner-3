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

	usageQuery := `INSERT INTO ingredient_usages (id, recipe_version_id, ingredient_id, quantity, unit) VALUES ($1, $2, $3, $4, $5)`
	// TODO: may want to batch
	for _, usage := range recipeVersion.Ingredients {
		_, err := tx.ExecContext(ctx, usageQuery, usage.ID, recipeVersion.ID, usage.IngredientID, usage.Quantity, usage.Unit)
		if err != nil {
			return nil, err
		}
	}
	if recipeVersion.Source != nil {
		sourceQuery := `INSERT INTO recipe_sources (recipe_version_id, type, url, book_title, book_page, instructions) VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = tx.ExecContext(ctx, sourceQuery, recipeVersion.ID, int(recipeVersion.Source.Type), recipeVersion.Source.URL, recipeVersion.Source.BookTitle, recipeVersion.Source.BookPage, recipeVersion.Source.Instructions)
		if err != nil {
			return nil, err
		}
	}
	return &dbRecipeContainer, nil
}

func (r *Repo) GetRecipeByID(ctx context.Context, db db.DBTX, id string) (*RecipeContainer, error) {
	var recipeContainer RecipeContainer
	var recipeVersion RecipeVersion
	row := db.QueryRowContext(ctx,
		`SELECT 
		rc.id, rc.user_id, rc.created_at, rc.current_version_id,
		rv.id, rv.name, rv.prep_mins, rv.cook_mins, rv.portions, rv.created_at, rv.version
		FROM recipe_containers rc
		JOIN recipe_versions rv ON rc.current_version_id = rv.id
		WHERE rc.id = $1`,
		id,
	)
	err := row.Scan(&recipeContainer.ID, &recipeContainer.UserID, &recipeContainer.CreatedAt, &recipeContainer.CurrentVersionID,
		&recipeVersion.ID, &recipeVersion.Name, &recipeVersion.PrepMins, &recipeVersion.CookMins, &recipeVersion.Portions, &recipeVersion.CreatedAt, &recipeVersion.Version)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	recipeContainer.CurrentVersion = &recipeVersion
	return &recipeContainer, nil
}

// func (r *Repo) GetDeletedRecipeByID(ctx context.Context, db db.DBTX, id string) (*RecipeVersion, error) {
// 	var recipe RecipeVersion
// 	row := db.QueryRowContext(ctx,
// 		`SELECT id, name, prep_mins, cook_mins, portions
// 		FROM recipe_versions
// 		WHERE id = $1
// 		AND deleted_on IS NOT NULL`,
// 		id,
// 	)
// 	err := row.Scan(&recipe.ID, &recipe.Name, &recipe.PrepMins, &recipe.CookMins, &recipe.Portions)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return &recipe, nil
// }

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
		var recipeContainer RecipeContainer
		var recipeVersion RecipeVersion
		if err := rows.Scan(&recipeContainer.ID, &recipeContainer.UserID, &recipeContainer.CreatedAt, &recipeContainer.CurrentVersionID,
			&recipeVersion.ID, &recipeVersion.Name, &recipeVersion.PrepMins, &recipeVersion.CookMins, &recipeVersion.Portions, &recipeVersion.CreatedAt, &recipeVersion.Version); err != nil {
			return nil, err
		}
		recipeContainer.CurrentVersion = &recipeVersion
		recipes = append(recipes, &recipeContainer)
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

// func (r *Repo) DeleteRecipe(ctx context.Context, db db.DBTX, recipeID string) (*RecipeVersion, error) {
// 	var recipe RecipeVersion
// 	row := db.QueryRowContext(ctx, "UPDATE recipe_versions SET deleted_on = NOW() WHERE id = $1 RETURNING id, user_id, name, prep_mins, cook_mins, portions, deleted_on", recipeID)
// 	err := row.Scan(&recipe.ID, &recipe.UserID, &recipe.Name, &recipe.PrepMins, &recipe.CookMins, &recipe.Portions, &recipe.DeletedOn)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &recipe, nil
// }

// func (r *Repo) UndeleteRecipe(ctx context.Context, db db.DBTX, recipeID string) (*RecipeVersion, error) {
// 	var recipe RecipeVersion
// 	row := db.QueryRowContext(ctx, "UPDATE recipe_versions SET deleted_on = NULL WHERE id = $1 RETURNING id, user_id, name, prep_mins, cook_mins, portions, deleted_on", recipeID)
// 	err := row.Scan(&recipe.ID, &recipe.UserID, &recipe.Name, &recipe.PrepMins, &recipe.CookMins, &recipe.Portions, &recipe.DeletedOn)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &recipe, nil
// }

// func (s *Repo) DeleteOldRecipes(ctx context.Context, db db.DBTX, retentionDays int) (int64, error) {
// 	result, err := db.ExecContext(ctx, "DELETE FROM recipe_versions WHERE deleted_on IS NOT NULL AND deleted_on < NOW() - $1 * INTERVAL '1 day'", retentionDays)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return result.RowsAffected()
// }

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
		var recipeContainer RecipeContainer
		var recipeVersion RecipeVersion
		if err := rows.Scan(&recipeContainer.ID, &recipeContainer.UserID, &recipeContainer.CreatedAt, &recipeContainer.CurrentVersionID,
			&recipeVersion.ID, &recipeVersion.Name, &recipeVersion.PrepMins, &recipeVersion.CookMins, &recipeVersion.Portions, &recipeVersion.CreatedAt, &recipeVersion.Version); err != nil {
			return nil, err
		}
		recipeContainer.CurrentVersion = &recipeVersion
		recipes = append(recipes, &recipeContainer)
	}
	return recipes, nil
}

// func (r *Repo) GetDeletedRecipesByUserID(ctx context.Context, db db.DBTX, userID string) ([]*RecipeVersion, error) {
// 	rows, err := db.QueryContext(ctx,
// 		`SELECT id, user_id, name, prep_mins, cook_mins, portions, deleted_on
// 		FROM recipe_versions
// 		WHERE user_id = $1
// 		AND deleted_on IS NOT NULL`,
// 		userID,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var recipes []*RecipeVersion
// 	for rows.Next() {
// 		var recipe RecipeVersion
// 		if err := rows.Scan(&recipe.ID, &recipe.UserID, &recipe.Name, &recipe.PrepMins, &recipe.CookMins, &recipe.Portions, &recipe.DeletedOn); err != nil {
// 			return nil, err
// 		}
// 		recipes = append(recipes, &recipe)
// 	}
// 	return recipes, nil
// }

func (r *Repo) GetRecipeVersionByID(ctx context.Context, db db.DBTX, id string) (*RecipeVersion, error) {
	var recipeVersion RecipeVersion
	row := db.QueryRowContext(ctx, `SELECT id, recipe_id, name, prep_mins, cook_mins, portions, created_at, version FROM recipe_versions WHERE id = $1`, id)
	err := row.Scan(&recipeVersion.ID, &recipeVersion.RecipeID, &recipeVersion.Name, &recipeVersion.PrepMins, &recipeVersion.CookMins, &recipeVersion.Portions, &recipeVersion.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &recipeVersion, nil
}
