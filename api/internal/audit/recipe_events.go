package audit

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

func NewRecipeCreatedEvent(correlationID, userID, recipeID, versionID uuid.UUID, ipAddress, userAgent string, occurredAtUTC time.Time) (*AuditEntry, error) {

	newState, err := json.Marshal(struct {
		RecipeID  uuid.UUID `json:"recipe_id"`
		VersionID uuid.UUID `json:"version_id"`
		UserID    uuid.UUID `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
	}{
		RecipeID:  recipeID,
		VersionID: versionID,
		UserID:    userID,
		CreatedAt: occurredAtUTC,
	})
	if err != nil {
		return nil, err
	}

	contextData, err := json.Marshal(struct {
		Source    string `json:"source"`
		Operation string `json:"operation"`
		IPAddress string `json:"ip_address"`
		UserAgent string `json:"user_agent"`
	}{
		Source:    "graphql",
		Operation: "recipe_created",
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})
	if err != nil {
		return nil, err
	}

	return &AuditEntry{
		ID:            uuid.New(),
		CorrelationID: correlationID,
		ActorID:       &userID,
		ResourceType:  ResourceTypeRecipe,
		ResourceID:    &recipeID,
		Action:        ActionRecipeCreated,
		Result:        ResultSuccess,
		CreatedAt:     time.Now().UTC(),
		NewState:      newState,
		Context:       contextData,
	}, nil
}

func NewRecipeUpdatedEvent(correlationID, userID, recipeID, versionID uuid.UUID, ipAddress, userAgent string, occurredAtUTC time.Time) (*AuditEntry, error) {
	newState, err := json.Marshal(struct {
		RecipeID  uuid.UUID `json:"recipe_id"`
		VersionID uuid.UUID `json:"version_id"`
		UserID    uuid.UUID `json:"user_id"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		RecipeID:  recipeID,
		VersionID: versionID,
		UserID:    userID,
		UpdatedAt: occurredAtUTC,
	})
	if err != nil {
		return nil, err
	}

	contextData, err := json.Marshal(struct {
		Source    string `json:"source"`
		Operation string `json:"operation"`
		IPAddress string `json:"ip_address"`
		UserAgent string `json:"user_agent"`
	}{
		Source:    "graphql",
		Operation: "recipe_updated",
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})
	if err != nil {
		return nil, err
	}

	return &AuditEntry{
		ID:            uuid.New(),
		CorrelationID: correlationID,
		ActorID:       &userID,
		ResourceType:  ResourceTypeRecipe,
		ResourceID:    &recipeID,
		Action:        ActionRecipeUpdated,
		Result:        ResultSuccess,
		CreatedAt:     time.Now().UTC(),
		NewState:      newState,
		Context:       contextData,
	}, nil
}
