package audit

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

func NewRecipeCreatedEvent(correlationID, userID, recipeID, versionID uuid.UUID, ipAddress, userAgent string) (*AuditEntry, error) {
	createdAt := time.Now().UTC()

	newState, err := json.Marshal(struct {
		RecipeID  uuid.UUID `json:"recipe_id"`
		VersionID uuid.UUID `json:"version_id"`
		UserID    uuid.UUID `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
	}{
		RecipeID:  recipeID,
		VersionID: versionID,
		UserID:    userID,
		CreatedAt: createdAt,
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
		CreatedAt:     createdAt,
		NewState:      newState,
		Context:       contextData,
	}, nil
}

func NewUserSignupEvent(correlationID uuid.UUID, userID uuid.UUID, username, ipAddress string) (*AuditEntry, error) {
	createdAt := time.Now().UTC()

	newState, err := json.Marshal(struct {
		UserID    uuid.UUID `json:"user_id"`
		Username  string    `json:"username"`
		CreatedAt time.Time `json:"created_at"`
	}{
		UserID:    userID,
		Username:  username,
		CreatedAt: createdAt,
	})
	if err != nil {
		return nil, err
	}

	contextData, err := json.Marshal(struct {
		Source    string `json:"source"`
		Operation string `json:"operation"`
		IPAddress string `json:"ip_address"`
	}{
		Source:    "graphql",
		Operation: "signup",
		IPAddress: ipAddress,
	})
	if err != nil {
		return nil, err
	}

	return &AuditEntry{
		ID:            uuid.New(),
		CorrelationID: correlationID,
		ActorID:       &userID,
		ResourceType:  ResourceTypeUser,
		ResourceID:    &userID,
		Action:        ActionUserSignup,
		Result:        ResultSuccess,
		CreatedAt:     createdAt,
		NewState:      newState,
		Context:       contextData,
	}, nil
}

func NewUserSigninEvent(correlationID uuid.UUID, userID uuid.UUID, username, ipAddress, userAgent string) (*AuditEntry, error) {
	createdAt := time.Now().UTC()

	contextData, err := json.Marshal(struct {
		Source    string `json:"source"`
		Operation string `json:"operation"`
		IPAddress string `json:"ip_address"`
		UserAgent string `json:"user_agent"`
	}{
		Source:    "graphql",
		Operation: "signin",
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
		ResourceType:  ResourceTypeUser,
		ResourceID:    &userID,
		Action:        ActionUserSignin,
		Result:        ResultSuccess,
		CreatedAt:     createdAt,
		Context:       contextData,
	}, nil
}

func NewUserSigninFailureEvent(correlationID uuid.UUID, userID *uuid.UUID, email, ipAddress, userAgent, failureReason string) (*AuditEntry, error) {
	createdAt := time.Now().UTC()

	contextData, err := json.Marshal(struct {
		Source        string `json:"source"`
		Operation     string `json:"operation"`
		IPAddress     string `json:"ip_address"`
		UserAgent     string `json:"user_agent"`
		Email         string `json:"email"`
		FailureReason string `json:"failure_reason"`
	}{
		Source:        "graphql",
		Operation:     "signin",
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		Email:         email,
		FailureReason: failureReason,
	})
	if err != nil {
		return nil, err
	}

	return &AuditEntry{
		ID:            uuid.New(),
		CorrelationID: correlationID,
		ActorID:       userID,
		ResourceType:  ResourceTypeUser,
		ResourceID:    userID,
		Action:        ActionUserSignin,
		Result:        ResultFailure,
		CreatedAt:     createdAt,
		Reason:        &failureReason,
		Context:       contextData,
	}, nil
}
