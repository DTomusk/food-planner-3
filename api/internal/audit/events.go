package audit

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

func NewUserSignupEvent(correlationID uuid.UUID, userID uuid.UUID, newState json.RawMessage) *AuditEntry {
	return &AuditEntry{
		ID:            uuid.New(),
		CorrelationID: correlationID,
		ActorID:       &userID,
		ResourceType:  ResourceTypeUser,
		ResourceID:    &userID,
		Action:        ActionUserSignup,
		Result:        ResultSuccess,
		CreatedAt:     time.Now(),
		NewState:      newState,
		Context:       nil,
	}
}
