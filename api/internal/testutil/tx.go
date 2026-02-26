package testutil

import (
	"context"
	"database/sql"
	"testing"

	"foodplanner/internal/db"
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

type TestTxRunner struct {
	tx *sql.Tx
}

func NewTestTxRunner(tx *sql.Tx) *TestTxRunner {
	return &TestTxRunner{tx: tx}
}

func (t *TestTxRunner) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	return fn(t.tx)
}

func (t *TestTxRunner) DB() db.DBTX {
	return t.tx
}
