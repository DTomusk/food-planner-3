package audit

import (
	"context"
	"foodplanner/internal/db"
	"foodplanner/internal/events"
)

type RecipeUpdatedEventHandler struct {
	auditService *AuditService
}

func NewRecipeUpdatedEventHandler(auditService *AuditService) *RecipeUpdatedEventHandler {
	return &RecipeUpdatedEventHandler{auditService: auditService}
}

func (h *RecipeUpdatedEventHandler) Handle(ctx context.Context, tx db.DBTX, event events.Event) error {
	recipeUpdatedEvent, ok := event.(events.RecipeUpdatedEvent)
	if !ok {
		return nil
	}

	entry, err := NewRecipeUpdatedEvent(
		recipeUpdatedEvent.Meta.CorrelationID,
		recipeUpdatedEvent.UserID,
		recipeUpdatedEvent.RecipeID,
		recipeUpdatedEvent.VersionID,
		recipeUpdatedEvent.IPAddress,
		recipeUpdatedEvent.UserAgent,
		recipeUpdatedEvent.Meta.OccurredAtUTC,
	)
	if err != nil {
		return err
	}

	return h.auditService.Log(ctx, tx, entry)
}
