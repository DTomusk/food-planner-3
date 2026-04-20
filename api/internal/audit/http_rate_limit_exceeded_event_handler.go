package audit

import (
	"context"
	"foodplanner/internal/db"
	"foodplanner/internal/events"
)

type RateLimitExceededEventHandler struct {
	auditService *AuditService
}

func NewRateLimitExceededEventHandler(auditService *AuditService) *RateLimitExceededEventHandler {
	return &RateLimitExceededEventHandler{
		auditService: auditService,
	}
}

func (h *RateLimitExceededEventHandler) Handle(ctx context.Context, tx db.DBTX, event events.Event) error {
	rateLimitEvent, ok := event.(events.RateLimitExceededEvent)
	if !ok {
		return nil
	}

	entry, err := NewRateLimitExceededEvent(
		rateLimitEvent.Meta.CorrelationID,
		rateLimitEvent.Meta.ActorUserID,
		rateLimitEvent.Subject,
		string(rateLimitEvent.SubjectType),
		rateLimitEvent.IPAddress,
		rateLimitEvent.UserAgent,
		rateLimitEvent.Method,
		rateLimitEvent.Path,
		rateLimitEvent.Limit,
		rateLimitEvent.Count,
		rateLimitEvent.WindowSeconds,
		rateLimitEvent.RetryAfterSeconds,
		rateLimitEvent.Meta.OccurredAtUTC,
	)
	if err != nil {
		return err
	}

	return h.auditService.Log(context.WithoutCancel(ctx), tx, entry)
}
