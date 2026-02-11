package db

import (
	"context"
	"database/sql"
)

// Executes the provided function within a transaction
// Roll back if function errors
func WithTx(
	ctx context.Context,
	db *sql.DB,
	fn func(*sql.Tx) error,
) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}
