package audit

import (
	"context"

	"foodplanner/internal/db"
	"foodplanner/internal/events"
)

type GraphQLRequestRejectedEventHandler struct {
	auditService *AuditService
}

func NewGraphQLRequestRejectedEventHandler(auditService *AuditService) *GraphQLRequestRejectedEventHandler {
	return &GraphQLRequestRejectedEventHandler{auditService: auditService}
}

func (h *GraphQLRequestRejectedEventHandler) Handle(ctx context.Context, tx db.DBTX, event events.Event) error {
	rejectedEvent, ok := event.(events.GraphQLRequestRejectedEvent)
	if !ok {
		return nil
	}

	entry, err := NewGraphQLRequestRejectedEvent(
		rejectedEvent.Meta.CorrelationID,
		rejectedEvent.Meta.ActorUserID,
		rejectedEvent.OperationName,
		rejectedEvent.OperationType,
		rejectedEvent.IPAddress,
		rejectedEvent.UserAgent,
		rejectedEvent.Path,
		rejectedEvent.Reason,
		rejectedEvent.QueryHash,
		rejectedEvent.MaxComplexity,
		rejectedEvent.PresentedMessage,
	)
	if err != nil {
		return err
	}

	return h.auditService.Log(context.WithoutCancel(ctx), tx, entry)
}
