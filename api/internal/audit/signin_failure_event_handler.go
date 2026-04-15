package audit

import (
	"context"

	"foodplanner/internal/db"
	"foodplanner/internal/events"
)

type SigninFailureEventHandler struct {
	auditService *AuditService
}

func NewSigninFailureEventHandler(auditService *AuditService) *SigninFailureEventHandler {
	return &SigninFailureEventHandler{auditService: auditService}
}

func (h *SigninFailureEventHandler) Handle(ctx context.Context, tx db.DBTX, event events.Event) error {
	signinFailureEvent, ok := event.(events.UserSigninFailedEvent)
	if !ok {
		return nil
	}

	entry, err := NewUserSigninFailureEvent(
		signinFailureEvent.Meta.CorrelationID,
		signinFailureEvent.UserID,
		signinFailureEvent.Email,
		signinFailureEvent.IPAddress,
		signinFailureEvent.FailureReason,
	)
	if err != nil {
		return err
	}

	return h.auditService.Log(context.WithoutCancel(ctx), tx, entry)
}
