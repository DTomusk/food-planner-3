package audit

import (
	"context"
	"foodplanner/internal/db"
	"foodplanner/internal/events"
)

type RecipeCreatedEventHandler struct {
	auditService *AuditService
}

func NewRecipeCreatedEventHandler(auditService *AuditService) *RecipeCreatedEventHandler {
	return &RecipeCreatedEventHandler{auditService: auditService}
}

func (h *RecipeCreatedEventHandler) Handle(ctx context.Context, tx db.DBTX, event events.Event) error {
	recipeCreatedEvent, ok := event.(events.RecipeCreatedEvent)
	if !ok {
		return nil
	}

	entry, err := NewRecipeCreatedEvent(
		recipeCreatedEvent.Meta.CorrelationID,
		recipeCreatedEvent.UserID,
		recipeCreatedEvent.RecipeID,
		recipeCreatedEvent.VersionID,
		recipeCreatedEvent.IPAddress,
		recipeCreatedEvent.UserAgent,
	)
	if err != nil {
		return err
	}

	return h.auditService.Log(ctx, tx, entry)
}
