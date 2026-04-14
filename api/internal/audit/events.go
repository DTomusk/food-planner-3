package audit

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

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
