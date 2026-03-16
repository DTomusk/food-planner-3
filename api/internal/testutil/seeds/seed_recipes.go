package seeds

import (
	"context"
	"foodplanner/internal/db"

	"github.com/google/uuid"
)

func InsertRecipeContainer(ctx context.Context, db db.DBTX, id, userID uuid.UUID) error {
	query := `INSERT INTO recipe_containers (id, user_id) VALUES ($1, $2)`
	_, err := db.ExecContext(ctx, query, id, userID)
	return err
}

func InsertRecipeVersion(ctx context.Context, db db.DBTX, id, recipeID uuid.UUID, name string, prepMins, cookMins, portions, version int) error {
	query := `INSERT INTO recipe_versions
	(id, recipe_id, name, prep_mins, cook_mins, portions, version) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.ExecContext(ctx, query, id, recipeID, name, prepMins, cookMins, portions, version)
	return err
}

func SetRecipeContainerCurrentVersion(ctx context.Context, db db.DBTX, recipeID, currentVersionID uuid.UUID) error {
	query := `UPDATE recipe_containers SET current_version_id = $1 WHERE id = $2`
	_, err := db.ExecContext(ctx, query, currentVersionID, recipeID)
	return err
}

func InsertRecipeSource(ctx context.Context, db db.DBTX, recipeVersionID uuid.UUID, sourceType int, url, bookTitle *string, bookPage *int, instructions *string) error {
	query := `INSERT INTO recipe_sources (recipe_version_id, type, url, book_title, book_page, instructions) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.ExecContext(ctx, query, recipeVersionID, sourceType, url, bookTitle, bookPage, instructions)
	return err
}
