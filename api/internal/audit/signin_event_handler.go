package audit

import (
	"context"

	"foodplanner/internal/db"
	"foodplanner/internal/events"
)

type SigninEventHandler struct {
	auditService *AuditService
}

func NewSigninEventHandler(auditService *AuditService) *SigninEventHandler {
	return &SigninEventHandler{auditService: auditService}
}

func (h *SigninEventHandler) Handle(ctx context.Context, tx db.DBTX, event events.Event) error {
	signinEvent, ok := event.(events.UserSignedInEvent)
	if !ok {
		return nil
	}

	entry, err := NewUserSigninEvent(
		signinEvent.Meta.CorrelationID,
		signinEvent.UserID,
		signinEvent.Username,
		signinEvent.IPAddress,
		signinEvent.UserAgent,
	)
	if err != nil {
		return err
	}

	return h.auditService.Log(context.WithoutCancel(ctx), tx, entry)
}
