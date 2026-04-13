package audit

import (
	"context"
	"database/sql"
	"foodplanner/internal/testutil/seeds"
	"testing"

	"foodplanner/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestLog_PersistsAuditEntry(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		service := NewAuditService(tx, NewRepo())

		testUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		entry, err := NewUserSignupEvent(uuid.New(), testUser.ID, testUser.Username, "127.0.0.1")
		require.NoError(t, err)

		err = service.Log(ctx, entry)
		require.NoError(t, err)

		var storedCount int
		err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM audits WHERE id = $1`, entry.ID).Scan(&storedCount)
		require.NoError(t, err)
		require.Equal(t, 1, storedCount)
	})
}
