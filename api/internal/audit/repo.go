package audit

import (
	"context"
	"foodplanner/internal/db"
)

type repo struct{}

const insertAuditQuery = `
INSERT INTO audits (
	id,
	correlation_id,
	actor_id,
	resource_type,
	resource_id,
	action,
	created_at,
	result,
	old_state,
	new_state,
	reason,
	context
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
`

func NewRepo() *repo {
	return &repo{}
}

func (r *repo) Save(ctx context.Context, database db.DBTX, entry *AuditEntry) error {
	_, err := database.ExecContext(
		ctx,
		insertAuditQuery,
		entry.ID,
		entry.CorrelationID,
		entry.ActorID,
		entry.ResourceType,
		entry.ResourceID,
		entry.Action,
		entry.CreatedAt,
		entry.Result,
		entry.OldState,
		entry.NewState,
		entry.Reason,
		entry.Context,
	)

	return err
}
