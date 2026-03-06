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
func (r *Repo) CreateRecipe(ctx context.Context, tx *sql.Tx, recipeContainer *RecipeContainer, recipeVersion *RecipeVersion) (*RecipeVersion, error) {
	var dbRecipe RecipeVersion
	query := `INSERT INTO recipe_versions 
	(id, name, prep_mins, cook_mins, portions) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING id, name, prep_mins, cook_mins, portions`
	err := tx.QueryRowContext(
		ctx,
		query,
		recipeVersion.ID,
		recipeContainer.UserID,
		recipeVersion.Name,
		recipeVersion.PrepMins,
		recipeVersion.CookMins,
		recipeVersion.Portions,
	).Scan(
		&dbRecipe.ID,
		&dbRecipe.Name,
		&dbRecipe.PrepMins,
		&dbRecipe.CookMins,
		&dbRecipe.Portions,
	)
	if err != nil {
		return nil, err
	}
	usageQuery := `INSERT INTO ingredient_usages (id, recipe_version_id, ingredient_id, quantity, unit) VALUES ($1, $2, $3, $4, $5)`
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
	return &dbRecipe, nil
}

func (r *Repo) GetRecipeByID(ctx context.Context, db db.DBTX, id string) (*RecipeVersion, error) {
	var recipe RecipeVersion
	row := db.QueryRowContext(ctx,
		`SELECT id, name, prep_mins, cook_mins, portions
		FROM recipe_versions 
		WHERE id = $1
		AND deleted_on IS NULL`,
		id,
	)
	err := row.Scan(&recipe.ID, &recipe.Name, &recipe.PrepMins, &recipe.CookMins, &recipe.Portions)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &recipe, nil
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

func (r *Repo) GetAllRecipes(ctx context.Context, db db.DBTX) ([]*RecipeVersion, error) {
	rows, err := db.QueryContext(ctx,
		`SELECT id, name, prep_mins, cook_mins, portions
		FROM recipe_versions
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []*RecipeVersion
	for rows.Next() {
		var recipe RecipeVersion
		if err := rows.Scan(&recipe.ID, &recipe.Name, &recipe.PrepMins, &recipe.CookMins, &recipe.Portions); err != nil {
			return nil, err
		}
		recipes = append(recipes, &recipe)
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

func (r *Repo) GetRecipesByUserID(ctx context.Context, db db.DBTX, userID string) ([]*RecipeVersion, error) {
	rows, err := db.QueryContext(ctx,
		`SELECT id, user_id, name, prep_mins, cook_mins, portions, deleted_on
		FROM recipe_versions
		WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []*RecipeVersion
	for rows.Next() {
		var recipe RecipeVersion
		if err := rows.Scan(&recipe.ID, &recipe.Name, &recipe.PrepMins, &recipe.CookMins, &recipe.Portions); err != nil {
			return nil, err
		}
		recipes = append(recipes, &recipe)
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
