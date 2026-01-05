package testutil

import (
	"database/sql"
	"testing"
)

func WithTx(t *testing.T, fn func(tx *sql.Tx)) {
	t.Helper()
	db := GetTestDB()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin tx: %v", err)
	}

	defer tx.Rollback()

	fn(tx)
}
