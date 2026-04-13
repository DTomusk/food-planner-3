package audit

import (
	"context"
	"database/sql"
	"foodplanner/internal/testutil/seeds"
	"testing"
	"time"

	"foodplanner/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSave_PersistsAuditEntry(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		ctx := context.Background()
		repo := NewRepo()

		actorUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		resourceUser, err := seeds.SeedTestUser(ctx, tx)
		require.NoError(t, err)

		actorID := actorUser.ID
		resourceID := resourceUser.ID
		reason := "user created account"
		entry := &AuditEntry{
			ID:            uuid.New(),
			CorrelationID: uuid.New(),
			ActorID:       &actorID,
			ResourceType:  ResourceTypeUser,
			ResourceID:    &resourceID,
			Action:        ActionUserSignup,
			CreatedAt:     time.Now().UTC().Truncate(time.Microsecond),
			Result:        ResultSuccess,
			OldState:      []byte(`{"status":"pending"}`),
			NewState:      []byte(`{"status":"active"}`),
			Reason:        &reason,
			Context:       []byte(`{"source":"graphql","operation":"signup"}`),
		}

		err = repo.Save(ctx, tx, entry)
		require.NoError(t, err)

		var (
			storedID            uuid.UUID
			storedCorrelationID uuid.UUID
			storedActorID       sql.NullString
			storedResourceType  string
			storedResourceID    sql.NullString
			storedAction        string
			storedResult        string
			storedOldState      []byte
			storedNewState      []byte
			storedReason        sql.NullString
			storedContext       []byte
		)

		err = tx.QueryRowContext(ctx, `
			SELECT id, correlation_id, actor_id::text, resource_type, resource_id::text, action, result, old_state, new_state, reason, context
			FROM audits
			WHERE id = $1
		`, entry.ID).Scan(
			&storedID,
			&storedCorrelationID,
			&storedActorID,
			&storedResourceType,
			&storedResourceID,
			&storedAction,
			&storedResult,
			&storedOldState,
			&storedNewState,
			&storedReason,
			&storedContext,
		)
		require.NoError(t, err)

		require.Equal(t, entry.ID, storedID)
		require.Equal(t, entry.CorrelationID, storedCorrelationID)
		require.True(t, storedActorID.Valid)
		require.Equal(t, entry.ActorID.String(), storedActorID.String)
		require.Equal(t, string(entry.ResourceType), storedResourceType)
		require.True(t, storedResourceID.Valid)
		require.Equal(t, entry.ResourceID.String(), storedResourceID.String)
		require.Equal(t, string(entry.Action), storedAction)
		require.Equal(t, string(entry.Result), storedResult)
		require.JSONEq(t, string(entry.OldState), string(storedOldState))
		require.JSONEq(t, string(entry.NewState), string(storedNewState))
		require.True(t, storedReason.Valid)
		require.Equal(t, reason, storedReason.String)
		require.JSONEq(t, string(entry.Context), string(storedContext))
	})
}
