package audit

import (
	"context"

	"foodplanner/internal/events"
)

type SignupEventHandler struct {
	auditService *AuditService
}

func NewSignupEventHandler(auditService *AuditService) *SignupEventHandler {
	return &SignupEventHandler{auditService: auditService}
}

func (h *SignupEventHandler) Handle(ctx context.Context, event events.Event) error {
	signupEvent, ok := event.(events.UserSignedUpEvent)
	if !ok {
		return nil
	}

	entry, err := NewUserSignupEvent(signupEvent.Meta.CorrelationID, signupEvent.UserID, signupEvent.Username, signupEvent.IPAddress)
	if err != nil {
		return err
	}

	return h.auditService.Log(context.WithoutCancel(ctx), entry)
}
