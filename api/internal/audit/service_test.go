package audit

import (
	"context"
	"database/sql"
	"testing"

	"foodplanner/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestLog_PersistsAuditEntry(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		service := NewAuditService(tx, NewRepo())

		entry, err := NewUserSignupEvent(uuid.New(), uuid.New(), "test-user", "127.0.0.1")
		require.NoError(t, err)

		err = service.Log(ctx, entry)
		require.NoError(t, err)

		var storedCount int
		err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM audits WHERE id = $1`, entry.ID).Scan(&storedCount)
		require.NoError(t, err)
		require.Equal(t, 1, storedCount)
	})
}
