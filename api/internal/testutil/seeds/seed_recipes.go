package seeds

import (
	"context"
	"foodplanner/internal/db"
	"time"
)

func InsertRecipe(ctx context.Context, db db.DBTX, id, userID, name string, prepMins, cookMins, portions int, deletedOn *time.Time) error {
	query := `INSERT INTO recipe_versions 
	(id, user_id, name, prep_mins, cook_mins, portions, deleted_on) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.ExecContext(ctx, query, id, userID, name, prepMins, cookMins, portions, deletedOn)
	return err
}
