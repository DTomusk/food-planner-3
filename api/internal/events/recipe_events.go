package events

import (
	"time"

	"github.com/google/uuid"
)

const (
	RecipeCreatedEventType = "recipe.created"
	RecipeUpdatedEventType = "recipe.updated"
)

// Record basic data about a recipe being created
// We don't record content (name, description, ingredients) as recipe versions are immutable
type RecipeCreatedEvent struct {
	Meta      Metadata
	RecipeID  uuid.UUID
	VersionID uuid.UUID
	UserID    uuid.UUID
	IPAddress string
	UserAgent string
}

func (e RecipeCreatedEvent) Metadata() Metadata {
	return e.Meta
}

func NewRecipeCreatedEvent(correlationID, recipeID, versionID, userID uuid.UUID, ipAddress, userAgent string) RecipeCreatedEvent {
	actorID := userID

	return RecipeCreatedEvent{
		Meta: Metadata{
			ID:            uuid.New(),
			Type:          RecipeCreatedEventType,
			Version:       1,
			OccurredAtUTC: time.Now().UTC(),
			CorrelationID: correlationID,
			ActorUserID:   &actorID,
			Source:        RecipeServiceSource,
		},
		RecipeID:  recipeID,
		VersionID: versionID,
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
}

type RecipeUpdatedEvent struct {
	Meta      Metadata
	RecipeID  uuid.UUID
	VersionID uuid.UUID
	UserID    uuid.UUID
	IPAddress string
	UserAgent string
}

func (e RecipeUpdatedEvent) Metadata() Metadata {
	return e.Meta
}

func NewRecipeUpdatedEvent(correlationID, recipeID, versionID, userID uuid.UUID, ipAddress, userAgent string) RecipeUpdatedEvent {
	actorID := userID

	return RecipeUpdatedEvent{
		Meta: Metadata{
			ID:            uuid.New(),
			Type:          RecipeUpdatedEventType,
			Version:       1,
			OccurredAtUTC: time.Now().UTC(),
			CorrelationID: correlationID,
			ActorUserID:   &actorID,
			Source:        RecipeServiceSource,
		},
		RecipeID:  recipeID,
		VersionID: versionID,
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
}
